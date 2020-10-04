Pipes
=====

---

![Screenshot](https://i.imgur.com/luk3uzj.png)

---

## Description

### TL;DR

This is a simple _pipe_ generator thingy just for fun.

## Making it works

### Building requirements
- Go

And nothing else is needed.

### Testing

```sh
go run pipes.go
```

### Building

```sh
git clone https://github.com/llathasa-veleth/pipes
cd pipes
go build
```

### Running

You can grab latest binaries from [releases](https://github.com/llathasa-veleth/pipes/releases).

```sh
chmod +x pipes
pipes --help  # seek for help
pipes         # run with default style
pipes -s thin # run with <thin> style
```

## TODO

- Makefile
- More options/flags
- Releases
- Windows support
