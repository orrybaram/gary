package ui

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// spinnerFrames is the braille spinner cycle.
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// spinnerInterval is how often the frame advances.
const spinnerInterval = 100 * time.Millisecond

// Spinner draws an animated "working" indicator on a single terminal line
// until Stop is called. Stop is synchronous: it guarantees the goroutine has
// finished writing and the line has been erased before it returns, so callers
// can print immediately afterwards without interleaving.
type Spinner struct {
	stop chan struct{}
	done chan struct{}
	once sync.Once
}

// StartSpinner begins animating "<frame> <label>… (Ns)" on stdout.
// If stdout is not a terminal (piped or redirected) it returns an inert
// spinner, so machine-readable output never gets escape codes in it.
func StartSpinner(label string) *Spinner {
	s := &Spinner{
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}

	if !isTerminal() {
		close(s.done)
		return s
	}

	go s.run(label)
	return s
}

func (s *Spinner) run(label string) {
	defer close(s.done)

	ticker := time.NewTicker(spinnerInterval)
	defer ticker.Stop()

	start := time.Now()
	frame := 0

	for {
		select {
		case <-s.stop:
			eraseLine()
			return
		case <-ticker.C:
			s.draw(label, frame, time.Since(start))
			frame++
		}
	}
}

// draw renders one frame in place, using \r to return to the line start.
func (s *Spinner) draw(label string, frame int, elapsed time.Duration) {
	glyph := spinnerFrames[frame%len(spinnerFrames)]
	fmt.Printf("\r%s%s %s… (%ds)%s\x1b[K",
		colorDim, glyph, label, int(elapsed.Seconds()), colorReset)
}

// Stop halts the animation and clears the line. It blocks until the drawing
// goroutine has exited, so no stray frame can land after subsequent output.
// Safe to call more than once.
func (s *Spinner) Stop() {
	s.once.Do(func() { close(s.stop) })
	<-s.done
}

// eraseLine returns the cursor to column 0 and clears to end of line.
func eraseLine() {
	fmt.Print("\r\x1b[K")
}

// isTerminal reports whether stdout is a character device, i.e. a real TTY
// rather than a pipe or file.
func isTerminal() bool {
	info, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
