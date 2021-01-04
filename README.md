# gokraken [![GoDoc](https://godoc.org/github.com/yaustn/gokraken?status.svg)](https://pkg.go.dev/github.com/yaustn/gokraken) [![Go Report Card](https://goreportcard.com/badge/github.com/yaustn/gokraken)](https://goreportcard.com/report/github.com/yaustn/gokraken)
A lightweight Golang implementation of the ![Kraken REST API specification](https://www.kraken.com/en-us/features/api).

## Usage

Add the latest version to your go.mod
```
require github.com/yaustn/gokraken v1.0.0
```

Example REST API Call:
```
import "github.com/yaustn/gokraken"

func main() {
    client := gokraken.NewClient(<API Key>, <API Secret>)

    orders, err := client.GetOrders()
	if err != nil {
        // Handle errors
	}

	// Process orders
}
```
