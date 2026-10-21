package csvwriter

import (
	"bytes"
	"strings"
	"testing"
)

func TestWritesHeaderAndRows(t *testing.T) {
	var buf bytes.Buffer
	err := WriteRecords(&buf, []Record{
		{Name: "Ada", Email: "ada@example.com"},
		{Name: "Grace", Email: "grace@example.com"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "name,email\nAda,ada@example.com\nGrace,grace@example.com\n"
	if buf.String() != want {
		t.Fatalf("got %q, want %q", buf.String(), want)
	}
}

func TestEmptySliceWritesHeader(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteRecords(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != "name,email\n" {
		t.Fatalf("got %q", buf.String())
	}
}

func TestSimpleNames(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteRecords(&buf, []Record{{Name: "Linus", Email: "linus@example.com"}}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "Linus,linus@example.com") {
		t.Fatalf("row missing: %q", buf.String())
	}
}
