# Destructor

Destructor creates a function for the given type, that returns all the fields.

## Example

```go
//go:gombok Destructor
type Name struct {
    firstname string
    lastname string
}

func main() {
    name := &Name{
        firstname: "John",
        lastname: "Smith",
    }
    firstname, lastname := name.Destruct()
}
```

## Parameters

Struct parameters:

| Parameter | Description                                   |
|-----------|-----------------------------------------------|
| private   | If this is set the field will not be exported |

Field parameters:

| Parameter | Description                                                  |
|-----------|--------------------------------------------------------------|
| exclude   | If this is set the field will not be added to the destructor |
