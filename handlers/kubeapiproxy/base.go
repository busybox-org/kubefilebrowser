package kubeapiproxy

import (
	"github.com/sirupsen/logrus"
	"github.com/xmapst/kubefilebrowser/internal"
	"io"
	"os"
	"sync"
)

type copyReq struct {
	namespace string
	pod       string
	srcPath   string
	destPath  string
}

type copyRes struct {
	Pod       string `json:"pod"`
	Container string `json:"container"`
	State     bool   `json:"state"`
	Msg       string `json:"msg"`
}

func (c *copyReq) toPodContainer(container string) (res *copyRes) {
	pb := internal.PodBase{
		Namespace: c.namespace,
		PodName:   c.pod,
		Container: container,
	}
	reader, writer := io.Pipe()
	cp := pb.NewPodExec()
	cp.Stdin = reader

	go func() {
		defer func() {
			_ = writer.Close()
		}()
		tarFile, err := os.Open(c.srcPath)
		if err != nil {
			logrus.Error(err)
			return
		}
		_, err = io.Copy(writer, tarFile)
		if err != nil {
			logrus.Error(err)
		}
	}()
	res = &copyRes{
		Pod:       c.pod,
		Container: container,
		State:     true,
	}
	err := cp.ToPodContainer(c.destPath)
	if err != nil {
		res.State = false
		res.Msg = err.Error()
	}
	return
}

func (c *copyReq) toPodContainers(containers []string) (res []*copyRes) {
	var wg = new(sync.WaitGroup)
	var lock = new(sync.Mutex)
	for _, container := range containers {
		wg.Add(1)
		go func(container string) {
			defer wg.Done()
			_res := c.toPodContainer(container)
			lock.Lock()
			defer lock.Unlock()
			res = append(res, _res)
		}(container)
	}
	wg.Wait()
	return
}
