package handlers

import (
	"archive/tar"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	"github.com/xmapst/kubefilebrowser/configs"
	"github.com/xmapst/kubefilebrowser/utils"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	CodeSuccess = 0
	CodeErrApp  = iota + 1000
	CodeErrMsg
	CodeErrParam
	CodeErrNoPriv
)

var MsgFlags = map[int]string{
	CodeErrApp:    "内部错误",
	CodeSuccess:   "成功",
	CodeErrMsg:    "未知错误",
	CodeErrParam:  "参数错误",
	CodeErrNoPriv: "沒有权限",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[CodeErrApp]
}

type Gin struct {
	*gin.Context
}

type JSONResult struct {
	Code    int         `json:"code" description:"返回码" example:"0000"`
	Message string      `json:"message,omitempty" description:"消息" example:"消息"`
	Data    interface{} `json:"data" description:"数据"`
}

type Info struct {
	Ok      bool   `json:"ok" description:"状态" example:"true"`
	Message string `json:"message,omitempty" description:"消息" example:"消息"`
}

func NewRes(data interface{}, err error, code int) *JSONResult {
	if code == 200 {
		code = 0
	}
	codeMsg := ""
	if configs.Config.RunMode == gin.ReleaseMode && code != 0 {
		codeMsg = GetMsg(code)
	}

	return &JSONResult{
		Data: data,
		Code: code,
		Message: func() string {
			result := NewInfo(err)
			if codeMsg != "" && result != "" {
				result += ", " + codeMsg
			} else if codeMsg != "" {
				result = codeMsg
			}
			return strings.TrimSpace(result)
		}(),
	}
}

func NewInfo(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// SetRes Response res
func (g *Gin) SetRes(res interface{}, err error, code int) {
	g.JSON(http.StatusOK, NewRes(res, err, code))
}

// SetJson Set Json
func (g *Gin) SetJson(res interface{}) {
	g.SetRes(res, nil, CodeSuccess)
}

// SetError Check Error
func (g *Gin) SetError(code int, err error) {
	g.SetRes(nil, err, code)
	g.Abort()
}

// SaveToTarFile 保存为tar压缩文件
func (g *Gin) SaveToTarFile(filePath string) error {
	form, err := g.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]
	if len(files) == 0 {
		return fmt.Errorf("files is null")
	}
	fw, err := os.Create(filePath)
	if err != nil {
		return err
	}

	tw := tar.NewWriter(fw)
	defer func() {
		_ = tw.Close()
		_ = fw.Close()
	}()
	for _, f := range files {
		//file.Filename does not contain the directory path
		// RFC 7578, Section 4.2 requires that if a filename is provided, the
		// directory path information must not be used.
		// 🙂🙂🙂🙂🙂🙂
		v := f.Header.Get("Content-Disposition")
		_, dispositionParams, err := mime.ParseMediaType(v)
		if err != nil {
			return err
		}
		fileName, ok := dispositionParams["filename"]
		if !ok {
			return fmt.Errorf("filename does not exist")
		}

		hdr := &tar.Header{
			Name: fileName,
			Mode: 0644,
			Size: f.Size,
		}
		err = tw.WriteHeader(hdr)
		if err != nil {
			return err
		}
		_f, err := f.Open()
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, _f)
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveFiles 保存为本地文件
func (g *Gin) SaveFiles(filePath string) error {
	form, err := g.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]
	if len(files) == 0 {
		return fmt.Errorf("files is null")
	}
	var wg = new(sync.WaitGroup)
	var errs error

	// create _tmpSaveDir
	if !utils.FileOrPathExist(filePath) {
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	for _, f := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()
			//file.Filename does not contain the directory path
			// RFC 7578, Section 4.2 requires that if a filename is provided, the
			// directory path information must not be used.
			// 🙂🙂🙂🙂🙂🙂
			v := file.Header.Get("Content-Disposition")
			_, dispositionParams, err := mime.ParseMediaType(v)
			if err != nil {
				logrus.Error(err)
				return
			}
			fileName, ok := dispositionParams["filename"]
			if !ok {
				logrus.Error("filename does not exist")
				return
			}

			file.Filename = fileName
			_filePath := filePath
			// Default save path
			uploadFileName := filepath.Base(file.Filename)
			uploadFPath := filepath.Dir(file.Filename)
			// Process folder upload
			if uploadFPath != "." {
				_filePath = filepath.Join(filePath, uploadFPath)
				if !utils.FileOrPathExist(_filePath) {
					_ = os.MkdirAll(_filePath, os.ModePerm)
				}
			}

			// save file to local in _tp
			err = g.SaveUploadedFile(file, filepath.Join(_filePath, uploadFileName))
			if err != nil {
				errs = multierror.Append(errs, fmt.Errorf(file.Filename, err.Error()))
			}
		}(f)
	}
	wg.Wait()
	if errs != nil {
		return errs
	}
	return nil
}
