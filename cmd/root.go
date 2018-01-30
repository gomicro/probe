package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/certifi/gocertifi"
	"github.com/spf13/cobra"
)

var (
	verbose    bool
	skipVerify bool
)

func init() {
	cobra.OnInitialize(initEnvs)

	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "show more verbose output")
	RootCmd.Flags().BoolVarP(&skipVerify, "insecure", "k", false, "permit operations for servers otherwise considered insecure")
}

func initEnvs() {
}

var RootCmd = &cobra.Command{
	Use:   "probe [flags] URL",
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
		printf("Failed to execute: %v\n", err.Error())
		os.Exit(1)
	}
}

func probe(cmd *cobra.Command, args []string) {
	pool, err := gocertifi.CACerts()
	if err != nil {
		printf("Error: failed to create cert pool: %v\n", err.Error())
	}

	tlsConfig := &tls.Config{
		RootCAs: pool,
	}

	if skipVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{Transport: transport}

	resp, err := client.Get(args[0])
	if err != nil {
		verbosef("error: http get: %v", err.Error())
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		verbosef("error: status %v", resp.StatusCode)
		os.Exit(1)
	}
}

func printf(f string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(f, args...))
}

func verbosef(f string, args ...interface{}) {
	if verbose {
		fmt.Println(fmt.Sprintf(f, args...))
	}
}
