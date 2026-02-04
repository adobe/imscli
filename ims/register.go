package ims

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

func RegisterCmd(imsConfig *Config) *cobra.Command {
	var registerURL, clientName string
	var redirectURIs []string

	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a client.",
		Long:  `Register a new OAuth client using Dynamic Client Registration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Build redirect URIs JSON array
			redirectURIsJSON := "["
			for i, uri := range redirectURIs {
				if i > 0 {
					redirectURIsJSON += ","
				}
				redirectURIsJSON += fmt.Sprintf(`"%s"`, uri)
			}
			redirectURIsJSON += "]"

			payload := strings.NewReader(fmt.Sprintf(`{
  "client_name": "%s",
  "redirect_uris": %s
}`, clientName, redirectURIsJSON))
			req, _ := http.NewRequest("POST", registerURL, payload)
			req.Header.Add("content-type", "application/json")
			res, _ := http.DefaultClient.Do(req)
			defer res.Body.Close()
			body, _ := io.ReadAll(res.Body)
			fmt.Println(res)
			fmt.Println(string(body))
			return nil
		},
	}

	cmd.Flags().StringVarP(&registerURL, "url", "u", "", "registration endpoint URL.")
	cmd.Flags().StringVarP(&clientName, "clientName", "n", "", "Client application name.")
	cmd.Flags().StringSliceVarP(&redirectURIs, "redirectURIs", "r", []string{}, "Redirect URIs (comma-separated or multiple flags).")

	return cmd
}
