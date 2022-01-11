package symwalk

import (
	"os"
	"path/filepath"
)

// symwalkFunc 为常规文件调用提供的 WalkFn。
// 但是，当它遇到符号链接时，它会使用
// filepath.EvalSymlinks 函数并在解析路径上递归调用 symwalk.Walk。
// 这样可以确保 unlink filepath.Walk，遍历不会在符号链接处停止。
//
// 请注意，如果有任何非终止循环，symwalk.Walk 不会终止
// 文件结构。
func walk(filename string, linkDirname string, walkFn filepath.WalkFunc) error {
	symWalkFunc := func(path string, info os.FileInfo, err error) error {
		if fname, err := filepath.Rel(filename, path); err == nil {
			path = filepath.Join(linkDirname, fname)
		} else {
			return err
		}

		if err == nil && info.Mode()&os.ModeSymlink == os.ModeSymlink {
			finalPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}
			info, err := os.Lstat(finalPath)
			if err != nil {
				return walkFn(path, info, err)
			}
			if info.IsDir() {
				return walk(finalPath, path, walkFn)
			}
		}
		return walkFn(path, info, err)
	}
	return filepath.Walk(filename, symWalkFunc)
}

// Walk 扩展 filepath.Walk 也遵循符号链接
func Walk(path string, walkFn filepath.WalkFunc) error {
	return walk(path, path, walkFn)
}
