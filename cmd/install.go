package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	tmpl "text/template"
	"time"

	cnf "github.com/akash-aman/tx/config"
	c "github.com/akash-aman/tx/constants"
	hd "github.com/akash-aman/tx/helper/directory"
	"github.com/akash-aman/tx/ui/loading"
	"github.com/akash-aman/tx/ui/progress"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "i",
	Short: "Install tx command line tool",
	Long:  `Install tx command line tool for managing projects templates.`,

	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 && args[0] == "" {
			fmt.Println("Please provide a template name")
			os.Exit(1)
		}

		template := args[0]

		if template != "" {

			var templatePath = filepath.Join(os.Getenv("HOME"), ".tx/templates", template)
			var desFolder = hd.GetCwd(".")

			/**
			 * Check if the template exists.
			 */
			if _, err := os.Stat(templatePath); os.IsNotExist(err) {
				fmt.Println("Template does not exist")
				os.Exit(1)
			}

			cnf.GenConfig.Path = desFolder
			cnf.GenConfig.Load()

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
			count := hd.CountFiles(templatePath)
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

			err := filepath.Walk(templatePath, func(srcPath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				/**
				 * Skip the directory.
				 */
				if info.IsDir() {
					return nil
				}

				/**
				 * Get the relative path of the file.
				 */
				relPath, err := filepath.Rel(templatePath, srcPath)
				if err != nil {
					return err
				}
				destPath := filepath.Join(desFolder, relPath)

				/**
				 * loop over ignores and skip the file if it matches.
				 */
				destDir := filepath.Dir(destPath)
				if err := os.MkdirAll(destDir, 0755); err != nil {
					return fmt.Errorf("error creating directories: %s", err)
				}

				fileName := info.Name()
				if strings.HasSuffix(fileName, ".txt") && strings.Count(fileName, ".") > 1 {

					/**
					 * Read the template file.
					 */
					tmplContent, err := os.ReadFile(srcPath)
					if err != nil {
						log.Fatalf("error reading template file: %s", err)
					}

					/**
					 * Parse the template.
					 */
					templateParsed, err := tmpl.New(relPath).Parse(string(tmplContent))
					if err != nil {
						log.Fatalf("error parsing template: %s", err)
					}

					/**
					 * Create the file.
					 */
					outFile, err := os.Create(strings.TrimSuffix(destPath, ".txt"))
					if err != nil {
						return fmt.Errorf("error creating file: %s", err)
					}
					defer outFile.Close()

					/**
					 * Execute the template with the data.
					 */
					err = templateParsed.Execute(outFile, cnf.GenConfig.Template)
					if err != nil {
						return fmt.Errorf("error executing template: %s", err)
					}
				}

				/**
				 * Last step of the progressBar.
				 */
				progressBar.UpdateProgress()
				return nil
			})

			progressBar.UpdateProgress()

			if err != nil {
				fmt.Println("Error counting files:", err)
				os.Exit(1)
			}
		} else {

			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			/**
			 * Create recursive folder.
			 */
			templateDir := filepath.Join(homeDir, ".tx/templates")
			if _, err := os.Stat(templateDir); os.IsNotExist(err) {
				err = os.MkdirAll(templateDir, 0755)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			/**
			 * Create txconfig.json file.
			 */
			configFile := filepath.Join(homeDir, ".tx/txconfig.json")
			if _, err := os.Stat(configFile); os.IsNotExist(err) {
				file, err := os.Create(configFile)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				defer file.Close()
				file.WriteString(`{}`)
			}
		}
	},
}
