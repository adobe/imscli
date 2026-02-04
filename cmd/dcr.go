package cmd

import (
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
		ims.RegisterCmd(imsConfig),
	)
	return cmd
}
