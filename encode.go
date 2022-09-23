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

	// v0.3.1
	// Since base64 encoding represents 3 bytes using 4 characters. To simplify the logic, the buffer size is set to a multiple of 3.
	// However according to the doc, the function NewReaderSize() "...returns a new Reader whose buffer has at least the specified
	// size. If the argument io.Reader is already a Reader with large enough size, it returns the underlying Reader..."
	// Testings show that apparently the reader returned by os.Open() has a buffer size of 16. Therefore if the specified 'basesf'
	// buffer size is less than 16, the buffered reader will have a size of 16 instead of the desired size. This will cause problem
	// when encoding, as the reader will pause after reading a total of 16 characters, cutting a 3-character unit into 2 parts.
	// Therefore need to re-create a new buffered reader if the reader's size does not match the specified 'basesf' buffer size.
	//
	// v0.4.0
	// Relying on the buffer size is a multiple of 3 is not safe. For example when receiving input from pipe, it seems that both in
	// Linux and MaxOS the shell can only pipe 65536 characters at a time. According to the doc the buffered reader may also return
	// less characters (and not multiple of 3) for unspecified reasons. It is more robust to make sure the number of characters
	// processed is multiple of 3 except for the last round.

	rdr := bufio.NewReaderSize(inp, cfg.Buffer)
	if rdr.Size() != cfg.Buffer {
		if cfg.Verbose {
			fmt.Printf("Read buffer size %v mismatching with the specified size %v, changing buffer size...\n", rdr.Size(), cfg.Buffer)
		}
		cfg.Buffer = rdr.Size()
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

	cnt, cnt1, cnt2, off, maxl := 0, 0, 0, 0, 0
	var err1 error
	buf, buf1 := make([]byte, 0, cfg.Buffer), make([]byte, 0, cfg.Buffer*2)

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
			cnt2 = cnt1 - (cnt1 % 3)
			if err != nil {
				cnt2 = cnt1 // Process everything when looping for the last time
			}

			encoded := base64.StdEncoding.EncodeToString(buf1[:cnt2])
			if len(encoded) > maxl {
				maxl = len(encoded)
			}

			if cfg.Verbose {
				verboseDtls(encoded, buf1[:cnt2], maxl, true, cfg)
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

		buf1 = append(buf1[cnt2:cnt1], buf[:cnt]...)
		cnt1 = cnt1 - cnt2 + cnt
		err1 = err
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
