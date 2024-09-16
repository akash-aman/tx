package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove specific templates",
	Long:  `Remove specific templates from the list of available templates.`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 && args[0] == "" {
			fmt.Println("Please provide a template name")
			os.Exit(1)
		}

		template := args[0]

		var templatePath = filepath.Join(os.Getenv("HOME"), ".tx/templates", template)

		/**
		 * Check if the template exists.
		 */
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			fmt.Println("Template does not exist")
			os.Exit(1)
		}

		/**
		 * Remove the template.
		 */
		os.RemoveAll(templatePath)

		fmt.Println("Template removed successfully :" + template)
	},
}

func init() {
	rmCmd.Flags().StringP("name", "n", "", "Name of the template")
	rmCmd.Flags().BoolP("all", "a", false, "Path to the template")
}
