package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/CHTJonas/whoami-dns"
	"github.com/spf13/cobra"
)

// Software version defaults to the value below but is overridden by the compiler in Makefile.
var version = "dev-edge"

// The dnstap UNIX socket path and HTTP webserver port.
var sockPath, webPort string

// If the --version or -v command line option has been invoked.
var printVerAndExit bool

var rootCmd = &cobra.Command{
	Use:   "whoami-dns",
	Short: "HTTP server that works in tandem with an authoritative DNS server using dnstap",
	Long: "whoami-dns is a clever webserver that works in tandem with an authoritative DNS " +
		"server using dnstap and wildcard domains so that clients' recursive DNS servers can " +
		"be identified by the source IP addresses of their queries.",
	Run: func(cmd *cobra.Command, args []string) {
		if printVerAndExit {
			fmt.Println("whoami-dns version", version)
			os.Exit(0)
		}

		serv := whoami.NewServer()
		pwrBy := fmt.Sprintf("whoami-dns/%s Go/%s (+https://github.com/CHTJonas/whoami-dns)",
			version, strings.TrimPrefix(runtime.Version(), "go"))
		serv.SetHeader("X-Powered-By", pwrBy)

		serv.OpenSocket(sockPath)
		serv.Start(webPort)

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT)
		signal.Notify(quit, syscall.SIGTERM)
		<-quit

		serv.CloseSocket()
		serv.Stop()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVarP(&sockPath, "bind", "b", "/var/lib/knot/dnstap.sock", "path to dnstap UNIX socket")
	rootCmd.Flags().StringVarP(&webPort, "port", "p", "6780", "port on which to listen for HTTP requests")
	rootCmd.Flags().BoolVarP(&printVerAndExit, "version", "v", false, "print version and exit")
}

func initConfig() {
	// TODO
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
