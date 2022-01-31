# gombok

Gombok is a go code generation library inspired by [Project Lombok](https://projectlombok.org/).

## How to install

Install go version 1.13 or newer.
You can use a newer version if you switch to 1.16 or newer.
You can use generics if you switch to the 1.18 or newer version [here](https://github.com/csutorasa/gombok).

Install the generator with the following command:

```bash
go install github.com/csutorasa/gombok@v0.13.0
```

## How to use

For command line options you may run `gombok -h`.

Enable default generation:

```go
//go:generate gombok
```

This enables processing all files in the project. To disable processing for specific files you can add:

```go
//gombok:ignore
```

All generated files have this ignore flag.


You can add the gombok comments or tags to enable code generation.

```go
//gombok:Stringer
type Example struct {
    //gombok:Getter
    firstname string
    lastname string `gombokSetter:""`
}
```

## Available generators

- [Builder](docs/builder.md)
- [Constructor](docs/constructor.md)
- [Destructor](docs/destructor.md)
- [Getter](docs/getter.md)
- [Setter](docs/setter.md)
- [Stringer](docs/stringer.md)
- [Wither](docs/wither.md)
