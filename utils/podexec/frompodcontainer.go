package podexec

import (
	"bytes"
	"fmt"
)

// FromPodContainer 从pod内拷贝到io.Writer
func (p *PodExec) FromPodContainer(dest []string, style string) error {
	switch style {
	case "tar":
		p.Command = append([]string{"tar", "cf", "-"}, dest...)
	case "zip":
		p.Command = append([]string{"/kftools", "zip"}, dest...)
	default:
		p.Command = append([]string{"tar", "cf", "-"}, dest...)
	}
	p.Tty = false
	var stderr bytes.Buffer
	p.Stderr = &stderr
	err := p.Exec()
	if err != nil {
		return fmt.Errorf(err.Error(), stderr)
	}
	if err != nil {
		if len(stderr.Bytes()) != 0 {
			return fmt.Errorf("STDERR: " + stderr.String())
		}
		return err
	}
	// 三次重试
	//attempts := 3
	//attempt := 0
	//for attempt < attempts {
	//	attempt++
	//
	//	stderr, err := c.Exec()
	//	logs.Error(err, string(stderr))
	//	if attempt == attempts {
	//		if err != nil {
	//			return err
	//		}
	//		if len(stderr) != 0 {
	//			return fmt.Errorf("STDERR: " + string(stderr))
	//		}
	//	}
	//	if err == nil {
	//		return nil
	//	}
	//	time.Sleep(time.Duration(attempt) * time.Second)
	//}
	return nil
}
