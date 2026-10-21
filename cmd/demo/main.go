package main

import (
	"os"

	"csvwriter"
)

func main() {
	recs := []csvwriter.Record{
		{Name: "Ada", Email: "ada@example.com"},
		{Name: "Grace", Email: "grace@example.com"},
	}
	if err := csvwriter.WriteRecords(os.Stdout, recs); err != nil {
		panic(err)
	}
}
