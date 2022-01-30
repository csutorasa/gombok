# Getter

Getter creates getter function for a given field.

## Example

```go
//gombok:Getter
type Name struct {
	firstname string
	lastname string
}

func main() {
    name := &Name{
        firstname: "John",
        lastname: "Smith",
    }
    name.GetFirstname()
}
```

## Parameters

Getter has the same parameters on struct or field level.

Parameters:

| Parameter | Description                                   |
|-----------|-----------------------------------------------|
| private   | If this is set the field will not be exported |
