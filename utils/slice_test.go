package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	t.Run("Remove at beginning", func(t *testing.T) {
		s := []string{"1", "2", "3", "4", "5"}
		Remove(&s, 0)
		assert.Equal(t, []string{"2", "3", "4", "5"}, s)
	})

	t.Run("Remove at end", func(t *testing.T) {
		s := []string{"1", "2", "3", "4", "5"}
		Remove(&s, 4)
		assert.Equal(t, []string{"1", "2", "3", "4"}, s)
	})

	t.Run("Remove in middle", func(t *testing.T) {
		s := []string{"1", "2", "3", "4", "5"}
		Remove(&s, 2)
		assert.Equal(t, []string{"1", "2", "4", "5"}, s)
	})
}
