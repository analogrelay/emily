package scanner

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmptyWindow(t *testing.T) {
	win := newWindow("")
	require.Equal(t, "", win.Rest())
	require.Equal(t, "", win.String())

	// Advance is a no-op at the end
	win.Advance()
	require.Equal(t, "", win.Rest())
	require.Equal(t, "", win.String())

	// Next/Peek return 0 and don't advance
	require.Equal(t, rune(0), win.Peek())
	require.Equal(t, rune(0), win.Next())
	require.Equal(t, "", win.Rest())
	require.Equal(t, "", win.String())
}

func TestAdvancingWindow(t *testing.T) {
	win := newWindow("abcdefg")

	require.Equal(t, 'a', win.Next())
	require.Equal(t, 'b', win.Next())
	require.Equal(t, 'c', win.Next())
	require.Equal(t, "abc", win.String())
	require.Equal(t, "defg", win.Rest())

	win.Advance()
	require.Equal(t, "", win.String())
	require.Equal(t, "defg", win.Rest())
	require.Equal(t, 'd', win.Next())
	require.Equal(t, 'e', win.Next())
	require.Equal(t, 'f', win.Next())
	require.Equal(t, "def", win.String())
	require.Equal(t, "g", win.Rest())

	win.Advance()
	require.Equal(t, "", win.String())
	require.Equal(t, "g", win.Rest())
	require.Equal(t, 'g', win.Next())
	require.Equal(t, rune(0), win.Next())
	require.Equal(t, "g", win.String())
	require.Equal(t, "", win.Rest())

	win.Advance()
	require.Equal(t, rune(0), win.Next())
	require.Equal(t, "", win.String())
	require.Equal(t, "", win.Rest())
}

func TestPeek(t *testing.T) {
	win := newWindow("abc")
	require.Equal(t, 'a', win.Peek())
	require.Equal(t, 'a', win.Peek())
	require.Equal(t, 0, win.Len())
	require.Equal(t, 3, win.Remaining())
}

func TestRewind(t *testing.T) {
	win := newWindow("abc")
	require.Equal(t, 'a', win.Next())
	require.Equal(t, 'b', win.Next())
	require.Equal(t, 'c', win.Next())
	win.Rewind(2)
	require.Equal(t, "a", win.String())
	require.Equal(t, 2, win.Remaining())
	require.Equal(t, 'b', win.Next())
	require.Equal(t, 'c', win.Next())
}

func TestUtf8(t *testing.T) {
	win := newWindow("10¢, 50€")
	win.Take(2)
	require.Equal(t, "10", win.String())
	require.Equal(t, '¢', win.Next())
	require.Equal(t, "10¢", win.String())
	require.Equal(t, ", 50€", win.Rest())
	require.Equal(t, 4, win.Len()) // Len is in bytes not runes!
	win.Rewind(1)
	require.Equal(t, "10", win.String())
	require.Equal(t, "¢, 50€", win.Rest())
	win.Take(5)
	require.Equal(t, "10¢, 50", win.String())
	require.Equal(t, '€', win.Next())
	require.Equal(t, "10¢, 50€", win.String())
	require.Equal(t, 11, win.Len())
	win.Rewind(2)
	require.Equal(t, "10¢, 5", win.String())
	require.Equal(t, "0€", win.Rest())
	require.Equal(t, 7, win.Len())
	win.Take(2)
	win.Rewind(8)
	require.Equal(t, "", win.String())
	require.Equal(t, "10¢, 50€", win.Rest())
}

func TestTakeWhile(t *testing.T) {
	win := newWindow("aaaaabbbb")
	win.TakeWhile(func(r rune) bool { return r == 'a' })
	require.Equal(t, "aaaaa", win.String())
	require.Equal(t, "bbbb", win.Rest())
	win.TakeWhile(func(r rune) bool { return r == 'b' })
	require.Equal(t, "aaaaabbbb", win.String())
	require.Equal(t, "", win.Rest())
}

func TestTakeIf(t *testing.T) {
	win := newWindow("abcde")
	require.True(t, win.TakeIf('a', 'e'))
	require.Equal(t, "a", win.String())
	require.Equal(t, "bcde", win.Rest())
	require.False(t, win.TakeIf('c', 'd'))
	require.Equal(t, "a", win.String())
	require.Equal(t, "bcde", win.Rest())
	require.True(t, win.TakeIf('e', 'b'))
	require.Equal(t, "ab", win.String())
	require.Equal(t, "cde", win.Rest())
}
