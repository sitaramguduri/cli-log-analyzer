package golden

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI_Golden(t *testing.T) {
	root := filepath.FromSlash("../")
	logf := filepath.Join(root, "testdata", "sample.log")

	cmd := exec.Command("go", "run", "./cmd/analyze", "-file", logf, "-q", "0.5,0.95", "-top", "2")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Dir = root

	if err := cmd.Run(); err != nil {
		t.Fatalf("run: %v\n%s", err, out.String())
	}

	got := normalize(out.String())
	wantB, err := os.ReadFile(filepath.Join(root, "golden", "summary.txt"))
	if err != nil {
		t.Fatal(err)
	}
	want := normalize(string(wantB))

	if got != want {
		t.Fatalf("golden mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func normalize(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return strings.Join(lines, "\n")
}
