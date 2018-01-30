package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/certifi/gocertifi"
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initEnvs)
}

func initEnvs() {
}

var RootCmd = &cobra.Command{
	Use:   "probe",
	Short: "Lightweight healthchecker for scratch containers",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Expected url to access")
		}

		return nil
	},
	Run: probe,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Printf("Failed to execute: %v\n", err.Error())
		os.Exit(1)
	}
}

func probe(cmd *cobra.Command, args []string) {
	pool, err := gocertifi.CACerts()
	if err != nil {
		fmt.Printf("Error: failed to create cert pool: %v\n", err.Error())
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: true,
		},
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
