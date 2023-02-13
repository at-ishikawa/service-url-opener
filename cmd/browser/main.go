package main

import (
	"fmt"
	"os"

	"github.com/at-ishikawa/service-url-opener/internal/command"
	"github.com/spf13/cobra"
)

func main() {
	if err := runMain(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)

}

func runMain() error {
	rootCmd := &cobra.Command{
		Use:   "browser",
		Short: "Open a page on a browser",
	}
	gcpCommandFactory := command.InitGCPCommandFactory()
	rootCmd.AddCommand(gcpCommandFactory.Create())
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("rootCmd.Execute > %w", err)
	}

	return nil
}
