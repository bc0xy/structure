package structure_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/bc0xy/structure/internal/structure"
)

// TestGetImportPackages 测试imports.go的依赖库
func TestGetImportPackages(t *testing.T) {
	got, err := GetImportPackages("imports.go")
	assert.Nil(t, err, "获取当前目录下的依赖包失败")
	assert.Equal(t, []string{
		"os", "regexp", "strings",
		"github.com/pkg/errors",
		"github.com/samber/lo",
	}, got, "获取go文件失败")
}
