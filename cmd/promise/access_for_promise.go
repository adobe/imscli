package promise

import (
	"encoding/json"
	"fmt"

	"github.com/adobe/imscli/ims"
	"github.com/spf13/cobra"
)

func AccessForPromiseCmd(imsConfig *ims.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "access-for-promise",
		Aliases: []string{"afp"},
		Short:   "Exchange an access token for a promise token.",
		Long:    `Perform the IMS Promise Token grant (grant_type=promise) to exchange an authenticating access token for a promise token.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			resp, err := imsConfig.PromiseTokenExchange()
			if err != nil {
				return fmt.Errorf("error during the promise token exchange: %w", err)
			}
			if imsConfig.FullOutput {
				data := struct {
					PromiseToken string `json:"promise_token"`
					TokenType    string `json:"token_type"`
					Scope        string `json:"scope"`
					ExpiresIn    int    `json:"expires_in"`
				}{resp.PromiseToken, resp.TokenType, resp.Scope, resp.ExpiresIn}
				jsonData, err := json.MarshalIndent(data, "", "  ")
				if err != nil {
					return fmt.Errorf("error marshalling full JSON response: %w", err)
				}
				fmt.Printf("%s\n", jsonData)
				return nil
			}
			fmt.Println(resp.PromiseToken)
			return nil
		},
	}

	cmd.Flags().StringVarP(&imsConfig.ClientID, "clientID", "c", "", "IMS client ID.")
	cmd.Flags().StringVarP(&imsConfig.ClientSecret, "clientSecret", "p", "", "IMS client secret.")
	cmd.Flags().StringVarP(&imsConfig.AccessToken, "accessToken", "t", "", "Authenticating access token (only access tokens are accepted).")
	cmd.Flags().StringVarP(&imsConfig.PromiseDefinitionID, "promiseDefinitionID", "d", "", "Promise definition ID.")
	cmd.Flags().StringSliceVarP(&imsConfig.Scopes, "scopes", "s", nil, "Scopes to request in the promise token.")
	cmd.Flags().BoolVarP(&imsConfig.FullOutput, "fullOutput", "F", false, "Output a JSON with the promise token, token type, scope, and expiration.")

	return cmd
}
