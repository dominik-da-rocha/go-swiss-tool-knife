package databox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Dialect_Qt(t *testing.T) {
	d := NewDialect()
	assert.Equal(t, `"test"`, d.Qt("test"), "Qt as \"\" to string")
}

func Test_Dialect_QtAll(t *testing.T) {
	d := NewDialect()
	list := []string{"test1", "test2", "test3", "test4"}

	assert.Equal(t,
		[]string{`"test1"`, `"test2"`, `"test3"`, `"test4"`},
		d.QtAll(list...),
		"QtAll works like Qt with slices")

}
