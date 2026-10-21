package heldout

import (
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"errors"
	"io"
	"strings"
	"testing"

	"csvwriter"
)

// A writer that fails on the first call.
type failWriter struct{ err error }

func (f failWriter) Write(p []byte) (int, error) { return 0, f.err }

// A writer that accepts everything but reports a short write.
type shortWriter struct{}

func (shortWriter) Write(p []byte) (int, error) { return len(p) / 2, nil }

func TestHeldoutGolden(t *testing.T) {
	var buf bytes.Buffer
	err := csvwriter.WriteRecords(&buf, []csvwriter.Record{{Name: "Ada", Email: "ada@example.com"}})
	if err != nil {
		t.Fatal(err)
	}
	want := "name,email\nAda,ada@example.com\n"
	if buf.String() != want {
		t.Fatalf("got %q want %q", buf.String(), want)
	}
}

func TestHeldoutCommaIsQuoted(t *testing.T) {
	var buf bytes.Buffer
	if err := csvwriter.WriteRecords(&buf, []csvwriter.Record{{Name: "Doe, John", Email: "j@example.com"}}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"Doe, John"`) {
		t.Fatalf("comma not quoted: %q", buf.String())
	}
}

func TestHeldoutQuoteIsDoubled(t *testing.T) {
	var buf bytes.Buffer
	if err := csvwriter.WriteRecords(&buf, []csvwriter.Record{{Name: `He said "hi"`, Email: "h@example.com"}}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"He said ""hi"""`) {
		t.Fatalf("quote not doubled: %q", buf.String())
	}
}

func TestHeldoutNewlineIsQuoted(t *testing.T) {
	var buf bytes.Buffer
	if err := csvwriter.WriteRecords(&buf, []csvwriter.Record{{Name: "line one\nline two", Email: "l@example.com"}}); err != nil {
		t.Fatal(err)
	}
	r := csv.NewReader(strings.NewReader(buf.String()))
	rows, err := r.ReadAll()
	if err != nil {
		t.Fatalf("output is not valid CSV: %v", err)
	}
	if len(rows) != 2 || rows[1][0] != "line one\nline two" {
		t.Fatalf("newline field not preserved: %#v", rows)
	}
}

func TestHeldoutErrorIsReturned(t *testing.T) {
	want := errors.New("disk full")
	err := csvwriter.WriteRecords(failWriter{want}, []csvwriter.Record{{Name: "Ada", Email: "a@example.com"}})
	if !errors.Is(err, want) {
		t.Fatalf("writer error swallowed: got %v", err)
	}
}

func TestHeldoutShortWriteIsAnError(t *testing.T) {
	err := csvwriter.WriteRecords(shortWriter{}, []csvwriter.Record{{Name: strings.Repeat("x", 4096), Email: "a@example.com"}})
	if err == nil {
		t.Fatal("short write reported as success")
	}
}

func TestHeldoutHeaderOnly(t *testing.T) {
	var buf bytes.Buffer
	if err := csvwriter.WriteRecords(&buf, []csvwriter.Record{}); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "name,email\n" {
		t.Fatalf("got %q", buf.String())
	}
}

func TestHeldoutEmptyFields(t *testing.T) {
	var buf bytes.Buffer
	if err := csvwriter.WriteRecords(&buf, []csvwriter.Record{{Name: "", Email: ""}}); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "name,email\n,\n" {
		t.Fatalf("got %q", buf.String())
	}
}

func TestHeldoutUnicode(t *testing.T) {
	var buf bytes.Buffer
	if err := csvwriter.WriteRecords(&buf, []csvwriter.Record{{Name: "José", Email: "j@example.com"}}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "José") {
		t.Fatalf("unicode lost: %q", buf.String())
	}
}

func TestHeldoutComposesWithGzip(t *testing.T) {
	var raw bytes.Buffer
	gz := gzip.NewWriter(&raw)
	if err := csvwriter.WriteRecords(gz, []csvwriter.Record{{Name: "Ada", Email: "ada@example.com"}}); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	zr, err := gzip.NewReader(&raw)
	if err != nil {
		t.Fatal(err)
	}
	out, err := io.ReadAll(zr)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "name,email\nAda,ada@example.com\n" {
		t.Fatalf("gzip payload wrong: %q", out)
	}
}
