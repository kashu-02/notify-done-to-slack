package app

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ConfigureWebhookURL(cmd *cobra.Command, args []string) error {
	webhookURL, _ := cmd.Flags().GetString("webhook-url")
	if webhookURL == "" {
		return fmt.Errorf("please specify the Slack Webhook URL")
	}

	viper.Set("webhook-url", webhookURL)

	// 設定ファイルに書き込み
	if err := viper.WriteConfig(); err != nil {
		fmt.Println("Error writing config:", err)
	} else {
		fmt.Println("Config written successfully")
	}
	return nil
}
