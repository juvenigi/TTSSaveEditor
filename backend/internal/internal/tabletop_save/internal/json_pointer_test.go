package internal

import (
	"testing"

	jsonpatch "github.com/evanphx/json-patch"
)

var sample = `{
  "foo": [
    "bar",
    "baz"
  ],
  "": 0,
  "a/b": 1,
  "c%d": 2,
  "e^f": 3,
  "g|h": 4,
  "i\\j": 5,
  "k\"l": 6,
  " ": 7,
  "m~n": 8
}`

var modified = `{
  "foo": [
    "quox",
    "baz"
  ],
  "": 1,
  "a/b": 2,
  "c%d": 3,
  "e^f": 4,
  "g|h": 5,
  "i\\j": 6,
  "k\"l": 7,
  " ": 8,
  "m~n": 9
}`

var patch = `[
  { "op": "replace", "path": "/foo/0",  "value": "quox" },
  { "op": "replace", "path": "/",       "value": 1 },
  { "op": "replace", "path": "/a~1b",   "value": 2 },
  { "op": "replace", "path": "/c%d",    "value": 3 },
  { "op": "replace", "path": "/e^f",    "value": 4 },
  { "op": "replace", "path": "/g|h",    "value": 5 },
  { "op": "replace", "path": "/i\\j",   "value": 6 },
  { "op": "replace", "path": "/k\"l",   "value": 7 },
  { "op": "replace", "path": "/ ",      "value": 8 },
  { "op": "replace", "path": "/m~0n",   "value": 9 }
]`

func TestJsonMergePatch(t *testing.T) {
	t.Logf("\n%s", sample)
	patch, err := jsonpatch.CreateMergePatch([]byte(sample), []byte(modified))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%s", patch)
}

func TestJsonPatch(t *testing.T) {
	t.Logf("\n%s", sample)
	patches, err := jsonpatch.DecodePatch([]byte(patch))
	if err != nil {
		t.Fatal(err)
	}
	patched, err := patches.ApplyIndent([]byte(sample), "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%s", patched)
}

func TestJsonPointer_Clone(t *testing.T) {
	var pointer JsonPointer
	pointer.Append("foo")
	t.Logf("\n%s", pointer.BuildPointer())

	var clone = pointer.clone()
	clone.Append("bar")
	t.Logf("\n%s", clone.BuildPointer())

}

func TestJsonPointer_CloneFrom(t *testing.T) {
	var src JsonPointer
	var target JsonPointer

	target.cloneFrom(&src)
	target.Append("bar")
	src.Append("quox")
	t.Logf("\n%s", target.BuildPointer())
	t.Logf("\n%s", src.BuildPointer())
}
