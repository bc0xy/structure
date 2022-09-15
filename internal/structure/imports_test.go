package structure_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/bc0xy/structure/internal/structure"
)

func TestGetImportPackages(t *testing.T) {
	goFiles, err := GetGoFiles("")
	assert.Nil(t, err, "获取当前目录go文件失败")
	got, err := GetImportPackages(goFiles...)
	assert.Nil(t, err, "获取当前目录下的依赖包失败")
	assert.Equal(t, goFiles, got, "获取go文件失败")
}
