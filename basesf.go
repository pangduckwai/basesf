package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
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
	return "0.4.0beta1"
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

func verboseHead(idx, cnt int, cfg *Config) {
	digits := int(math.Log10(float64(cfg.Buffer))) + 1
	plr := "s"
	if cnt < 2 {
		plr = " "
	}

	fmt.Printf("%4v - read "+fmt.Sprintf("%%%dv", digits)+"/%v byte%v, ", idx, cnt, cfg.Buffer, plr) //, lf)
}

func verboseDtls(encoded string, decoded []byte, maxlen int, isEncoded bool, cfg *Config) {
	digits := int(math.Log10(float64(cfg.Buffer*2))) + 1

	msg := "decode"
	cnt := len(encoded)
	dirn := "->"
	if isEncoded {
		msg = "encode"
		cnt = len(decoded)
		dirn = "<-"
	}

	plr := "s"
	if cnt < 2 {
		plr = " "
	}

	display0 := maxlen
	if maxlen > dISPLAY1 {
		display0 = dISPLAY1
	}

	display1 := encoded
	if len(encoded) > dISPLAY1 {
		display1 = encoded[0:dISPLAY1-3] + "..."
	}

	display2 := fmt.Sprintf("%v", decoded)
	if len(decoded) > dISPLAY2 {
		tmp := fmt.Sprintf("%v", decoded[0:dISPLAY2])
		lidx := strings.LastIndex(tmp, " ")
		display2 = tmp[0:lidx] + " ...]"
	}

	fmt.Printf(fmt.Sprintf("%v %%%dv byte%v | %%-%dv", msg, digits, plr, display0)+" %v %v\n", cnt, display1, dirn, display2)
}

type Err struct {
	Code uint8
	Msg  string
}

func (e *Err) Error() string {
	return fmt.Sprintf("%v", e.Msg)
}
