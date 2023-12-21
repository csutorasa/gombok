# Setter

Setter creates setter function for a given field.

## Example

```go
//go:gombok Setter
type Name struct {
    firstname string
    lastname string
}

func main() {
    name := &Name{
        firstname: "John",
        lastname: "Smith",
    }
    name.SetFirstname("Adam")
}
```

## Parameters

Setter has the same parameters on struct or field level.

Parameters:

| Parameter | Description                                                    |
|-----------|----------------------------------------------------------------|
| private   | If this is set the field will not be exported                  |
| chained   | If this is set the field will return a reference to the struct |
