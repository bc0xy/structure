package structure_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/bc0xy/structure/internal/structure"
)

func TestGetGoFiles(t *testing.T) {
	dir, goFiles, err := getTestDirAndGoFiles(t)
	assert.Nil(t, err, "准备测试目录失败")
	got, err := GetGoFiles(dir)
	assert.Nil(t, err, "获取目录(%s)下的go文件失败", dir)
	assert.Equal(t, goFiles, got, "获取go文件失败")
}

// getTestDirAndGoFiles 获取测试目录和go文件列表
func getTestDirAndGoFiles(t *testing.T) (string, []string, error) {
	t.Helper()
	dir, err := os.MkdirTemp("", "structure*")
	assert.Nil(t, err, "创建测试目录失败")

	fileA := filepath.Join(dir, "a.go")
	a, err := os.Create(fileA)
	assert.Nil(t, err, "创建a.go失败")
	defer a.Close()

	fileB := filepath.Join(dir, "b.go")
	b, err := os.Create(fileB)
	assert.Nil(t, err, "创建b.go失败")
	defer b.Close()

	return dir, []string{fileA, fileB}, nil
}

func TestGetGoFilesInCurrentDir(t *testing.T) {
	got, err := GetGoFiles("")
	assert.Nil(t, err, "获取当前目录go文件失败")
	want := []string{"dependence.go", "files.go", "imports.go", "ui.go"}
	assert.Equal(t, want, got, "获取的go文件列表不正确")
}
