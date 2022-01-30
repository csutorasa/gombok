# Stringer

Stringer implements the `fmt.Stringer` interface.

## Example

```go
//gombok:Stringer
type Name struct {
    firstname string
    lastname string
}

func main() {
    name := &Name{
        firstname: "John",
        lastname: "Smith",
    }
    fmt.Println(name.String())
}
```

## Parameters

Stringer has no parameters on struct level.

Field parameters:

| Parameter | Description                                              |
|-----------|----------------------------------------------------------|
| exclude   | If this is set the field will not be added to the string |
