package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func decode(cfg *Config) error {
	var err error
	inp := os.Stdin
	if cfg.Input != "" {
		inp, err = os.Open(cfg.Input)
		if err != nil {
			return err
		}
	}

	// According to the doc, the function NewReaderSize() "...returns a new Reader whose buffer has at least the specified size.
	// If the argument io.Reader is already a Reader with large enough size, it returns the underlying Reader..."
	// Testings show that apparently the reader returned by os.Open() has a buffer size of 16. Therefore if the specified 'basesf'
	// buffer size is less than 16, the buffered reader will have a size of 16 instead of the specified size. Re-creating the
	// buffered reader size to match the size returned by os.Open() since the memory is already allocated.
	cfg.Buffer = cfg.Buffer - (cfg.Buffer % 4) // Base64 encoding represents 3 bytes with 4 characters
	rdr := bufio.NewReaderSize(inp, cfg.Buffer)
	if rdr.Size() != cfg.Buffer {
		if cfg.Verbose {
			fmt.Printf("Read buffer size %v mismatching with the specified size %v, changing buffer size...\n", rdr.Size(), cfg.Buffer)
		}

		rmd := rdr.Size() % 4
		if rmd == 0 {
			cfg.Buffer = rdr.Size()
		} else {
			cfg.Buffer = rdr.Size() + 4 - rmd
		}
		rdr = bufio.NewReaderSize(inp, cfg.Buffer)
	}

	if cfg.Verbose {
		fmt.Printf("Decoding %v (buffer size: %v)...\n", display(cfg), cfg.Buffer)
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
			encoded := string(buf1[:cnt1])
			decoded, errr := base64.StdEncoding.DecodeString(encoded)
			if errr != nil {
				return &Err{21, errr.Error()}
			}
			if len(encoded) > maxl {
				maxl = len(encoded)
			}

			if cfg.Verbose {
				verboseDtls(encoded, decoded, maxl, false)
			}
			if wtr == nil {
				if !cfg.Verbose {
					fmt.Println(decoded) // Not terribly useful here...
				}
			} else {
				_, errr = wtr.Write(decoded) // Write must return error if # of bytes written < len(decoded), so the # of bytes returned can be ignored
				if errr != nil {
					return errr
				}
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
		fmt.Println("Decoding finished")
	}

	return nil
}
