package main

import (
	"fmt"
	"os"
	"os/exec"

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
	rootCmd.AddCommand(initGCPCommand())
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("rootCmd.Execute > %w", err)
	}

	return nil
}

func openBrowser(url string) error {
	// TODO: support multiple commands
	command := "open"

	if err := exec.Command(command, url).Run(); err != nil {
		return fmt.Errorf("exec.Command > %w", err)
	}
	return nil
}

func initGCPCommand() *cobra.Command {
	var projectId string
	gcpUrl := func(path string) string {
		baseURL := "https://console.cloud.google.com"
		return fmt.Sprintf("%s%s?project=%s", baseURL, path, projectId)
	}

	gcpCmd := &cobra.Command{
		Use:   "gcp",
		Short: "Open GCP console",
		Long: `print is for printing anything back to the screen.
For many years people have printed back to the screen.`,
		Args: cobra.MinimumNArgs(1),
	}

	gcpCmd.PersistentFlags().StringVarP(&projectId, "project", "", "", "Required. GCP project id")
	if err := gcpCmd.MarkPersistentFlagRequired("project"); err != nil {
		panic(err)
	}

	commands := map[string]string{
		"compute": "/compute/instances",
	}
	for serviceName, path := range commands {
		path := path
		gcpCmd.AddCommand(&cobra.Command{
			Use: serviceName,
			Run: func(cmd *cobra.Command, args []string) {
				if err := openBrowser(gcpUrl(path)); err != nil {
					fmt.Println(err)
				}
			},
		})
	}

	return gcpCmd
}
