package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
)

func main() {
	cfg, err := parse(os.Args)
	if err != nil {
		if errr, ok := err.(*Err); !ok || errr.Code > 1 {
			log.Fatal(err)
		}
		log.Fatalf("%v\n%v\n%v\n", err, app(), usage())
	}

	switch cfg.Command {
	case 0:
		validate(cfg)
		err = encode(cfg)
	case 1:
		validate(cfg)
		err = decode(cfg)
	case 2:
		fmt.Printf("%v\n%v\n", app(), help())
	case 3:
		fmt.Println(app())
	}

	if err != nil {
		log.Fatal(err)
	}
}

func Version() string {
	return "0.2.1beta"
}

func app() string {
	return fmt.Sprintf("BASE-Sixty-Four encoding/decoding tool (version %v)", Version())
}

func validate(cfg *Config) {
	if cfg.Input != "" {
		if _, err := os.Stat(cfg.Input); errors.Is(err, os.ErrNotExist) {
			log.Fatalf("Input file '%v' does not exist\n", cfg.Input)
		} else if err != nil {
			log.Fatal(err)
		}
	}

	if cfg.Output != "" {
		if _, err := os.Stat(cfg.Output); err == nil {
			log.Fatalf("Output file '%v' already exists\n", cfg.Output)
		} else if !errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
		}
	}
}

func display(cfg *Config) string {
	inp := "stdin"
	if cfg.Input != "" {
		inp = cfg.Input
	}

	out := "stdout"
	if cfg.Output != "" {
		out = cfg.Output
	}

	return fmt.Sprintf("'%v' to '%v'", inp, out)
}

func verbose(idx, cnt int, cfg *Config, linefeed bool) {
	digits := int(math.Log10(float64(cfg.Buffer))) + 1
	format := fmt.Sprintf("%%%dv", digits)

	plr := "s"
	if cnt < 2 {
		plr = " "
	}

	lf := " | "
	if linefeed {
		lf = "\n"
	}

	fmt.Printf("%4v - read "+format+"/%v byte%v%v", idx, cnt, cfg.Buffer, plr, lf)
}

type Err struct {
	Code uint8
	Msg  string
}

func (e *Err) Error() string {
	return fmt.Sprintf("%v", e.Msg)
}
