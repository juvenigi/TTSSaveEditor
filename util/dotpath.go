package util

import "strings"

// LastSegment mimic filepath.Base
func LastSegment(s string) string {
	lastDot := strings.LastIndex(s, ".")
	if lastDot == -1 {
		return s
	} else {
		return s[lastDot+1:]
	}
}

// SegmentParent mimic filepath.Dir
func SegmentParent(s string, count int) string {
	parts := strings.Split(s, ".")
	if len(parts) < count {
		return s
	}
	return strings.Join(parts[:count], ".")
}
