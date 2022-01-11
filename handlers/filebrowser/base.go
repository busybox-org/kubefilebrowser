package filebrowser

import (
	"github.com/sirupsen/logrus"
	"github.com/xmapst/kubefilebrowser/internal"
	"io"
	"strings"
)

type Query struct {
	Namespace string    `json:"namespace" form:"namespace" binding:"required"`
	Pod       string    `json:"pod" form:"pod" binding:"required" binding:"required"`
	Container string    `json:"container" form:"container" binding:"required"`
	Path      string    `json:"path" form:"path" binding:"required" binding:"required"`
	OldPath   string    `json:"old_path,omitempty" form:"old_path"`
	Command   []string  `json:"-"`
	Stdin     io.Reader `json:"-"`
}

func (q *Query) fileBrowser() (res []byte, err error) {
	pb := internal.PodBase{
		Namespace: q.Namespace,
		PodName:   q.Pod,
		Container: q.Container,
	}
	res, err = pb.Exec(q.Stdin, q.Command...)
	if err != nil {
		if strings.Contains(err.Error(), "kf_tools") ||
			err.Error() == "command terminated with exit code 126" {
			err = pb.InstallKFTools()
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			res, err = pb.Exec(q.Stdin, q.Command...)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
		} else {
			logrus.Error(err)
			return nil, err
		}
	}
	return res, err
}
