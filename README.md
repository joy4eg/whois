# WHOIS Lookup Library

This library provides functionality for performing WHOIS lookups for domain names. It supports multiple TLD (Top Level Domain) adapters and can automatically detect the appropriate WHOIS server based on the domain name.

## Installation

To install the library, use the following command:

```sh
go get -u github.com/joy4eg/whois
```

## Usage
```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/joy4eg/whois"
)

func main() {
    client, err := whois.New()
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    result, err := client.Whois(context.Background(), "example.com")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(result)
}
```

## Testing
```go
go test ./...
```
