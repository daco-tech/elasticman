package general

// HasPrefix function checks if the prefix is contained in the string.
// Set verbose true if you want more output details.
func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}
