# `basesf`
### ***BASE***-***S***ixty-***F***our encoding / decoding tool

A simple tool for base64 encoding and decoding.

## Build
1. [Install](https://go.dev/doc/install) golang
> ```bash
> $ go version
> go version go1.17 xxx
> ```

1. [Clone](https://github.com/pangduckwai/basesf) the repository from GitHub
> ```bash
> $ git clone https://github.com/pangduckwai/basesf.git
> ```

1. Build the executable `basesf`
> ```bash
> $ cd .../basesf
> $ go build
```

## Usage
```
Usage:
 basesf [encode | decode | version | help]
   {-i FILE | --in=FILE}
   {-o FILE | --out=FILE}
   {-b SIZE | --buffer=SIZE}
```

- Commands
  - `encode`
    - convert input into base64 encoded string
  - `decode`
    - convert base64 encoded string back to the original form
  - `version`
    - display the current version
  - `help`
    - display the help message

- Options
  - `-i filename` | `--in=filename`
    - name of the input file, omitting means input from `stdin`
  - `-o filename` | `--out=filename`
    - name of the file to write the output to, omitting means output to `stdout`
  - `-b size` | `--buffer=size`
    - the buffer size used to read large inputs, automatically round down to multiple of 3 for encoding and multiple of 4 for decoding
  - `-v | --verbose`
    - display detail messages of the processing if specified

## Changelog
### v0.2.0
- add support of `stdin` as input.

### v0.1.1
- fix decoding output to `stdout`.
- fix encoding/decoding problem when buffer size smaller than the input file
- move `version` and `help` from options to commands

### v0.1.0
- first usable version
