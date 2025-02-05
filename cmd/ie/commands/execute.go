package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/Azure/InnovationEngine/internal/engine"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/spf13/cobra"
)

// / Register the command with our command runner.
func init() {
	rootCommand.AddCommand(executeCommand)

	// Bool flags
	executeCommand.PersistentFlags().Bool("verbose", false, "Enable verbose logging & standard output.")
	executeCommand.PersistentFlags().Bool("do-not-delete", false, "Do not delete the Azure resources created by the Azure CLI commands executed.")

	// String flags
	executeCommand.PersistentFlags().String("correlation-id", "", "Adds a correlation ID to the user agent used by a scenarios azure-cli commands.")
	executeCommand.PersistentFlags().String("subscription", "", "Sets the subscription ID used by a scenarios azure-cli commands. Will rely on the default subscription if not set.")
	executeCommand.PersistentFlags().String("working-directory", ".", "Sets the working directory for innovation engine to operate out of. Restores the current working directory when finished.")

	// StringArray flags
	executeCommand.PersistentFlags().StringArray("var", []string{}, "Sets an environment variable for the scenario. Format: --var <key>=<value>")
}

var executeCommand = &cobra.Command{
	Use:   "execute [markdown file]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Execute the commands for an Azure deployment scenario.",
	Run: func(cmd *cobra.Command, args []string) {
		markdownFile := args[0]
		if markdownFile == "" {
			logging.GlobalLogger.Errorf("Error: No markdown file specified.")
			cmd.Help()
			os.Exit(1)
		}

		verbose, _ := cmd.Flags().GetBool("verbose")
		doNotDelete, _ := cmd.Flags().GetBool("do-not-delete")

		subscription, _ := cmd.Flags().GetString("subscription")
		correlationId, _ := cmd.Flags().GetString("correlation-id")
		environment, _ := cmd.Flags().GetString("environment")
		workingDirectory, _ := cmd.Flags().GetString("working-directory")

		environmentVariables, _ := cmd.Flags().GetStringArray("var")

		// Parse the environment variables from the command line into a map
		cliEnvironmentVariables := make(map[string]string)
		for _, environmentVariable := range environmentVariables {
			keyValuePair := strings.SplitN(environmentVariable, "=", 2)
			if len(keyValuePair) != 2 {
				logging.GlobalLogger.Errorf("Error: Invalid environment variable format: %s", environmentVariable)
				fmt.Printf("Error: Invalid environment variable format: %s", environmentVariable)
				cmd.Help()
				os.Exit(1)
			}

			cliEnvironmentVariables[keyValuePair[0]] = keyValuePair[1]
		}

		// Parse the markdown file and create a scenario
		scenario, err := engine.CreateScenarioFromMarkdown(markdownFile, []string{"bash", "azurecli", "azurecli-interactive", "terraform"}, cliEnvironmentVariables)
		if err != nil {
			logging.GlobalLogger.Errorf("Error creating scenario: %s", err)
			fmt.Printf("Error creating scenario: %s", err)
			os.Exit(1)
		}

		innovationEngine, err := engine.NewEngine(engine.EngineConfiguration{
			Verbose:          verbose,
			DoNotDelete:      doNotDelete,
			Subscription:     subscription,
			CorrelationId:    correlationId,
			Environment:      environment,
			WorkingDirectory: workingDirectory,
		})

		if err != nil {
			logging.GlobalLogger.Errorf("Error creating engine: %s", err)
			fmt.Printf("Error creating engine: %s", err)
			os.Exit(1)
		}

		// Execute the scenario
		err = innovationEngine.ExecuteScenario(scenario)
		if err != nil {
			logging.GlobalLogger.Errorf("Error executing scenario: %s", err)
			fmt.Printf("Error executing scenario: %s", err)
			os.Exit(1)
		}
	},
}
