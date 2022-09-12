package structure

import (
	"io/fs"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

// GetGoFiles 读取目录下的go文件
func GetGoFiles(dir string) ([]string, error) {
	if len(dir) == 0 {
		dir = "."
	}

	dirEntry, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "读取目录(%s)下的文件失败", dir)
	}
	return lo.FilterMap(dirEntry, func(x fs.DirEntry, _ int) (string, bool) {
		if !x.IsDir() && strings.HasSuffix(x.Name(), ".go") && !strings.HasSuffix(x.Name(), "_test.go") {
			return x.Name(), true
		}
		return "", false
	}), nil
}
