package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kashu-02/notify-done-to-slack/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	SlackWebhookURL string
}

var config Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "notify-done-to-slack",
	Short: "Notify the completion of command execution by Slack.",
	Long:  `Notify the completion of command execution by Slack.`,
	Args:  cobra.ArbitraryArgs,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	RunE: app.NotifyDoneToSlack,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolP("output", "o", false, "Output the result of standard output.")

	rootCmd.Flags().Int32P("head", "e", 0, "Print the first n lines of standard output.")
	rootCmd.Flags().Int32P("tail", "t", 10, "Print the last n lines of standard output.")

	rootCmd.Flags().SetInterspersed(false)
}

func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cfgFilePath := home + "/.config/notify-done-to-slack/"
	viper.AddConfigPath(cfgFilePath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("Config file is not found.")
			configDir := filepath.Dir(cfgFilePath)
			if _, err := os.Stat(configDir); os.IsNotExist(err) {
				if err := os.MkdirAll(configDir, 0700); err != nil {
					fmt.Println("Error creating config directory:", err)
					os.Exit(1)
				}
			}
			if err := viper.SafeWriteConfig(); err != nil {
				fmt.Println("Error creating config file:", err)
			}
			os.Exit(1)
		} else {
			// Config file was found but another error was produced
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if err := viper.ReadInConfig(); err == nil {
        fmt.Println("Using config file:", viper.ConfigFileUsed())
    }
}
