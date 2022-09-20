package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func encode(cfg *Config) error {
	cfg.Buffer = cfg.Buffer - (cfg.Buffer % 3) // Base64 encoding represents 3 bytes using 4 characters

	if cfg.Verbose {
		fmt.Printf("Encoding %v (buffer size: %v)...\n", display(cfg), cfg.Buffer)
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

	cnt, cnt1, off, maxl := 0, 0, 0, 0
	var err1 error
	buf, buf1 := make([]byte, 0, cfg.Buffer), make([]byte, 0, cfg.Buffer)

	for idx := 0; ; idx++ {
		if err1 == nil { // When loop for the last time, skip read
			cnt, err = rdr.Read(buf[:cap(buf)])
			if cfg.Verbose {
				verbose(idx, cnt, cfg, (wtr != nil))
			}
		}

		if cnt > 0 && cfg.Input == "" {
			// If getting input from stdin interactively, pressing <enter> would signify the end of an input line.
			if buf[:cnt][0] == 46 { // ASCII code 46 is period ('.')
				if cnt == 2 && buf[:cnt][1] == 10 { // ASCII code 10 is line feed LF ('\n')
					cnt = 0
					off = 1
					err = io.EOF
				} else if cnt == 3 && buf[:cnt][1] == 13 && buf[:cnt][2] == 10 { // ASCII code 13 is carriage return CR
					cnt = 0
					off = 2
					err = io.EOF
				}
			}
		}

		// As described in the doc, handle read data first if n > 0 before handling error,
		// it is because the returned error could have been EOF
		// NOTE: use cnt1, err1 and buf1 here because trying to delay processing for 1 cycle
		cnt1 -= off
		if cnt1 > 0 {
			encoded := base64.StdEncoding.EncodeToString(buf1[:cnt1])
			if len(encoded) > maxl {
				maxl = len(encoded)
			}
			format := fmt.Sprintf("%%-%dv", maxl)
			if wtr == nil {
				if cfg.Verbose {
					fmt.Printf(format+" %v\n", encoded, buf1[:cnt1])
				} else {
					fmt.Print(encoded)
				}
			} else {
				if cfg.Verbose {
					fmt.Printf(format+" %v\n", encoded, buf1[:cnt1])
				}
				fmt.Fprint(wtr, encoded)
			}
		}
		if cfg.Verbose && idx == 0 {
			fmt.Println()
		}

		if err1 != nil {
			if err1 == io.EOF {
				break // Done
			} else {
				return err1
			}
		}

		cnt1 = cnt
		err1 = err
		copy(buf1[:cnt], buf[:cnt])
	}

	if wtr == nil {
		fmt.Println()
	} else {
		wtr.Flush()
	}

	if cfg.Verbose {
		fmt.Println("Encoding finished")
	}

	return nil
}
