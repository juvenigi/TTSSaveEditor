package properties

import (
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func newYAMLFixture(content string) io.Reader {
	return strings.NewReader(content)
}

func TestParseProperties_ValidYAMLConfigFile(t *testing.T) {
	properties, err2 := writeDefaultProperties()
	if err2 != nil {
		t.Fatal(err2)
	}
	file, err := os.Open(properties)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := newProperties(file)
	if err != nil {
		t.Fatal(err)
	}

	if len(strings.TrimSpace(actual.GameDir)) == 0 {
		t.Fatal("game dir is empty")
	}
	if len(strings.TrimSpace(actual.PackDataDir)) == 0 {
		t.Fatal("pack data dir is empty")
	}
}

func TestParseProperties_InvalidYAML(t *testing.T) {
	badYAML := `: invalid yaml`
	file := newYAMLFixture(badYAML)

	_, err := newProperties(file)
	if err == nil {
		t.Fatal("Expected error for invalid YAML, got nil")
	}
}

func TestParseProperties_MissingPropertyKeys(t *testing.T) {
	file := newYAMLFixture("foo: true\n")

	_, err := newProperties(file)
	if err == nil {
		t.Fatal("Expected error for missing property keys, got nil")
	}
}

// todo: dont do this without flags
func TestGetConfigFilePath(t *testing.T) {
	origArgs := os.Args
	t.Log(origArgs)
	defer func() { os.Args = origArgs }()

	cfgAbs, err := filepath.Abs("./../../../resources/config")
	if err != nil {
		t.Fatal(err)
	}
	os.Args = []string{"", cfgAbs}

	path := getCfgFilePath()
	if path == "" {
		t.Fatal("Expected GetConfigFilePath, got nothing")
	}
	t.Log(path)
}

// todo: verify after fixing
func TestNewPropertiesController(t *testing.T) {
	controller := NewPropertiesController()
	properties := controller.GetApplicationProperties()
	t.Log(properties)
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	executableDir := filepath.Dir(executable)
	files, err := os.ReadDir(executableDir)
	if err != nil {
		t.Fatal(err)
	}
	cfgFilenames := []string{ConfigFileYaml}
	found := false
	for _, file := range files {
		if slices.ContainsFunc(cfgFilenames, func(s string) bool {
			return strings.HasSuffix(s, file.Name())
		}) {
			t.Logf("file found: %v", filepath.Join(executableDir, file.Name()))
			found = true
			continue
		}
	}
	if !found {
		t.Fatal("Could not find config files")
	}

	if len(properties.GameDir) == 0 {
		t.Fatal("GameDir is undefined")
	}
}

func TestPathSeparator(t *testing.T) {
	given := "C://foo/bar"
	actual := filepath.FromSlash(given)
	if strings.Compare(actual, "C:\\\\foo\\bar") != 0 {
		t.Fatal("Not what expected:", actual)
	}
}

func TestPathSeparatorIdentity(t *testing.T) {
	given := "C:\\\\foo\\bar"
	actual := filepath.FromSlash(given)
	if strings.Compare(actual, given) != 0 {
		t.Fatal("Expected no change, got:", actual)
	}

}
