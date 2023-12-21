# Wither

Wither creates a function that returns a new instance with a modified field.

## Example

```go
//go:gombok Wither
type Name struct {
    firstname string
    lastname string
}

func main() {
    name := &Name{
        firstname: "John",
        lastname: "Smith",
    }
    name2 := name.WithFirstname("Adam")
}
```

## Parameters

Wither has the same parameters on struct or field level.

Parameters:

| Parameter | Description                                   |
|-----------|-----------------------------------------------|
| private   | If this is set the field will not be exported |
