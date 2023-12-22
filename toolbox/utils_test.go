package toolbox_test

import (
	"errors"
	"go-tool-box/toolbox"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Uups_Panics(t *testing.T) {
	assert.Panics(t, func() {
		toolbox.Uups(errors.New("test"))
	})
}

func Test_MustBeTrue_Panics(t *testing.T) {
	assert.Panics(t, func() {
		toolbox.MustBeTrue(false, "doom to fail")
	})
}
