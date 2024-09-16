package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	cnf "github.com/akash-aman/tx/config"
	c "github.com/akash-aman/tx/constants"
	hd "github.com/akash-aman/tx/helper/directory"
	"github.com/akash-aman/tx/ui/loading"
	"github.com/akash-aman/tx/ui/progress"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [name] [path]",
	Short: "Add a new template",
	Long:  `Add a new template to the list of available templates.`,

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"name", "path"}, cobra.ShellCompDirectiveNoFileComp
	},

	Run: func(cmd *cobra.Command, args []string) {
		nameFlag, _ := cmd.Flags().GetString("name")
		pathFlag, _ := cmd.Flags().GetString("path")

		var name, path string

		// Determine the path to the template
		if pathFlag != "." {
			path = hd.GetCwd(pathFlag)
		} else if len(args) > 1 {
			path = hd.GetCwd(args[1])
		} else {
			path = hd.GetCwd(".")
		}

		// Determine the folder name
		if nameFlag != "" {
			name = nameFlag
		} else if len(args) > 0 && args[0] != "." {
			name = args[0]
		} else if len(args) > 0 && args[0] == "." {
			name = hd.GetDir(path)
		} else {
			name = hd.GetDir(path)
		}

		cnf.AddConfig.Path = path
		cnf.AddConfig.Load()

		/**
		 * Get the count of files in the template directory.
		 *
		 * main
		 *	|\
		 * 	| \
		 * 	|  |
		 * 	|  |- loading thread.
		 * 	|  |
		 * 	| /
		 * 	|/
		 */
		loading := loading.NewLoading(cmd.Context())
		loading.Messgage = c.CALCULATE_COUNT_MESSAGE
		loading.Run()
		time.Sleep(500 * time.Millisecond)
		count := hd.CountFiles(path)
		loading.Messgage = c.TEMPLATE_COUNT_MESSAGE(count)
		time.Sleep(500 * time.Millisecond)
		loading.End()
		loading.Wg.Wait()

		/**
		 * Start configuring the template.
		 *
		 * main
		 * 	|\
		 * 	| \
		 * 	|  |
		 * 	|  |- progress thread.
		 * 	|  |
		 * 	| /
		 * 	|/
		 */
		progressBar := progress.NewProgress(cmd.Context())
		progressBar.Count = count
		progressBar.Run()
		defer progressBar.Wg.Wait()

		/**
		 * Create the template directory.
		 */
		desFolder := filepath.Join(cnf.TemplateDir, name)

		err := filepath.Walk(path, func(srcPath string, info os.FileInfo, err error) error {

			/**
			 * loop over ignores and skip the file if it matches.
			 */
			for _, ignore := range cnf.AddConfig.Ignore {
				if strings.Contains(srcPath, ignore) {
					return nil
				}
			}

			if err != nil {
				return err
			}

			/**
			 * Skip directories.
			 */
			if info.IsDir() {
				return nil
			}

			/**
			 * Get the relative path of the file.
			 */
			relPath, err := filepath.Rel(path, srcPath)
			if err != nil {
				return err
			}

			/**
			 * Create the destination path.
			 */
			destPath := filepath.Join(desFolder, relPath)

			/**
			 * Create the directory if it doesn't exist.
			 */
			fileName := info.Name()
			if strings.HasSuffix(fileName, ".txt") && strings.Count(fileName, ".") > 1 ||
				(info.Name() == "tx.json" && path == hd.GetCwd(".")) {

				/**
				 * Create the directory if it doesn't exist.
				 */
				err := hd.CopyFile(srcPath, destPath)
				if err != nil {
					return err
				}
			} else {

				/**
				 * Create the directory if it doesn't exist.
				 */
				destPath = destPath + ".txt"
				err := hd.CopyFile(srcPath, destPath)
				if err != nil {
					return err
				}
			}

			progressBar.UpdateProgress()
			return nil
		})

		progressBar.UpdateProgress()

		if err != nil {
			fmt.Println("Error counting files:", err)
			os.Exit(1)
		}
	},
}

func init() {

	addCmd.Flags().StringP("name", "n", "", "Name of the template. (default: path directory name)")
	addCmd.Flags().StringP("path", "p", ".", "Path to the template directory")
}
