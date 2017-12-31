package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/certifi/gocertifi"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Expected url to access.")
		os.Exit(1)
	}

	pool, err := gocertifi.CACerts()
	if err != nil {
		fmt.Printf("Error: failed to create cert pool: %v\n", err.Error())
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: pool},
	}
	client := &http.Client{Transport: transport}

	resp, err := client.Get(args[0])
	if err != nil {
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}
}
