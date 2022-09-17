package structure

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

// DependPackage 包的依赖关系
type DependPackage struct {
	// Name 包名
	Name string
	// Depends 依赖包
	Depends []*DependPackage
}

// FindModule 查找module信息
// HACK: 多次依赖一个包,这个包会递归多次,影响性能,可以优化
func FindModule(dir string) (moduleName, packageName string, err error) {
	modPath := strings.Split(dir, "/")
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

		// HACK: 拆解为更小的函数
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

		packageName = moduleName + dir[len(modDir):]
		return
	}

	return "", "", errors.Errorf("在目录(%s)的各层目录中未找到go.mod", dir)
}

// GetDepends 获取依赖关系
func GetDepends(dir string) (_ *DependPackage, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "获取依赖关系失败")
		}
	}()

	absDir, err := filepath.Abs(dir)
	if err != nil {
		err = errors.Wrapf(err, "目录(%s)转换成绝对路径失败", dir)
		return
	}

	gofiles, err := GetGoFiles(absDir)
	if err != nil {
		return nil, err
	}

	// HACK: 递归之前,可以将目录先准备好,避免每次迭代都会找go.mod
	importPackages, err := GetImportPackages(gofiles...)
	if err != nil {
		return nil, err
	}

	moduleName, packageName, err := FindModule(absDir)
	if err != nil {
		return nil, err
	}

	depends := &DependPackage{
		Name:    moduleName,
		Depends: make([]*DependPackage, 0),
	}
	depends.Depends = lo.Map(importPackages, func(x string, _ int) *DependPackage {
		return &DependPackage{Name: x, Depends: make([]*DependPackage, 0)}
	})

	for _, depend := range depends.Depends {
		if strings.HasPrefix(depend.Name, moduleName) {
			dependPackagePath := filepath.Join(
				strings.TrimSuffix(absDir, strings.TrimPrefix(packageName, moduleName)),
				strings.TrimPrefix(depend.Name, moduleName),
			)
			importDepend, err := GetDepends(dependPackagePath)
			if err != nil {
				return nil, errors.WithMessagef(err, "获取子依赖(%s)的依赖关系失败", dependPackagePath)
			}
			depend.Depends = importDepend.Depends
		}
	}

	return depends, nil
}
