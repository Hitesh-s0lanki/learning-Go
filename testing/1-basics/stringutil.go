// Package basics shows the anatomy of a Go test: the code under test lives in
// stringutil.go and its tests live next to it in stringutil_test.go.
package basics

// Reverse returns s with its characters in reverse order. It works on runes,
// not bytes, so multi-byte characters (é, 世, emoji) are preserved.
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
