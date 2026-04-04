package input

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParsePRNumber_Valid(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"1", 1},
		{"123", 123},
		{"99999", 99999},
	}
	for _, tt := range tests {
		got, err := ParsePRNumber(tt.input)
		if err != nil {
			t.Errorf("ParsePRNumber(%q) error = %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ParsePRNumber(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestParsePRNumber_Invalid(t *testing.T) {
	tests := []string{
		"",
		"abc",
		"-1",
		"0",
		"1.5",
		"123abc",
	}
	for _, input := range tests {
		_, err := ParsePRNumber(input)
		if err == nil {
			t.Errorf("ParsePRNumber(%q) expected error", input)
		}
	}
}

func TestReadBody_Text(t *testing.T) {
	body, err := ReadBody(BodySource{Text: "hello world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body != "hello world" {
		t.Errorf("body = %q, want %q", body, "hello world")
	}
}

func TestReadBody_File(t *testing.T) {
	// Create temp file
	dir := t.TempDir()
	path := filepath.Join(dir, "body.txt")
	if err := os.WriteFile(path, []byte("file content"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	body, err := ReadBody(BodySource{File: path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body != "file content" {
		t.Errorf("body = %q, want %q", body, "file content")
	}
}

func TestReadBody_FileNotFound(t *testing.T) {
	_, err := ReadBody(BodySource{File: "/nonexistent/path/file.txt"})
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestReadBody_Stdin(t *testing.T) {
	stdin := strings.NewReader("stdin content")
	body, err := ReadBody(BodySource{
		Stdin:       stdin,
		StdinIsPipe: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body != "stdin content" {
		t.Errorf("body = %q, want %q", body, "stdin content")
	}
}

func TestReadBody_StdinNotPipe(t *testing.T) {
	stdin := strings.NewReader("stdin content")
	_, err := ReadBody(BodySource{
		Stdin:       stdin,
		StdinIsPipe: false, // Not a pipe, should be ignored
	})
	if err == nil {
		t.Fatal("expected error when stdin is not a pipe")
	}
}

func TestReadBody_TextOverridesFile(t *testing.T) {
	// Create temp file
	dir := t.TempDir()
	path := filepath.Join(dir, "body.txt")
	if err := os.WriteFile(path, []byte("file content"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	body, err := ReadBody(BodySource{
		Text: "text content",
		File: path,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body != "text content" {
		t.Errorf("body = %q, want %q (text should override file)", body, "text content")
	}
}

func TestReadBody_FileOverridesStdin(t *testing.T) {
	// Create temp file
	dir := t.TempDir()
	path := filepath.Join(dir, "body.txt")
	if err := os.WriteFile(path, []byte("file content"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	stdin := strings.NewReader("stdin content")
	body, err := ReadBody(BodySource{
		File:        path,
		Stdin:       stdin,
		StdinIsPipe: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body != "file content" {
		t.Errorf("body = %q, want %q (file should override stdin)", body, "file content")
	}
}

func TestReadBody_NoSource(t *testing.T) {
	_, err := ReadBody(BodySource{})
	if err == nil {
		t.Fatal("expected error when no source provided")
	}
}

func TestValidateReviewEvent_Valid(t *testing.T) {
	events := []string{"APPROVE", "REQUEST_CHANGES", "COMMENT"}
	for _, event := range events {
		if err := ValidateReviewEvent(event); err != nil {
			t.Errorf("ValidateReviewEvent(%q) unexpected error: %v", event, err)
		}
	}
}

func TestValidateReviewEvent_Invalid(t *testing.T) {
	events := []string{"", "approve", "REJECT", "LGTM", "invalid"}
	for _, event := range events {
		if err := ValidateReviewEvent(event); err == nil {
			t.Errorf("ValidateReviewEvent(%q) expected error", event)
		}
	}
}
