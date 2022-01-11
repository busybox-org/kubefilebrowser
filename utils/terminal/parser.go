package terminal

import (
	"github.com/sirupsen/logrus"
	"github.com/xmapst/kubefilebrowser/utils/terminalparser"
	"strings"
)

func ParseCmdOutput(p []byte) string {
	o := parse(p)
	return strings.Join(o, "\r\n")
}

func ParseCmdInput(p []byte) string {
	o := parse(p)
	return strings.Join(o, "")
}

func GetPs1(p []byte) string {
	lines := parse(p)
	if len(lines) == 0 {
		return ""
	}
	return lines[len(lines)-1]
}

func parse(p []byte) []string {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error(r)
		}
	}()
	s := terminalparser.Screen{
		Rows:   make([]*terminalparser.Row, 0, 1024),
		Cursor: &terminalparser.Cursor{},
	}
	return s.Parse(p)
}
