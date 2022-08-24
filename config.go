package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Config struct {
	Command uint8  // 1 - encode; 2 - decode
	Input   string // nil - stdin
	Output  string // nil - stdout
	Buffer  int    // buffer size
}

func parse(args []string) (cfg *Config, err error) {
	cfg = &Config{
		Buffer: 1024,
	}

	if args[1] == "encode" {
		cfg.Command = 1
	} else if args[1] == "decode" {
		cfg.Command = 2
	} else if args[1] == "test" {
		cfg.Command = 0
	} else {
		err = &Err{1, fmt.Sprintf("Invalid command '%v'", args[1])}
	}

	var val int
	for i := 2; i < len(args); i++ {
		switch {
		case args[i] == "-h" || args[i] == "--help":
			cfg.Command = 3
			cfg.Input = ""
			cfg.Output = ""
			err = nil
			return
		case args[i] == "-v" || args[i] == "--version":
			cfg.Command = 4
			cfg.Input = ""
			cfg.Output = ""
			err = nil
			return
		case args[i] == "-i":
			i++
			if i >= len(args) {
				err = &Err{2, "Missing input filename argument"}
			} else {
				cfg.Input = args[i]
			}
		case strings.HasPrefix(args[i], "--in="):
			if len(args[i]) <= 5 {
				err = &Err{2, "Missing input filename"}
			} else {
				cfg.Input = args[i][5:]
			}
		case args[i] == "-o":
			i++
			if i >= len(args) {
				err = &Err{3, "Missing output filename argument"}
			} else {
				cfg.Output = args[i]
			}
		case strings.HasPrefix(args[i], "--out="):
			if len(args[i]) <= 6 {
				err = &Err{3, "Missing out filename"}
			} else {
				cfg.Output = args[i][6:]
			}
		case args[i] == "-b":
			i++
			if i >= len(args) {
				err = &Err{4, "Missing buffer size argument"}
			} else {
				val, err = strconv.Atoi(args[i])
				if err == nil {
					cfg.Buffer = val
				}
			}
		case strings.HasPrefix(args[i], "--buffer="):
			if len(args[i]) <= 9 {
				err = &Err{4, "Missing buffer size"}
			} else {
				val, err = strconv.Atoi(args[i][9:])
				if err == nil {
					cfg.Buffer = val
				}
			}
		default:
			err = &Err{0, fmt.Sprintf("Invalid option '%v'", args[i])}
		}
	}

	return
}
