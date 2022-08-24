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
		return nil, &Err{1, fmt.Sprintf("Invalid command '%v'", args[1])}
	}

	for i := 2; i < len(args); i++ {
		switch {
		case args[i] == "-i":
			i++
			if i >= len(args) {
				return nil, &Err{2, "Missing input filename argument"}
			}
			cfg.Input = args[i]
		case strings.HasPrefix(args[i], "--in="):
			if len(args[i]) <= 5 {
				return nil, &Err{2, "Missing input filename"}
			}
			cfg.Input = args[i][5:]
		case args[i] == "-o":
			i++
			if i >= len(args) {
				return nil, &Err{3, "Missing output filename argument"}
			}
			cfg.Output = args[i]
		case strings.HasPrefix(args[i], "--out="):
			if len(args[i]) <= 6 {
				return nil, &Err{3, "Missing out filename"}
			}
			cfg.Output = args[i][6:]
		case strings.HasPrefix(args[i], "--buffer="):
			if len(args[i]) <= 9 {
				return nil, &Err{4, "Missing buffer size"}
			}
			val, err := strconv.Atoi(args[i][9:])
			if err != nil {
				return nil, err
			}
			cfg.Buffer = val
		default:
			return nil, &Err{0, fmt.Sprintf("Invalid option '%v'", args[i])}
		}
	}

	return
}
