package strings_test

import (
	"go-tool-box/strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RemoveFromStings_Found(t *testing.T) {
	list := []string{
		"test-1",
		"test-a",
		"test-v",
		"test-4",
		"test-d",
	}
	copy := strings.RemoveFromStings(list, "test-v")
	assert.Equal(t, 5, len(list))
	assert.Equal(t, 4, len(copy))
	assert.Equal(t, "test-1", copy[0])
	assert.Equal(t, "test-a", copy[1])
	assert.Equal(t, "test-4", copy[2])
	assert.Equal(t, "test-d", copy[3])
}

func Test_RemoveFromStings_NotFound(t *testing.T) {
	list := []string{
		"test-1",
		"test-a",
		"test-v",
		"test-4",
		"test-d",
	}
	copy := strings.RemoveFromStings(list, "unknown")
	assert.Equal(t, list, copy)
}

func Test_SortedStringsContains_Found(t *testing.T) {
	list := []string{
		"test-1",
		"test-a",
		"test-v",
		"test-4",
		"test-d",
	}
	found := strings.SortedStringsContains(list, "test-a")
	assert.Equal(t, true, found)
}

func Test_SortedStringsContains_NotFound(t *testing.T) {
	list := []string{
		"test-1",
		"test-a",
		"test-v",
		"test-4",
		"test-d",
	}
	found := strings.SortedStringsContains(list, "unknown")
	assert.Equal(t, false, found)
}

func Test_IndexOf_NotFound(t *testing.T) {
	list := []string{
		"test-1",
		"test-a",
		"test-v",
		"test-4",
		"test-d",
	}
	found := strings.IndexOfString(list, "unknown")
	assert.Equal(t, -1, found)
}

func Test_IndexOf_Found(t *testing.T) {
	list := []string{
		"test-1",
		"test-a",
		"test-v",
		"test-4",
		"test-d",
	}
	found := strings.IndexOfString(list, "test-v")
	assert.Equal(t, 2, found)
}

func Test_IsNilOrEmpty_Nil(t *testing.T) {
	assert.Equal(t, true, strings.IsNilOrEmpty(nil))
}

func Test_IsNilOrEmpty_Empty(t *testing.T) {
	empty := ""
	assert.Equal(t, true, strings.IsNilOrEmpty(&empty))
}

func Test_IsNilOrEmpty_Some(t *testing.T) {
	some := "some"
	assert.Equal(t, false, strings.IsNilOrEmpty(&some))
}
