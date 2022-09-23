# `basesf`
### ***BASE***-***S***ixty-***F***our encoding / decoding tool

A simple tool for base64 encoding and decoding.

## Build and Install
1. [Install](https://go.dev/doc/install) golang
> ```bash
> $ go version
> go version go1.17 xxx
> ```

2. [Clone](https://github.com/pangduckwai/basesf) the repository from GitHub
> ```bash
> $ git clone https://github.com/pangduckwai/basesf.git
> ```

3. Build the executable `basesf`
> ```bash
> $ cd .../basesf
> $ go build
> $ sudo ln -s basesf /usr/local/bin/basesf # for example
> ```

## Usage
```
Usage:
 basesf [encode | decode | version | help]
   {-i FILE | --in=FILE}
   {-o FILE | --out=FILE}
   {-b SIZE | --buffer=SIZE}
   {-v | --verbose}
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
  - `-v` | `--verbose`
    - display detail operation messages during processing if specified

- Notes
  - When inputting from `stdin` interactively, type a period (.) then press <enter> at a new line indicates there is no more input

## Changelog
### v0.4.0
- base64 encoding processes 3 bytes at a time and if data size is larger than buffer size, need to ensure the size of each chunk of data read is a multiple of 3. To simplify the logic, the older versions archieve this by controlling the buffer size of the reader. This version change to a more robust solution to control the actual number of bytes processed for each chunk

### v0.3.1
- fix the problem when specified buffer size less than 16

### v0.3.0
- change handling of trailing CR and/or LF: To ignore trailing CR/LF when reading interactively from stdin
- allows multi-line inputs when reading interactively from stdin

### v0.2.0
- add support of `stdin` as input.

### v0.1.1
- fix decoding output to `stdout`.
- fix encoding/decoding problem when buffer size smaller than the input file
- move `version` and `help` from options to commands

### v0.1.0
- first usable version
