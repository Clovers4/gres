package container

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCList(t *testing.T) {
	cls := NewCList()
	assert.NotNil(t, cls)
}

func TestCList_LPush(t *testing.T) {
	cls := NewCList()
	cls.LPush("a")
	assert.Equal(t, "[a]", fmt.Sprintf("%v", cls.Range(0, -1)))
	cls.LPush("b", "c", "d")
	assert.Equal(t, "[d c b a]", fmt.Sprintf("%v", cls.Range(0, -1)))
}

func TestCList_LPop(t *testing.T) {
	cls := NewCList()
	cls.LPush("b", "c", "d")
	cls.LPop()
	assert.Equal(t, "[c b]", fmt.Sprintf("%v", cls.Range(0, -1)))
}

func TestCList_RPush(t *testing.T) {
	cls := NewCList()
	cls.RPush("a")
	assert.Equal(t, "[a]", fmt.Sprintf("%v", cls.Range(0, -1)))
	cls.RPush("b", "c", "d")
	assert.Equal(t, "[a b c d]", fmt.Sprintf("%v", cls.Range(0, -1)))
}

func TestCList_RPop(t *testing.T) {
	cls := NewCList()
	cls.RPush("b", "c", "d")
	cls.RPop()
	assert.Equal(t, "[b c]", fmt.Sprintf("%v", cls.Range(0, -1)))
}

func TestCList_Range(t *testing.T) {
	cls := NewCList()
	cls.LPush("a", "b", "c")
	assert.Equal(t, "[c b a]", fmt.Sprintf("%v", cls.Range(0, -1)))
}
