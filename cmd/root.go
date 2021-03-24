package cmd

import (
	"crypto/tls"
	"errors"
	"net/http"
	"os"

	"golang.org/x/net/http2"

	"github.com/certifi/gocertifi"
	"github.com/spf13/cobra"

	"github.com/gomicro/probe/fmt"
)

var (
	verbose    bool
	skipVerify bool
	grpc       bool
)

func init() {
	cobra.OnInitialize(initEnvs)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "show more verbose output")
	rootCmd.Flags().BoolVarP(&skipVerify, "insecure", "k", false, "permit operations for servers otherwise considered insecure")
	rootCmd.Flags().BoolVarP(&grpc, "grpc", "g", false, "use http2 transport for grpc health check")
}

func initEnvs() {
}

var rootCmd = &cobra.Command{
	Use:   "probe [flags] URL",
	Short: "Lightweight healthchecker for scratch containers",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Expected url to access")
		}

		return nil
	},
	Run: probe,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Failed to execute: %v\n", err.Error())
		os.Exit(1)
	}
}

func probe(cmd *cobra.Command, args []string) {
	pool, err := gocertifi.CACerts()
	if err != nil {
		fmt.Printf("Error: failed to create cert pool: %v\n", err.Error())
	}

	tlsConfig := &tls.Config{
		RootCAs: pool,
	}

	if skipVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	var transport http.RoundTripper

	switch grpc {
	case true:
		transport = &http2.Transport{
			TLSClientConfig: tlsConfig,
		}
	default:
		transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	client := &http.Client{Transport: transport}

	resp, err := client.Get(args[0])
	if err != nil {
		fmt.Verbosef("error: http get: %v", err.Error())
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Verbosef("error: status %v", resp.StatusCode)
		os.Exit(1)
	}
}
