// Package cmd is responsible for all the cli commands
// to manage the conduit service
package cmd

import (
	"github.com/joho/godotenv"
	"github.com/rsb/realworld-golang/app"
	"github.com/rsb/realworld-golang/app/conf"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("[init] no .env file used.")
	}

	cobra.OnInitialize(initConfig)
	template := `{{printf "%s: %s - version %s\n" .Name .Short .Version}}`
	rootCmd.SetVersionTemplate(template)

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix(strings.ToUpper(app.ServiceName))
}

var (
	cfgFile string
	build   = "develop"
)

// rootCmd is the top level cli command for conduit
var rootCmd = &cobra.Command{
	Use:   "conduit",
	Short: "cli tool to help develop, deploy and test the lola platform",
	Long: `lola platform cli tool is used for the following:
- Aid in development of the services
- API server management
- Admin tasks
`,
	Version: build,
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		root := app.RootDir()
		configName := conf.ConfigFileName
		viper.AddConfigPath(root)
		viper.SetConfigType("toml")
		viper.SetConfigName(configName)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("cli-init,initConfig, config-file:", viper.ConfigFileUsed())
	}
}

func Execute(b string) {
	build = b
	cobra.CheckErr(rootCmd.Execute())
}
