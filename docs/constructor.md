# Constructor

Constructor creates a constructor function for the given type.

## Example

```go
//go:gombok Constructor
type Name struct {
    firstname string
    lastname string
}

func main() {
    name := NewName("John", "Smith")
}
```

## Parameters

Constructor has no parameters on struct level.

Field parameters:

| Parameter | Description                                                   |
|-----------|---------------------------------------------------------------|
| exclude   | If this is set the field will not be added to the constructor |
