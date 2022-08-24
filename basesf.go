package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
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

	if cfg.Input != "" {
		if _, err := os.Stat(cfg.Input); errors.Is(err, os.ErrNotExist) {
			log.Fatalf("Input file '%v' does not exist\n", cfg.Input)
		} else if err != nil {
			log.Fatal(err)
		}
	}

	if cfg.Output != "" {
		if _, err := os.Stat(cfg.Output); err == nil {
			log.Fatalf("Out file '%v' already exists\n", cfg.Output)
		} else if !errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
		}
	}

	switch cfg.Command {
	case 0:
		test()
	case 1:
		err = encode(cfg)
	case 2:
		err = decode(cfg)
	case 3:
		fmt.Printf("%v\n%v\n", app(), usage())
	case 4:
		fmt.Println(Version())
	}

	if err != nil {
		log.Fatal(err)
	}
}

func Version() string {
	return "0.1.0"
}

func app() string {
	return fmt.Sprintf("basesf: BASESixtyFour encoding/decoding tool (version %v)", Version())
}

func usage() string {
	return "Usage:\n basesf [encode|decode]\n" +
		"   {-h | --help} {-v | --version}\n" +
		"   {-i file | --in=file}\n" +
		"   {-o file | --out=file}\n" +
		"   {-b size | --buffer=size}"
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

	return fmt.Sprintf("'%v' to '%v'...", inp, out)
}

type Err struct {
	Code uint8
	Msg  string
}

func (e *Err) Error() string {
	return fmt.Sprintf("%v", e.Msg)
}

func test() {
	content := []byte{128, 182, 109, 169, 39, 17, 65, 10, 93, 201, 88, 143, 79, 5}

	// // Write content to file for testing
	// err := os.WriteFile("test.bin", content, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println(content)
	fmt.Println(content[:9])
	fmt.Println(content[9:])

	// Encode as base64.
	encoder := base64.StdEncoding //.WithPadding(base64.NoPadding)
	encoded0 := encoder.EncodeToString(content)
	encoded1 := encoder.EncodeToString(content[:9])
	encoded2 := encoder.EncodeToString(content[9:])

	// Print encoded data to console.
	fmt.Println("ENCODED 0: " + encoded0)
	fmt.Println("ENCODED 1: " + encoded1)
	fmt.Println("ENCODED 2: " + encoded2)
}
