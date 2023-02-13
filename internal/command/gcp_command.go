package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

type GCPCommandFactory struct {
	projectId   string
	isNoBrowser bool
}

type subCommandMetadata struct {
	path         string
	generatePath func(projectId string) string
	subCommands  map[string]subCommandMetadata
}

var (
	// command: url paths
	subCommands = map[string]subCommandMetadata{
		"compute": {path: "/compute/instances"},
		"sql":     {path: "/sql/instances"},
		"bq":      {path: "/bigquery"},

		"docker": {generatePath: func(projectId string) string {
			return "/gcr/images/" + projectId
		}},

		"container": {subCommands: map[string]subCommandMetadata{
			"clusters": {path: "/kubernetes/list/overview"},
		}},
	}
)

func InitGCPCommandFactory() GCPCommandFactory {
	return GCPCommandFactory{}
}

func (factory *GCPCommandFactory) Create() *cobra.Command {
	gcpCmd := &cobra.Command{
		Use:   "gcp",
		Short: "Open GCP console",
		Long: `print is for printing anything back to the screen.
For many years people have printed back to the screen.`,
		Args: cobra.MinimumNArgs(1),
	}

	gcpCmd.PersistentFlags().StringVarP(&factory.projectId, "project", "", "", "Required. GCP project id")
	gcpCmd.PersistentFlags().BoolVarP(&factory.isNoBrowser, "no-browser", "", false, "Optional. Doesn't open a browser if it's true")
	if err := gcpCmd.MarkPersistentFlagRequired("project"); err != nil {
		panic(err)
	}

	factory.addCommands(gcpCmd, subCommands)
	return gcpCmd
}

func (factory *GCPCommandFactory) addCommands(parentCommand *cobra.Command, subCommands map[string]subCommandMetadata) {
	for serviceName, subCommand := range subCommands {
		if subCommand.path != "" {
			path := subCommand.path
			parentCommand.AddCommand(&cobra.Command{
				Use: serviceName,
				Run: func(cmd *cobra.Command, args []string) {
					url := factory.gcpUrl(path)
					if factory.isNoBrowser {
						fmt.Println(url)
						return
					}

					if err := openBrowser(url); err != nil {
						fmt.Println(err)
					}
				},
			})
			continue
		}
		if subCommand.generatePath != nil {
			generatePath := subCommand.generatePath
			parentCommand.AddCommand(&cobra.Command{
				Use: serviceName,
				Run: func(cmd *cobra.Command, args []string) {
					url := factory.gcpUrl(generatePath(factory.projectId))
					if factory.isNoBrowser {
						fmt.Println(url)
						return
					}

					if err := openBrowser(url); err != nil {
						fmt.Println(err)
					}
				},
			})
			continue
		}
		if subCommand.subCommands != nil {
			subCobraCommand := &cobra.Command{
				Use: serviceName,
			}
			factory.addCommands(subCobraCommand, subCommand.subCommands)
			parentCommand.AddCommand(subCobraCommand)
			continue
		}
	}
}

func (factory GCPCommandFactory) gcpUrl(path string) string {
	baseURL := "https://console.cloud.google.com"
	return fmt.Sprintf("%s%s?project=%s", baseURL, path, factory.projectId)
}
