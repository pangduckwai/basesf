package main

import (
	"fmt"
	"strconv"
	"strings"
)

const bUFFER = 1048576

type Config struct {
	Command uint8  // 1 - encode; 2 - decode
	Input   string // nil - stdin
	Output  string // nil - stdout
	Buffer  int    // buffer size
	Verbose bool
}

func usage() string {
	return "Usage:\n basesf [encode | decode | version | help]\n" +
		"   {-i FILE | --in=FILE}\n" +
		"   {-o FILE | --out=FILE}\n" +
		"   {-b SIZE | --buffer=SIZE}\n" +
		"   {-v | --verbose}"
}

func help() string {
	return fmt.Sprintf("Usage: basesf [commands] {options}\n"+
		"  * commands:\n"+
		"    encode  - convert input into base64 encoded string\n"+
		"    decode  - convert base64 encoded string back to the original form\n"+
		"    version - display current version of 'basesf'\n"+
		"    help    - display this message\n"+
		"  * options:\n"+
		"    -i FILE, --in=FILE\n"+
		"       name of the input file, omitting means input from stdin\n"+
		"    -o FILE, --out=FILE\n"+
		"       name of the output file, omitting means output to stdout\n"+
		"    -b SIZE, --buffer=SIZE\n"+
		"       size of the read buffer (SIZE default: %vKB)\n"+
		"    -v, --verbose\n"+
		"       display detail operation messages during processing", bUFFER/1024)
}

func parse(args []string) (cfg *Config, err error) {
	if len(args) < 2 {
		return nil, &Err{1, "Command missing"}
	}

	cfg = &Config{
		Buffer:  bUFFER,
		Verbose: false,
	}

	switch args[1] {
	case "test":
		cfg.Command = 0
	case "encode":
		cfg.Command = 1
	case "decode":
		cfg.Command = 2
	case "help":
		cfg.Command = 3
	case "version":
		cfg.Command = 4
	default:
		return nil, &Err{1, fmt.Sprintf("Invalid command '%v'", args[1])}
	}

	var val int
	for i := 2; i < len(args); i++ {
		switch {
		case args[i] == "-v" || args[i] == "--verbose":
			cfg.Verbose = true
		case args[i] == "-i":
			i++
			if i >= len(args) {
				return nil, &Err{2, "Missing input filename argument"}
			} else {
				cfg.Input = args[i]
			}
		case strings.HasPrefix(args[i], "--in="):
			if len(args[i]) <= 5 {
				return nil, &Err{2, "Missing input filename"}
			} else {
				cfg.Input = args[i][5:]
			}
		case args[i] == "-o":
			i++
			if i >= len(args) {
				return nil, &Err{3, "Missing output filename argument"}
			} else {
				cfg.Output = args[i]
			}
		case strings.HasPrefix(args[i], "--out="):
			if len(args[i]) <= 6 {
				return nil, &Err{3, "Missing out filename"}
			} else {
				cfg.Output = args[i][6:]
			}
		case args[i] == "-b":
			i++
			if i >= len(args) {
				return nil, &Err{4, "Missing buffer size argument"}
			} else {
				val, err = strconv.Atoi(args[i])
				if err == nil {
					cfg.Buffer = val
				}
			}
		case strings.HasPrefix(args[i], "--buffer="):
			if len(args[i]) <= 9 {
				return nil, &Err{4, "Missing buffer size"}
			} else {
				val, err = strconv.Atoi(args[i][9:])
				if err == nil {
					cfg.Buffer = val
				}
			}
		default:
			return nil, &Err{0, fmt.Sprintf("Invalid option '%v'", args[i])}
		}
	}

	return
}
