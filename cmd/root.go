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

var rootCmd = &cobra.Command{
	Use:   "mrcon",
	Short: "Minecraft RCON client",
	Long:  `mrcon is a simple Minecraft RCON client written in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		if showVer {
			fmt.Printf("mrcon version: %s\ncommit: %s\nbuild date: %s\n", version, commitSHA, buildDate)
			return
		}

		if host == "" || port == 0 || password == "" || (!termMode && len(args) == 0) {
			fmt.Fprintln(os.Stderr, "Error: --host, --port, --password, and a command are required.")
			cmd.Help()
			os.Exit(1)
		}

		address := fmt.Sprintf("%s:%d", host, port)
		conn, err := rcon.Dial(address, password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close()

		if termMode {
			fmt.Println("Entering terminal mode. Type commands, Ctrl+C to exit.")
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
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					continue
				}
				if !silent {
					if raw {
						fmt.Print(resp)
					} else if noColor {
						fmt.Println(resp)
					} else {
						fmt.Println(resp) // Color output can be added here if needed
					}
				}
			}
			return
		}

		for _, cmdText := range args {
			resp, err := conn.Execute(cmdText)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}
			if !silent {
				if raw {
					fmt.Print(resp)
				} else if noColor {
					fmt.Println(resp)
				} else {
					fmt.Println(resp) // Color output can be added here if needed
				}
			}
			if wait > 0 {
				time.Sleep(time.Duration(wait) * time.Second)
			}
		}
	},
}

func init() {
	rootCmd.Flags().StringVar(&host, "host", "H", "RCON server host (required)")
	rootCmd.Flags().IntVarP(&port, "port", "p", 0, "RCON server port (required)")
	rootCmd.Flags().StringVar(&password, "P", "", "RCON password (required)")
	rootCmd.Flags().BoolVar(&showVer, "version", false, "Show version information")
	rootCmd.Flags().BoolVar(&raw, "raw", false, "Output raw response without formatting")
	rootCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	rootCmd.Flags().BoolVar(&silent, "silent", false, "Suppress command output")
	rootCmd.Flags().IntVar(&wait, "wait", 0, "Wait time in seconds between commands")
	rootCmd.Flags().BoolVar(&termMode, "terminal", false, "Enable terminal mode for interactive commands")

	// Bind environment variables
	cobra.MarkFlagRequired(rootCmd.Flags(), "password")
}

func Execute() {
	// Execute the root command
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
