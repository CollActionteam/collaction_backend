package main

import (
	"testing"
)

func TestHandler(t *testing.T) {
	t.Run("check sender", func(t *testing.T) {
		want := "hello@collaction.org"

		got := getSender()
		if want != got {
			t.Fatal("expected:", want, "got:", got)
		}
	})

	t.Run("check recipient", func(t *testing.T) {
		want := "hello@collaction.org"

		got := getSender()
		if want != got {
			t.Fatal("expected:", want, "got:", got)
		}
	})
}
