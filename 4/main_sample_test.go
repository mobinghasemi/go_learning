package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	b := NewBoard(10, 8)
	a1 := NewPaper(3, 2, 101)
	a1.addTo(b, 5, 4)

	a2 := NewPaper(3, 2, 102)
	a2.addTo(b, 5, 5)

	ids := b.getAnnouncementIDsAt(5, 4)
	assert.Equal(t, []int{101, 102}, ids)

	ids = a1.removeAndGetIDsOnTop()
	assert.Equal(t, []int{102}, ids)

	ids = b.getAnnouncementIDsAt(5, 4)
	assert.Equal(t, []int{102}, ids)
}
