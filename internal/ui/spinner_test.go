package ui

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

// captureStdout swaps os.Stdout for a pipe and returns what fn wrote.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stdout
	os.Stdout = w

	done := make(chan string)
	go func() {
		b, _ := io.ReadAll(r)
		done <- string(b)
	}()

	fn()

	os.Stdout = orig
	w.Close()
	return <-done
}

func TestDrawRendersFrameAndErase(t *testing.T) {
	out := captureStdout(t, func() {
		s := &Spinner{}
		s.draw("Thinking", 0, 3*time.Second)
		eraseLine()
	})

	if !strings.Contains(out, "⠋ Thinking… (3s)") {
		t.Errorf("missing frame text, got %q", out)
	}
	if !strings.HasPrefix(out, "\r") {
		t.Errorf("frame must start with CR, got %q", out)
	}
	if !strings.HasSuffix(out, "\r\x1b[K") {
		t.Errorf("must end by erasing the line, got %q", out)
	}
}

// TestStopIsSynchronous is the regression test for the interleaving bug: once
// Stop returns, no spinner frame may ever reach stdout again.
func TestStopIsSynchronous(t *testing.T) {
	out := captureStdout(t, func() {
		s := &Spinner{
			stop: make(chan struct{}),
			done: make(chan struct{}),
		}
		go s.run("Thinking")

		time.Sleep(250 * time.Millisecond) // let a few frames draw
		s.Stop()

		os.Stdout.WriteString("SENTINEL")
		time.Sleep(300 * time.Millisecond) // any late frame would land here
	})

	if !strings.Contains(out, "Thinking") {
		t.Fatalf("expected frames before stop, got %q", out)
	}
	if !strings.HasSuffix(out, "SENTINEL") {
		t.Errorf("output written after Stop returned: %q", out[strings.Index(out, "SENTINEL"):])
	}
}

// TestStopWithoutTTY covers the inert path: no goroutine, no output, and Stop
// must not block or panic when called repeatedly.
func TestStopWithoutTTY(t *testing.T) {
	out := captureStdout(t, func() {
		s := StartSpinner("Thinking")
		s.Stop()
		s.Stop()
	})
	if out != "" {
		t.Errorf("expected no output when stdout is not a tty, got %q", out)
	}
}
