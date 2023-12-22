package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestValue struct {
	val1   int `my-tag:"hello"`
	val2   int `my-tag:"world"`
	val3   int `my-tag:"!"`
	ignore int `my-tag:"-"`
}

func Test_GetStructNames(t *testing.T) {
	val := TestValue{
		val1: 1,
		val2: 2,
		val3: 3,
	}
	valPtr := &val
	names := GetStructTags(valPtr, "my-tag")
	assert.Equal(t, 3, len(names))
	assert.Equal(t, "hello", names[0])
	assert.Equal(t, "world", names[1])
	assert.Equal(t, "!", names[2])
}
