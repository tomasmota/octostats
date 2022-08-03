package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"n8m.io/octostats/cmd/deployments"
)

func Execute() {
	cmd := &cobra.Command{
		Use:     "octostats",
		Version: "0.0.1",
		Short:   "Cli for interacting with octopus deploy",
		Long:    `octostats is a cli for showing octopus deploy stats`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.PersistentFlags().String("apikey", os.Getenv("OCTOPUS_API_KEY"), "Octopus Deploy API Key")
	cmd.PersistentFlags().String("url", "http://octopus.dac.local", "Octopus Deploy Server Endpoint")
	cmd.PersistentFlags().String("format", "text", "Octopus Deploy Server Endpoint")

	viper.BindPFlag("apikey", cmd.PersistentFlags().Lookup("apikey"))
	viper.BindPFlag("url", cmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("format", cmd.PersistentFlags().Lookup("format"))

	cmd.AddCommand(deployments.NewDeploymentsCmd())

	cobra.CheckErr(cmd.Execute())
}
