package cmd

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/certifi/gocertifi"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	health_pb "google.golang.org/grpc/health/grpc_health_v1"

	ffmt "github.com/gomicro/probe/fmt"
)

var (
	verbose    bool
	skipVerify bool
	grpcFlag   bool

	// ErrGrpcBadStatus is the error returned when a non serving status is returend by a grpc service
	ErrGrpcBadStatus = errors.New("grpc status: not ok")
	// ErrGrpcConnFailure is the error returned when the grpc dial fails to connect to the specified host
	ErrGrpcConnFailure = errors.New("grpc conn: failure")
	// ErrHTTPBadStatus is the error returned when a non ok status is returned by a http service
	ErrHTTPBadStatus = errors.New("http status: not ok")
	// ErrHTTPGet is the errro returned when an http client encounters an error while performing a GET
	ErrHTTPGet = errors.New("http client: get")
)

func init() {
	cobra.OnInitialize(initEnvs)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "show more verbose output")
	rootCmd.Flags().BoolVarP(&skipVerify, "insecure", "k", false, "permit operations for servers otherwise considered insecure")
	rootCmd.Flags().BoolVarP(&grpcFlag, "grpc", "g", false, "use http2 transport for grpc health check")
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
	var err error

	if grpcFlag {
		err = probeGrpc(args[0])
	} else {
		err = probeHttp(args[0])
	}

	if err != nil {
		if errors.Is(err, ErrHTTPGet) || errors.Is(err, ErrGrpcConnFailure) {
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
		return fmt.Errorf("%w: %v", ErrHTTPGet, err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: status %v", ErrHTTPBadStatus, resp.StatusCode)
	}

	return nil
}

func probeGrpc(host string) error {
	ctx := context.Background()

	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}

	conn, err := grpc.DialContext(ctx, host, opts...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrGrpcConnFailure, err.Error())
	}
	defer conn.Close()

	resp, err := health_pb.NewHealthClient(conn).Check(ctx, &health_pb.HealthCheckRequest{})
	if err != nil {
		return ErrGrpcBadStatus
	}

	if resp.GetStatus() != health_pb.HealthCheckResponse_SERVING {
		return fmt.Errorf("%w: status %v", err, resp.GetStatus())
	}

	return nil
}
