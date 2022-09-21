package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func encode(cfg *Config) error {
	var err error
	inp := os.Stdin
	if cfg.Input != "" {
		inp, err = os.Open(cfg.Input)
		if err != nil {
			return err
		}
	}

	// Since base64 encoding processes 3 characters at a time. To simplify the logic, the buffer size is set to a multiple of 3.
	// However according to the doc, the function NewReaderSize() "...returns a new Reader whose buffer has at least the specified
	// size. If the argument io.Reader is already a Reader with large enough size, it returns the underlying Reader..."
	// Testings show that apparently the reader returned by os.Open() has a buffer size of 16. Therefore if the specified 'basesf'
	// buffer size is less than 16, the buffered reader will have a size of 16 instead of the desired size. This will cause problem
	// when encoding, as the reader will pause after reading a total of 16 characters, cutting a 3-character unit into 2 parts.
	// Therefore need to re-create a new buffered reader if the reader's size does not match the specified 'basesf' buffer size.
	cfg.Buffer = cfg.Buffer - (cfg.Buffer % 3) // Base64 encoding represents 3 bytes with 4 characters
	rdr := bufio.NewReaderSize(inp, cfg.Buffer)
	if rdr.Size() != cfg.Buffer {
		if cfg.Verbose {
			fmt.Printf("Read buffer size %v mismatching with the specified size %v, changing buffer size...\n", rdr.Size(), cfg.Buffer)
		}

		rmd := rdr.Size() % 3
		if rmd == 0 {
			cfg.Buffer = rdr.Size()
		} else {
			cfg.Buffer = rdr.Size() + 3 - rmd
		}
		rdr = bufio.NewReaderSize(inp, cfg.Buffer)
	}

	if cfg.Verbose {
		fmt.Printf("Encoding %v (buffer size: %v)...\n", display(cfg), cfg.Buffer)
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

	cnt, cnt1, off, maxl := 0, 0, 0, 0
	var err1 error
	buf, buf1 := make([]byte, 0, cfg.Buffer), make([]byte, 0, cfg.Buffer)

	for idx := 0; ; idx++ {
		if err1 == nil { // When loop for the last time, skip read
			cnt, err = rdr.Read(buf[:cap(buf)])
			if cfg.Verbose {
				verboseHead(idx, cnt, cfg)
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

			if cfg.Verbose {
				verboseDtls(encoded, buf1[:cnt1], maxl, true)
			}
			if wtr == nil {
				if !cfg.Verbose {
					fmt.Print(encoded)
				}
			} else {
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

	if wtr != nil {
		wtr.Flush()
	} else if !cfg.Verbose {
		fmt.Println()
	}

	if cfg.Verbose {
		fmt.Println("Encoding finished")
	}

	return nil
}
