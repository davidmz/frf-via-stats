package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "frf-via-stats clio/*.zip > stats.jsons")
		os.Exit(1)
	}

	archives, err := filepath.Glob(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid file mask:", err)
		fmt.Fprintln(os.Stderr,
			"See https://golang.org/pkg/path/filepath/#Match for the mask syntax.")
		os.Exit(1)
	}

	if len(archives) == 0 {
		fmt.Fprintln(os.Stderr, "None of archives were found")
		return
	}

	fmt.Fprintln(os.Stderr, len(archives), "archive(s) were found")

	outEnc := json.NewEncoder(os.Stdout)

	for n, archFile := range archives {
		fmt.Fprintf(os.Stderr, "Processing file #%d (%s)... ", n+1, filepath.Base(archFile))
		stats, err := processArchive(archFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		} else {
			fmt.Fprintln(os.Stderr, "OK")
		}
		outEnc.Encode(stats)
	}

	fmt.Fprintln(os.Stderr, "All done")
}
