package cmd

import (
	"io/ioutil"
	"log"

	"github.com/adobe/imscli/ims"
	"github.com/spf13/cobra"
)

func RootCmd(version string) *cobra.Command {
	var verbose bool = false
	var configFile string
	var imsConfig = &ims.Config{}

	cmd := &cobra.Command{
		Use:     "imscli",
		Short:   "imscli is a tool to interact with Adobe IMS",
		Long:    `imscli is a CLI tool to automate and troubleshoot Adobe's authentication and authorization service IMS.`,
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !verbose {
				log.SetOutput(ioutil.Discard)
			}
			// This call of the initParams will load all env vars, config file and flags.
			return initParams(cmd, imsConfig, configFile)
		},
	}
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output.")
	cmd.PersistentFlags().StringVarP(&imsConfig.URL, "url", "U", "https://ims-na1.adobelogin.com",
		"IMS Endpoint URL.")
	cmd.PersistentFlags().StringVarP(&imsConfig.ProxyURL, "proxyUrl", "P", "",
		"Connect to IMS through the specified proxy. Specified as http(s)://host:port.")
	cmd.PersistentFlags().BoolVarP(&imsConfig.ProxyIgnoreTLS, "proxyIgnoreTLS", "T", false,
		"Ignore TLS certificate verification (only valid when connecting through a proxy).")
	cmd.PersistentFlags().StringVarP(&configFile, "configFile", "f", "", "Configuration file.")

	cmd.AddCommand(
		authzCmd(imsConfig),
		profileCmd(imsConfig),
		organizationsCmd(imsConfig),
		validateCmd(imsConfig),
		exchangeCmd(imsConfig),
		invalidateCmd(imsConfig),
		decodeCmd(imsConfig),
		refreshCmd(imsConfig),
		adminCmd(imsConfig),
		dcrCmd(imsConfig),
	)
	return cmd
}
