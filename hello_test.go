package main

import "testing"

func TestHello(t *testing.T) {
	t.Run("hello to som", func(t *testing.T) {
	

	got := Hello("Chris")
	want := "Hello, Chris"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
})
	t.Run("say to empty string", func(t *testing.T) {
		got := Hello("")
		want := "Hello, "

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	func assertCorrectMessage(t *testing.T, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

	}
}