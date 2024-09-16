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

		if _, err := tea.NewProgram(list.NewModel(), tea.WithAltScreen()).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

	},
}
