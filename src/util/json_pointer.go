package util

import "strings"

type JsonPointer struct {
	sb strings.Builder
}

func (jp *JsonPointer) Clone() *JsonPointer {
	var out JsonPointer
	s := jp.sb.String()
	out.sb.Grow(len(s))
	out.sb.WriteString(s)
	return &out
}

func (jp *JsonPointer) Grow(size int) {
	jp.sb.Grow(size)
}
func (jp *JsonPointer) Reset() {
	jp.sb.Reset()
}

func (jp *JsonPointer) Append(part string) {
	jp.sb.WriteString("/")
	jp.sb.WriteString(normalize(part))
}

func (jp *JsonPointer) BuildPointer() string {
	return jp.sb.String()
}

func normalize(str string) string {
	str = strings.Replace(str, "/", "~1", -1)

	return strings.Replace(str, "~", "~0", -1)
}

// and this one works
func (jp *JsonPointer) cloneFrom(src *JsonPointer) {
	jp.sb.Reset()
	jp.sb.Grow(src.sb.Len())
	jp.sb.WriteString(src.sb.String())
}
