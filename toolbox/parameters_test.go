package toolbox_test

import (
	"testing"

	"github.com/dominik-da-rocha/go-toolbox/toolbox"
	"github.com/stretchr/testify/assert"
)

func Test_IsAnyOfOrDefault_NotFound(t *testing.T) {
	allowed := []string{
		"test-1",
		"test-a",
		"test-v",
		"test-4",
		"test-d",
	}
	in := "hello"
	value := toolbox.ToAnyStringOf(&in, allowed, "world")
	assert.Equal(t, "world", value)
}

func Test_IsAnyOfOrDefault_Nil(t *testing.T) {
	allowed := []string{
		"test-1",
		"test-a",
		"test-v",
		"test-4",
		"test-d",
	}
	value := toolbox.ToAnyStringOf(nil, allowed, "world")
	assert.Equal(t, "world", value)
}

func Test_IsAnyOfOrDefault_Found(t *testing.T) {
	allowed := []string{
		"test-1",
		"test-a",
		"test-v",
		"test-4",
		"test-d",
	}
	in := "test-v"
	value := toolbox.ToAnyStringOf(&in, allowed, "world")
	assert.Equal(t, "test-v", value)
}

func Test_ToInt64_Value(t *testing.T) {
	val := int64(123)
	value := toolbox.ToInt64(&val, 0)
	assert.Equal(t, int64(123), value)
}

func Test_ToInt64_Nil(t *testing.T) {
	value := toolbox.ToInt64(nil, 123)
	assert.Equal(t, int64(123), value)
}

func Test_ToString_Value(t *testing.T) {
	val := "hello"
	value := toolbox.ToString(&val, "world")
	assert.Equal(t, "hello", value)
}

func Test_ToString_Nil(t *testing.T) {
	value := toolbox.ToString(nil, "world")
	assert.Equal(t, "world", value)
}
