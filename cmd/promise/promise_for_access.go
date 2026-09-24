package promise

import (
	"encoding/json"
	"fmt"

	"github.com/adobe/imscli/ims"
	"github.com/spf13/cobra"
)

func PromiseForAccessCmd(imsConfig *ims.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "promise-for-access",
		Aliases: []string{"pfa"},
		Short:   "Exchange a promise token for an access token.",
		Long:    `Perform the IMS Promise Exchange grant (grant_type=promise_exchange) to redeem a promise token for a fresh access token and a rotated promise token.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			resp, err := imsConfig.PromiseExchange()
			if err != nil {
				return fmt.Errorf("error during the promise exchange: %w", err)
			}
			if imsConfig.FullOutput {
				data := struct {
					AccessToken           string `json:"access_token"`
					PromiseToken          string `json:"promise_token"`
					PromiseTokenID        string `json:"promise_token_id"`
					TokenType             string `json:"token_type"`
					Scope                 string `json:"scope"`
					ExpiresIn             int    `json:"expires_in"`
					PromiseTokenExpiresIn int    `json:"promise_token_expires_in"`
				}{
					resp.AccessToken, resp.PromiseToken, resp.PromiseTokenID, resp.TokenType,
					resp.Scope, resp.ExpiresIn, resp.PromiseTokenExpiresIn,
				}
				jsonData, err := json.MarshalIndent(data, "", "  ")
				if err != nil {
					return fmt.Errorf("error marshalling full JSON response: %w", err)
				}
				fmt.Printf("%s\n", jsonData)
				return nil
			}
			fmt.Println(resp.AccessToken)
			return nil
		},
	}

	cmd.Flags().StringVarP(&imsConfig.ClientID, "clientID", "c", "", "IMS client ID.")
	cmd.Flags().StringVarP(&imsConfig.ClientSecret, "clientSecret", "p", "", "IMS client secret.")
	cmd.Flags().StringVarP(&imsConfig.PromiseToken, "promiseToken", "t", "", "Promise token to exchange.")
	cmd.Flags().StringSliceVarP(&imsConfig.Scopes, "scopes", "s", nil, "Optional scopes to request in the new access token.")
	cmd.Flags().BoolVarP(&imsConfig.FullOutput, "fullOutput", "F", false, "Output a JSON with the access token, rotated promise token, and their expirations.")

	return cmd
}
