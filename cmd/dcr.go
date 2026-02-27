package cmd

import (
	"fmt"

	"github.com/adobe/imscli/ims"
	"github.com/spf13/cobra"
)

func dcrCmd(imsConfig *ims.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dcr",
		Short: "Dynamic Client Registration operations.",
		Long:  `The dcr command enables Dynamic Client Registration operations.`,
	}
	cmd.AddCommand(
		registerCmd(imsConfig),
	)
	return cmd
}

func registerCmd(imsConfig *ims.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a client.",
		Long:  `Register a new OAuth client using Dynamic Client Registration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true


			resp, err := imsConfig.Register()
			if err != nil {
				return fmt.Errorf("error during client registration: %v", err)
			}

			fmt.Printf("Status Code: %d\n", resp.StatusCode)
			fmt.Println(resp.Body)
			return nil
		},
	}

	cmd.Flags().StringVarP(&imsConfig.ClientName, "clientName", "n", "", "Client application name.")
	cmd.Flags().StringSliceVarP(&imsConfig.RedirectURIs, "redirectURIs", "r", []string{}, "Redirect URIs (comma-separated or multiple flags).")

	return cmd
}
