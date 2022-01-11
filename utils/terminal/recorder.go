package terminal

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xmapst/kubefilebrowser/configs"
	"github.com/xmapst/kubefilebrowser/utils"
	"os"
	"path/filepath"
	"time"
)

const header = `{"version":2,"width":%d,"height":%d,"timestamp":%d,"env":{"SHELL":"%s","TERM":"xterm-256color"}}\n`

type ReplyRecorder struct {
	sessionID     string
	shell         string
	timeStartNano int64
	absFilePath   string
	AbsGzFilePath string
	file          *os.File
	width         uint16
	height        uint16
}

func NewReplyRecord(sid, shell string) (recorder ReplyRecorder) {
	recorder = ReplyRecorder{sessionID: sid, shell: shell}
	recorder.initial()
	return recorder
}

func (r *ReplyRecorder) initial() {
	r.prepare()
	r.Record([]byte("Auditing is turned on, recording is turned on\r\n"))
	r.Record([]byte("Auditing is turned on, recording is turned on\r\n"))
	r.Record([]byte("Auditing is turned on, recording is turned on\r\n"))
}

func (r *ReplyRecorder) Record(b []byte) {
	if len(b) > 0 {
		delta := float64(time.Now().UnixNano()-r.timeStartNano) / 1000 / 1000 / 1000
		data, _ := json.Marshal(string(b))
		_, err := r.file.WriteString(fmt.Sprintf(`[%f,"o",%s]`, delta, data))
		if err != nil {
			logrus.Error(fmt.Sprintf("Session %s write replay to file failed: %s", r.sessionID, err))
		}
		_, _ = r.file.WriteString("\n")
	}
}

func (r *ReplyRecorder) prepare() {
	sessionID := r.sessionID
	rootPath := configs.Config.RootPath
	today := time.Now().UTC().Format("2006-01-02")
	gzFileName := sessionID + ".replay.gz"
	replayDir := filepath.Join(rootPath, "data", "replays", today)

	r.absFilePath = filepath.Join(replayDir, sessionID)
	r.AbsGzFilePath = filepath.Join(replayDir, gzFileName)
	r.timeStartNano = time.Now().UnixNano()

	err := utils.EnsureDirExist(replayDir)
	if err != nil {
		logrus.Error(fmt.Sprintf("Create dir %s error: %s\n", replayDir, err))
		return
	}

	r.file, err = os.Create(r.absFilePath)
	if err != nil {
		logrus.Error(fmt.Sprintf("Create file %s error: %s", r.absFilePath, err))
	}
	_, _ = r.file.WriteString("\n")
}

func (r *ReplyRecorder) End() {
	delta := float64(time.Now().UnixNano()-r.timeStartNano) / 1000 / 1000 / 1000
	_, _ = r.file.WriteString(fmt.Sprintf(`[%f,"o","Goodbye, record is saved\r\n"]`, delta))
	_ = r.file.Close()
	_ = utils.InsertStringToFile(r.absFilePath, fmt.Sprintf(header, r.width, r.height, time.Now().Unix(), r.shell), 1)
	if !utils.FileOrPathExist(r.AbsGzFilePath) {
		_ = utils.GzipCompressFile(r.absFilePath, r.AbsGzFilePath)
		_ = os.Remove(r.absFilePath)
	}
	//go r.UploadGzipFile()
}
