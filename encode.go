package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func encode(cfg *Config) error {
	if cfg.Input == "" {
		return &Err{11, "Reading input from stdin not yet supported"}
	}

	cfg.Buffer = cfg.Buffer - (cfg.Buffer % 3) // Base64 encoding represents 3 bytes using 4 characters

	if cfg.Verbose {
		fmt.Printf("Encoding %v (buffer size: %v)...\n", display(cfg), cfg.Buffer)
	}

	inp, err := os.Open(cfg.Input)
	if err != nil {
		return err
	}

	var wtr *bufio.Writer
	if cfg.Output != "" {
		out, err := os.Create(cfg.Output)
		if err != nil {
			return err
		}

		wtr = bufio.NewWriter(out)

		defer out.Close()
	}

	buf := make([]byte, 0, cfg.Buffer)
	rdr := bufio.NewReaderSize(inp, cfg.Buffer)
	for idx := 0; ; idx++ {
		n, err := rdr.Read(buf[:cap(buf)])
		if cfg.Verbose {
			verbose(idx, n, cfg)
		}

		// As described in the doc, process read data first if n > 0 before
		// handling error, which could have been EOF
		if n > 0 {
			encoded := base64.StdEncoding.EncodeToString(buf[:n])
			if wtr == nil {
				fmt.Print(encoded)
			} else {
				fmt.Fprint(wtr, encoded)
			}
		}

		if err != nil {
			if err == io.EOF {
				break // Done
			} else {
				return err
			}
		}
	}

	if wtr != nil {
		wtr.Flush()
	}

	if cfg.Verbose {
		fmt.Println("Encoding finished")
	}

	return nil
}
