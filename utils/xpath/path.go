package xpath

import "os"

// IsDir 是否是目录
func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
