package csvwriter

import (
	"io"
	"strings"
)

// Record is one row of the export.
type Record struct {
	Name  string
	Email string
}

// WriteRecords writes a CSV document to w. The full contract is in README.md.
func WriteRecords(w io.Writer, recs []Record) error {
	var b strings.Builder
	b.WriteString("name,email\n")
	for _, r := range recs {
		b.WriteString(r.Name)
		b.WriteString(",")
		b.WriteString(r.Email)
		b.WriteString("\n")
	}
	w.Write([]byte(b.String()))
	return nil
}
