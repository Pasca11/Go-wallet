package types

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAccount(t *testing.T) {
	acc, err := NewAccount("Amir", "Sha256", "Testovich", "pass123")
	assert.Nil(t, err)

	fmt.Printf("%+v\n", acc)
}
