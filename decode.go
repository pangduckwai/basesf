package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func decode(cfg *Config) error {
	if cfg.Input == "" {
		return &Err{21, "Reading input from stdin not yet supported"}
	}

	cfg.Buffer = cfg.Buffer - (cfg.Buffer % 4) // Base64 encoding represents 3 bytes using 4 characters

	fmt.Printf("Decoding %v to %v... (buffer size: %v)\n", cfg.Input, cfg.Output, cfg.Buffer)

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
	rdr := bufio.NewReader(inp)
	for {
		n, err := rdr.Read(buf[:cap(buf)])
		if n > 0 {
			encoded := string(buf[:n])
			dcd, errr := base64.StdEncoding.DecodeString(encoded)
			if errr != nil {
				return errr
			}
			_, errr = wtr.Write(dcd) // Write must return error if # of bytes written < len(dcd), so the # of bytes
			if errr != nil {
				return errr
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

	return nil
}
