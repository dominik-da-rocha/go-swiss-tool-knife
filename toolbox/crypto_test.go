package toolbox_test

import (
	"testing"

	"github.com/dominik-da-rocha/go-tool-box/toolbox"
	"github.com/stretchr/testify/assert"
)

func Test_Sha256_Hello(t *testing.T) {
	hash := toolbox.HashSha256([]byte("Hello"))
	assert.Equal(t, []byte{0x18, 0x5f, 0x8d, 0xb3, 0x22, 0x71, 0xfe, 0x25, 0xf5, 0x61, 0xa6, 0xfc, 0x93, 0x8b, 0x2e, 0x26, 0x43, 0x6, 0xec, 0x30, 0x4e, 0xda, 0x51, 0x80, 0x7, 0xd1, 0x76, 0x48, 0x26, 0x38, 0x19, 0x69}, hash)
}

func Test_NewRandomBytes_10(t *testing.T) {
	rand := toolbox.NewRandomBytes(10)
	assert.Equal(t, 10, len(rand))
}

func Test_NewRandomString_10(t *testing.T) {
	rand := toolbox.NewRandomString(1000)
	assert.Equal(t, 1000, len(rand))

}
