package main

import "strings"

// CleanPath is a utility function to clean up path by eliminating "..", ".",
// and prefixed "/".
func CleanPath(path string) string {
	parts := make([]string, 0, 10)
	for _, n := range strings.Split(path, "/") {
		switch n {
		case "": // ignore
		case ".": // ignore
		case "..":
			parts = parts[:len(parts)-1] // pop
		default:
			parts = append(parts, n) // push
		}
	}
	return "/" + strings.Join(parts, "/")
}

// JoinPath joins arguments with separator "/" and returns it witn cleaning-up.
func JoinPath(path ...string) string {
	return CleanPath(strings.Join(path, "/"))
}
