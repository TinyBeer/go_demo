package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
	as := assert.New(t)

	// assert equality
	as.Equal(123, 123, "they should be equal")

	// assert inequality
	as.NotEqual(123, 456, "they should be not equal")

	// assert nil object
	var e error
	as.Nil(e, "nil error")

	// assert for not nil
	e = io.EOF
	if as.NotNil(e) {
		as.Equal(e.Error(), "EOF")
	}

}
