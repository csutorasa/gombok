# gombok

Gombok is a go code generation library inspired by [Project Lombok](https://projectlombok.org/).

## How to install

Install go version 1.18 or newer.

You can downgrade to 1.13 but cannot use generics [here](https://github.com/csutorasa/gombok/tree/1.13).

Install the generator with the following command:

```bash
go install github.com/csutorasa/gombok@v1.18.0
```

## How to use

For command line options you may run `gombok -h`.

Enable default generation:

```go
//go:generate gombok
```

This enables processing all files in the project. To disable processing for specific files you can add:

```go
//go:gombok ignore
```

All generated files have this ignore flag.


You can add the gombok comments or tags to enable code generation.

```go
//go:gombok Stringer
type Example struct {
    //go:gombok Getter
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
