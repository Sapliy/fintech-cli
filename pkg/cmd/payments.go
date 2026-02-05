package cmd

import (
	"context"
	"fmt"
	"os"

	fintech "github.com/sapliy/fintech-sdk-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var paymentsCmd = &cobra.Command{
	Use:   "payments",
	Short: "Manage payments",
}

var createPaymentCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a payment",
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := viper.GetString("api_key")
		if apiKey == "" {
			fmt.Println("Error: Not authenticated. Use 'sapliy auth login'.")
			os.Exit(1)
		}

		amount, _ := cmd.Flags().GetInt64("amount")
		currency, _ := cmd.Flags().GetString("currency")

		client := fintech.NewClient(apiKey)
		zone := viper.GetString("current_zone")
		payment, err := client.Payments.CreateIntent(context.Background(), &fintech.PaymentIntentRequest{
			Amount:   amount,
			Currency: currency,
			ZoneID:   zone,
		})

		if err != nil {
			fmt.Printf("Error creating payment: %v\n", err)
			return
		}

		fmt.Printf("Payment created successfully! ID: %s\n", payment.ID)
	},
}

func init() {
	rootCmd.AddCommand(paymentsCmd)
	paymentsCmd.AddCommand(createPaymentCmd)
	createPaymentCmd.Flags().Int64P("amount", "a", 0, "Amount in cents")
	createPaymentCmd.Flags().StringP("currency", "c", "USD", "Currency code")
	createPaymentCmd.MarkFlagRequired("amount")
}
