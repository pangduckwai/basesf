package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func decode(cfg *Config) error {
	cfg.Buffer = cfg.Buffer - (cfg.Buffer % 4) // Base64 encoding represents 3 bytes using 4 characters

	if cfg.Verbose {
		fmt.Printf("Decoding %v (buffer size: %v)...\n", display(cfg), cfg.Buffer)
	}

	var err error
	inp := os.Stdin
	if cfg.Input != "" {
		inp, err = os.Open(cfg.Input)
		if err != nil {
			return err
		}
	}
	rdr := bufio.NewReaderSize(inp, cfg.Buffer)

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
	for idx := 0; ; idx++ {
		cnt, err := rdr.Read(buf[:cap(buf)])
		if cfg.Verbose {
			verbose(idx, cnt, cfg, (wtr != nil))
		}

		// As described in the doc, process read data first if n > 0 before
		// handling error, which could have been EOF
		if cnt > 0 {
			if cfg.Input == "" && buf[:cnt][cnt-1] == '\n' {
				err = io.EOF
			}

			decoded, errr := base64.StdEncoding.DecodeString(string(buf[:cnt]))
			if errr != nil {
				return &Err{21, errr.Error()}
			}
			if wtr == nil {
				fmt.Println(decoded) // Not terribly useful here...
			} else {
				_, errr = wtr.Write(decoded) // Write must return error if # of bytes written < len(decoded), so the # of bytes returned can be ignored
				if errr != nil {
					return errr
				}
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

	if wtr == nil {
		fmt.Println()
	} else {
		wtr.Flush()
	}

	if cfg.Verbose {
		fmt.Println("Decoding finished")
	}

	return nil
}
