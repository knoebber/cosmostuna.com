package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	invalid         = "invalid"
	valid           = "ca"
	validFull       = "caliForniA"
	validFullSpaces = " caliForniA "
	validCaps       = "IA"
	validFullCaps   = "CALIFORNIA"
)

func TestShippable(t *testing.T) {
	t.Run("returns false", func(t *testing.T) {
		assert.False(t, shippable(nil))
		assert.False(t, shippable(&invalid))
	})

	truths := []*string{&valid, &validFull, &validCaps, &validFullCaps, &validFullSpaces}
	t.Run("returns true", func(t *testing.T) {
		for _, truth := range truths {
			assert.True(t, shippable(truth))
		}
	})
}
