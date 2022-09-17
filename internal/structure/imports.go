package structure

import (
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

// GetImportPackages 获取go文件中依赖的库
func GetImportPackages(goFiles ...string) ([]string, error) {
	if len(goFiles) == 0 {
		return []string{}, nil
	}

	packages := make([]string, 0, 20)
	for _, goFile := range goFiles {
		content, err := os.ReadFile(goFile)
		if err != nil {
			return nil, errors.Wrapf(err, "读%s文件内容失败", goFile)
		}
		re := regexp.MustCompile(`import[\n| |\t]*(\((?sU:.*)\)|".*"$)`)
		importString := re.FindString(string(content))
		if len(importString) == 0 {
			continue
		}

		importString = strings.ReplaceAll(importString, "import", "")
		importString = strings.ReplaceAll(importString, "(", "")
		importString = strings.ReplaceAll(importString, ")", "")
		pkgs := strings.Split(importString, "\n")
		packages = append(packages, lo.FilterMap(pkgs, func(pkg string, _ int) (string, bool) {
			pkg = strings.TrimSpace(pkg)

			if len(pkg) == 0 || strings.HasPrefix(pkg, "//") {
				return "", false
			}
			start := strings.Index(pkg, "\"")
			if start == -1 {
				return "", false
			}
			stop := strings.Index(string(pkg[start+1:]), "\"")
			if start == -1 {
				return "", false
			}

			return string(pkg[start+1 : stop+1]), true
		})...)
	}

	return lo.Uniq(packages), nil
}
