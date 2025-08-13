package internal

import "strings"

type JsonPointer struct {
	sb strings.Builder
}

func (jp *JsonPointer) clone() *JsonPointer {
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
