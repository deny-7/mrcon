package cmd

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/gorcon/rcon"
	"github.com/spf13/cobra"
)

var (
	host     string
	port     int
	password string
	showVer  bool
	raw      bool
	noColor  bool
	silent   bool
	wait     int
	termMode bool

	version   = "dev"
	commitSHA = "none"
	buildDate = "unknown"
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mrcon",
		Short: "Minecraft RCON client",
		Long: `mrcon is a simple Minecraft RCON client written in Go.

Flag/env resolution order:
  1. CLI flags (e.g. --host, --port, --password)
  2. Environment variables: MRCON_HOST, MRCON_PORT, MRCON_PASSWORD
  3. If neither is set, the program will error.

Examples:
  mrcon --host 127.0.0.1 --port 25575 --password secret 'say hello'
  MRCON_HOST=127.0.0.1 MRCON_PORT=25575 MRCON_PASSWORD=secret mrcon 'say hello'
`,
	}
	cmd.Flags().StringVarP(&host, "host", "H", "", "RCON server host (required)")
	cmd.Flags().IntVarP(&port, "port", "p", 0, "RCON server port (required)")
	cmd.Flags().StringVarP(&password, "password", "P", "", "RCON password (required)")
	cmd.Flags().BoolVarP(&showVer, "version", "v", false, "Show version information")
	cmd.Flags().BoolVarP(&raw, "raw", "r", false, "Output raw response without formatting")
	cmd.Flags().BoolVarP(&noColor, "no-color", "n", false, "Disable colored output")
	cmd.Flags().BoolVarP(&silent, "silent", "s", false, "Suppress command output")
	cmd.Flags().IntVarP(&wait, "wait", "w", 0, "Wait time in seconds between commands")
	cmd.Flags().BoolVarP(&termMode, "terminal", "t", false, "Enable terminal mode for interactive commands")
	cmd.RunE = runRootCmd
	return cmd
}

func runRootCmd(cmd *cobra.Command, args []string) error {
	if showVer {
		fmt.Fprintf(cmd.OutOrStdout(), "mrcon version: %s\ncommit: %s\nbuild date: %s\n", version, commitSHA, buildDate)
		return nil
	}

	// Environment variable fallback
	if host == "" {
		host = os.Getenv("MRCON_HOST")
	}
	if port == 0 {
		if envPort := os.Getenv("MRCON_PORT"); envPort != "" {
			fmt.Sscanf(envPort, "%d", &port)
		}
	}
	if password == "" {
		password = os.Getenv("MRCON_PASSWORD")
	}

	if host == "" || port == 0 || password == "" || (!termMode && len(args) == 0) {
		fmt.Fprintln(cmd.ErrOrStderr(), "Error: --host, --port, --password, and a command are required. You can also set MRCON_HOST, MRCON_PORT, MRCON_PASSWORD.")
		cmd.Help()
		return fmt.Errorf("missing required flags")
	}

	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := rcon.Dial(address, password)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Failed to connect: %v\n", err)
		return err
	}
	defer conn.Close()

	if termMode {
		fmt.Fprintln(cmd.OutOrStdout(), "Entering terminal mode. Type commands, Ctrl+C to exit.")
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				break
			}
			cmdText := scanner.Text()
			if cmdText == "" {
				continue
			}
			resp, err := conn.Execute(cmdText)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
				continue
			}
			if !silent {
				if raw {
					fmt.Fprint(cmd.OutOrStdout(), resp)
				} else if noColor {
					fmt.Fprintln(cmd.OutOrStdout(), resp)
				} else {
					fmt.Fprintln(cmd.OutOrStdout(), resp) // Color output can be added here if needed
				}
			}
		}
		return nil
	}

	for _, cmdText := range args {
		resp, err := conn.Execute(cmdText)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
			continue
		}
		if !silent {
			if raw {
				fmt.Fprint(cmd.OutOrStdout(), resp)
			} else if noColor {
				fmt.Fprintln(cmd.OutOrStdout(), resp)
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), resp) // Color output can be added here if needed
			}
		}
		if wait > 0 {
			time.Sleep(time.Duration(wait) * time.Second)
		}
	}
	return nil
}

var rootCmd = newRootCmd()

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func SetVersion(v, c, d string) {
	version = v
	commitSHA = c
	buildDate = d
}
