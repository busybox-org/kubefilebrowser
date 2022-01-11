package utils

import (
	"archive/tar"
	"embed"
	_ "embed"
	"github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"path"
)

//go:embed kftoolsbinary
var kfBinaryEmbedFiles embed.FS

func TarKFTools(name string, writer io.Writer) error {
	fSys, err := fs.Sub(kfBinaryEmbedFiles, "kftoolsbinary")
	if err != nil {
		logrus.Error(err)
		return err
	}
	tw := tar.NewWriter(writer)
	// 如果关闭失败会造成tar包不完整
	defer tw.Close()
	f, err := fSys.Open(path.Base(name))
	if err != nil {
		logrus.Error(err)
		return err
	}
	f.Close()
	fi, err := f.Stat()
	if err != nil {
		logrus.Error(err)
		return err
	}
	hdr, err := tar.FileInfoHeader(fi, name)
	if err != nil {
		logrus.Error(err)
		return err
	}
	hdr.Name = "kftools"
	// 将tar的文件信息hdr写入到tw
	err = tw.WriteHeader(hdr)
	if err != nil {
		logrus.Error(err)
		return err
	}
	// 将文件数据写入
	_, err = io.Copy(tw, f)

	return err
}
