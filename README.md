# `basesf`
<u>BASE</u> <u>S</u>ixty <u>F</u>our encoding / decoding tool

A simple tool for base64 encoding and decoding.

## Usage
```
Usage:
 basesf [encode|decode]
   {-h | --help} {-v | --version}
   {-i file | --in=file}
   {-o file | --out=file}
   {-b size | --buffer=size}
```

- Commands
  - `encode` - convert input into base64 encoded string
  - `decode` - convert base64 encoded string back to the original value

- Options
  - Inputs
    - `-i filename` | `--in=filename` - name of the input file, omitting means input from stdin
  - Outputs
    - `-o filename` | `--out=filename` - name of the file to write the output to, omitting means output to stdout
  - Buffer size
    - `-b size` | `--buffer=size` - the buffer size `basesf` uses to process large inputs, automatically round down to multiple of 3 for encoding and multiple of 4 for decoding

## Changelog
### Unreleased 1
- add support of stdin as input.

### Unreleased 0
- fix decoding output to stdout.

### v0.1.0
- first usable version
