# gombok

Gombok is a go code generation library inspired by [Project Lombok](https://projectlombok.org/).

## How to install

Install the generator with the following command:

```bash
go install github.com/csutorasa/gombok
```

## How to use

For command line options you may run `gombok -h`.

Enable default generation:

```go
//go:generate gombok
```

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
