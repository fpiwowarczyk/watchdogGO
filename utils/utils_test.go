package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqualsStringsSlice_Success(t *testing.T) {

	output := Equals([]string{"a", "b"}, []string{"a", "b"})

	assert.True(t, output, "Slices with the same elements are equal")
}

func TestEqualsStringsSlice_Fail_DiffElems(t *testing.T) {

	output := Equals([]string{"a", "b"}, []string{"a", "c"})

	assert.False(t, output, "Slices with different elements aren't equal")
}

func TestEqualsStringsSlice_Fail_DiffLen(t *testing.T) {

	output := Equals([]string{"a", "b"}, []string{"a"})

	assert.False(t, output, "Slices with different len arent equal")
}
