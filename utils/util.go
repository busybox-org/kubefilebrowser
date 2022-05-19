package utils

import (
	"bufio"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func InSliceString(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func FileOrPathExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func EnsureDirExist(name string) error {
	if !FileOrPathExist(name) {
		return os.MkdirAll(name, os.ModePerm)
	}
	return nil
}

func GzipCompressFile(srcPath, dstPath string) error {
	sf, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer df.Close()
	writer := gzip.NewWriter(df)
	writer.Name = filepath.Base(srcPath)
	writer.ModTime = time.Now().UTC()
	_, err = io.Copy(writer, sf)
	if err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return nil
}

func File2lines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LinesFromReader(f)
}

func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// InsertStringToFile inserts a string to a file
/**
 * Insert sting to n-th line of file.
 * If you want to insert a line, append newline '\n' to the end of the string.
 */
func InsertStringToFile(path, str string, index int) error {
	lines, err := File2lines(path)
	if err != nil {
		return err
	}

	fileContent := ""
	for i, line := range lines {
		if i == index {
			fileContent += str
		}
		fileContent += line
		fileContent += "\n"
	}

	return ioutil.WriteFile(path, []byte(fileContent), 0644)
}

// ToLinuxPath windows 路径转Linux格式路径
func ToLinuxPath(path string) string {
	if runtime.GOOS == "windows" {
		if len(strings.Split(path, "C:")) != 0 {
			path = strings.ReplaceAll(path, "C:", "")
		}
		if len(strings.Split(path, "c:")) != 0 {
			path = strings.ReplaceAll(path, "c:", "")
		}
		path = strings.ReplaceAll(path, `\\`, "/")
		path = strings.ReplaceAll(path, `\`, "/")
	}
	return path
}
