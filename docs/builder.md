# Builder

Builder creates a builder for the given type.
The builder can be used to create multiple instaces that are built similarly.

## Example

```go
//go:gombok Builder
type Name struct {
    firstname string
    lastname string
}

func main() {
    name := NewNameBuilder()
        .Firstname("John")
        .Lastname("Smith")
        .Build()
}
```

## Parameters

Builder has no parameters on struct level.

Field parameters:

| Parameter | Description                                               |
|-----------|-----------------------------------------------------------|
| exclude   | If this is set the field will not be added to the builder |
