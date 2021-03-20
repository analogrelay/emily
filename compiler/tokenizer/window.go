package tokenizer

import (
	"unicode/utf8"
)

type window struct {
	str   string
	start int
	end   int
}

func newWindow(input string) window {
	return window{input, 0, 0}
}

// Len returns the length of the string in the window.
func (w *window) Len() int {
	return w.end - w.start
}

// Remaining returns the length in bytes of the remaining string after the end of the window
func (w *window) Remaining() int {
	return len(w.str) - w.end
}

// Rest returns the string that remains beyond the end of the window
func (w *window) Rest() string {
	return w.str[w.end:]
}

// String returns the string covered by the window
func (w *window) String() string {
	return w.str[w.start:w.end]
}

// Advance moves the window immediately beyond the current end and resets the length to 0
//
// Returns the start and end of the window that was just advanced past
func (w *window) Advance() (int, int) {
	s, e := w.start, w.end
	w.start = w.end
	return s, e
}

// Next advances the window forward to include the next rune in the string, and returns it.
// If there are no further runes, the window is not advanced and 0 is returned.
func (w *window) Next() rune {
	if len(w.str) <= w.end {
		return 0
	}

	r, size := utf8.DecodeRuneInString(w.str[w.end:])
	w.end += size
	return r
}

// Take takes the next `n` runes into the window
func (w *window) Take(n int) {
	for i := 0; i < n; i++ {
		_ = w.Next()
	}
}

// TakeIf takes a rune into the window if it matches any of the provided rune and returns a boolean indicating if the rune was taken
func (w *window) TakeIf(rs ...rune) bool {
	p := w.Peek()
	for _, r := range rs {
		if p == r {
			w.Take(1)
			return true
		}
	}
	return false
}

// TakeWhile takes runes into the window until either the end is reached, or the `fn` provided returns `false`
func (w *window) TakeWhile(fn func(rune) bool) {
	for {
		r := w.Peek()
		if r == 0 || !fn(r) {
			return
		}
		w.Next()
	}
}

// Peek returns the rune that would be added to the window next if `Next` is called.
func (w *window) Peek() rune {
	if len(w.str) <= w.end {
		return 0
	}

	r, _ := utf8.DecodeRuneInString(w.str[w.end:])
	return r
}

// Rewind moves the end of the window back by `n` runes.
func (w *window) Rewind(n int) {
	for i := 0; i < n; i++ {
		for !utf8.RuneStart(w.str[w.end-1]) {
			w.end--
		}

		// Move one byte before the start of the next char
		w.end--
	}
}
