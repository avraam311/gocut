package integrationtests

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const binName = "gocut"

func buildBinary(t *testing.T) string {
	t.Helper()
	binPath := filepath.Join(t.TempDir(), binName)
	cmd := exec.Command("go", "build", "-o", binPath, "../cmd/gocut")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, string(out))
	}
	return binPath
}

func runCmd(t *testing.T, bin string, args ...string) string {
	t.Helper()
	cmd := exec.Command(bin, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out.String())
	}
	return strings.TrimRight(out.String(), "\n")
}

func generateFields(fieldsPerLine int) []string {
	fields := make([]string, fieldsPerLine)
	for i := 1; i <= fieldsPerLine; i++ {
		fields[i-1] = fmt.Sprintf("f%d", i)
	}
	return fields
}

func generateTestFile(t *testing.T, numLines, fieldsPerLine int, delim string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "testinput.txt")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create file: %v", err)
	}
	defer f.Close()

	for i := 0; i < numLines; i++ {
		fields := generateFields(fieldsPerLine)
		line := strings.Join(fields, delim) + "\n"
		if _, err := f.WriteString(line); err != nil {
			t.Fatalf("write line: %v", err)
		}
	}
	return path
}

func TestEmptyFile(t *testing.T) {
	t.Parallel()
	bin := buildBinary(t)
	input := filepath.Join(t.TempDir(), "empty.txt")
	os.Create(input) // ignore err for simplicity
	out := runCmd(t, bin, "-f", "1", "-d", " ", input)
	if out != "" {
		t.Errorf("expected empty output, got: %q", out)
	}
}

func TestSingleLine(t *testing.T) {
	t.Parallel()
	bin := buildBinary(t)
	input := generateTestFile(t, 1, 5, " ")
	out := runCmd(t, bin, "-f", "2,4", "-d", " ", input)
	expected := "f2 f4"
	if out != expected {
		t.Errorf("SingleLine got %q want %q", out, expected)
	}
}

func TestShortLinesNoS(t *testing.T) {
	t.Parallel()
	bin := buildBinary(t)
	input := generateTestFile(t, 2, 1, " ")
	out := runCmd(t, bin, "-f", "1,2", "-d", " ", input)
	expected := "f1 \nf1 "
	if out != expected {
		t.Errorf("ShortNoS got %q want %q", out, expected)
	}
}

func TestShortLinesWithS(t *testing.T) {
	t.Parallel()
	bin := buildBinary(t)
	input := generateTestFile(t, 2, 1, " ")
	out := runCmd(t, bin, "-f", "1,2", "-d", " ", "-s", input)
	if out != "" {
		t.Errorf("ShortWithS expected empty, got %q", out)
	}
}

func TestWideLine(t *testing.T) {
	t.Parallel()
	bin := buildBinary(t)
	input := generateTestFile(t, 1, 20, " ")
	out := runCmd(t, bin, "-f", "5,15", "-d", " ", input)
	expected := "f5 f15"
	if out != expected {
		t.Errorf("Wide got %q want %q", out, expected)
	}
}

func TestLargeFile(t *testing.T) {
	t.Parallel()
	bin := buildBinary(t)
	input := generateTestFile(t, 1000, 10, " ")
	out := runCmd(t, bin, "-f", "2,5", "-d", " ", input)
	lines := strings.Count(out, "\n") + 1
	expectedLines := 1000
	if lines != expectedLines {
		firstLine := strings.Split(out, "\n")[0]
		t.Errorf("LargeFile expected %d lines, got %d: first %q", expectedLines, lines, firstLine)
	}
}

func TestSimpleGrep(t *testing.T) {
	t.Parallel()
	bin := buildBinary(t)
	out := runCmd(t, bin, "-f", "1,3", "-d", " ", "./test_files/simple.txt")
	expected := readFile(t, "./expected_files/simple.txt")
	if out != expected {
		t.Errorf("unexpected output:\nGot:\n%s\nWant:\n%s", out, expected)
	}
}

func TestFlagS(t *testing.T) {
	t.Parallel()
	bin := buildBinary(t)
	out := runCmd(t, bin, "-f", "1-3", "-d", " ", "-s", "./test_files/flagS.txt")
	expected := readFile(t, "./expected_files/flagS.txt")
	if out != expected {
		t.Errorf("unexpected output:\nGot:\n%s\nWant:\n%s", out, expected)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}
	return string(data)
}
