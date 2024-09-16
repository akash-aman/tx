package template

import (
	"log"
	"os"

	cnf "github.com/akash-aman/tx/config"
	"github.com/spf13/viper"
)

type Template struct {
	Name        string
	Description string
}

func GetAllTemplates() []Template {
	/**
	 * Templates array.
	 */
	var templates []Template

	/**
	 * Get all folders in the templates directory.
	 */
	entries, err := os.ReadDir(cnf.TemplateDir)

	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			viper.SetConfigName(cnf.ConfigFileName)
			viper.SetConfigType(cnf.ConfigFileExt)
			viper.AddConfigPath(cnf.TemplateDir + "/" + entry.Name())
			viper.ReadInConfig()

			var template Template

			err := viper.Unmarshal(&template)
			template.Name = entry.Name()
			if err != nil {
				log.Println("Error unmarshalling config file, ", err)
			}

			templates = append(templates, template)
		}
	}

	return templates
}
