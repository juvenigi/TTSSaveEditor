package internal

import "strings"

type JsonPointer struct {
	sb strings.Builder
}

func (p *JsonPointer) Grow(size int) {
	p.sb.Grow(size)
}
func (p *JsonPointer) Reset() {
	p.sb.Reset()
}

func (jp *JsonPointer) Append(part string) {
	jp.sb.WriteString(normalize(part))
}

func (jp *JsonPointer) BuildPointer() string {
	jsonString := jp.sb.String()
	if len(jsonString) == 0 {
		return ""
	} else {
		return "/" + jsonString
	}
}

func normalize(str string) string {
	str = strings.Replace(str, "/", "~1", -1)

	return strings.Replace(str, "~", "~0", -1)
}
