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
