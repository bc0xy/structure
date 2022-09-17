package structure

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// DependPackage 包的依赖关系
type DependPackage struct {
	// Name 包名
	Name string
	// Depends 依赖包
	Depends []DependPackage
}

// FindModule 查找module信息
func FindModule(dir string) (moduleName, packageName string, err error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		err = errors.Wrapf(err, "目录(%s)转换成绝对路径失败", dir)
		return
	}

	modPath := strings.Split(absDir, "/")
	var modFile string
	for i := len(modPath); i > 0; i-- {
		modDir := "/" + filepath.Join(modPath[:i]...)
		modFile = filepath.Join(modDir, "go.mod")
		var f *os.File
		f, err = os.Open(modFile)
		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			return "", "", errors.Wrapf(err, "打开文件(%s)失败", modFile)
		}
		defer f.Close()

		firstLine := make([]byte, 256)
		_, err = f.Read(firstLine)
		if err != nil {
			return "", "", errors.Wrapf(err, "读文件(%s)失败", modFile)
		}

		if len(firstLine) <= len("module \n") {
			return "", "", errors.Errorf("从go.mod中解析的首行长度太短(%s)", string(firstLine))
		}

		moduleName = string(firstLine)[len("module "):]
		newLine := strings.Index(moduleName, "\n")
		if newLine == -1 {
			return "", "", errors.Errorf("从go.mod中解析的首行长度太长(%s)", string(firstLine))
		}
		moduleName = moduleName[:newLine]

		packageName = moduleName + absDir[len(modDir):]
		return
	}

	return "", "", errors.Errorf("在目录(%s)的各层目录中未找到go.mod", dir)
}
