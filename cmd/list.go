package cmd

import (
	"fmt"
	"os"

	"github.com/akash-aman/tx/ui/list"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all templates",
	Long:  `List all templates available in the list of available templates.`,

	/**
	 * Run the command.
	 */
	Run: func(cmd *cobra.Command, args []string) {
		//nameFlag, _ := cmd.Flags().GetString("name")

		//var name string
		/**
		 * Determine the folder name.
		 */
		// if nameFlag != "" {
		// 	name = nameFlag
		// } else if len(args) > 0 {
		// 	name = args[0]
		// }

		if _, err := tea.NewProgram(list.NewModel(), tea.WithAltScreen()).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

	},
}

func init() {
	/**
	 * Add the flags to the command.
	 */
	//listCmd.Flags().StringP("name", "n", "", "Name of the template")
}
