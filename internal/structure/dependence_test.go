package structure_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/bc0xy/structure/internal/structure"
)

func TestFindModel(t *testing.T) {
	moduleName, packageName, err := FindModule(".")
	assert.Nil(t, err, "获取module信息失败")
	assert.Equal(t, "github.com/bc0xy/structure", moduleName, "module名不对")
	assert.Equal(t, "github.com/bc0xy/structure/internal/structure", packageName, "package名不对")
}

func TestGetDepends(t *testing.T) {
	depend, err := GetDepends("/home/go/workspace/gitlab/DevOpsService/cmd/devops_server")
	assert.Nil(t, err, "获取依赖关系失败")
	assert.Equal(t, &DependPackage{}, depend, "获取的依赖关系有误")
}
