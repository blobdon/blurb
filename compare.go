package main

import "strings"

// is containsAuthor reports whether the author's name is in the text
// This is currently the most naive version - need to test others
func containsAuthor(text, author string) (bool, error) {
	// func containsAuthor(text []byte, author string) (bool, error) {
	return strings.Contains(text, author), nil
	// return bytes.Contains(text, []byte(author)), nil
}

// Edge cases so far - naive Contains misses these
// All caps 			DAVID BROOKS vs David Brooks
// solution: Case-insensitive (all lower/upper)
// Missing part 		Nassim Taleb vs Nassim Nicholas Taleb
// solution: Some sort of name token-based search
// Spacing	 			Malcolm  Gladwell (extra space) vs Malcolm Gladwell
// solution: strip extra or all white space
