package cmd

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/certifi/gocertifi"
	"github.com/spf13/cobra"

	ffmt "github.com/gomicro/probe/fmt"
)

var (
	verbose    bool
	skipVerify bool
	grpc       bool

	ErrHttpGet       = errors.New("http client: get")
	ErrHttpBadStatus = errors.New("http status: not ok")
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
	err := probeHttp(args[0])
	if err != nil {
		if errors.Is(err, ErrHttpGet) {
			ffmt.Printf("%v", err.Error())
			os.Exit(2)
		}

		ffmt.Verbosef("%v", err.Error())
		os.Exit(1)
	}
}

func probeHttp(host string) error {
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

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{Transport: transport}

	resp, err := client.Get(host)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrHttpGet, err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: status %v", ErrHttpBadStatus, resp.StatusCode)
	}

	return nil
}
