package ims

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// RegisterResponse contains the response from DCR registration
type RegisterResponse struct {
	StatusCode int
	Body       string
}

func (i Config) validateRegisterConfig() error {
	switch {
	case i.URL == "":
		return fmt.Errorf("missing IMS base URL parameter")
	case i.ClientName == "":
		return fmt.Errorf("missing client name parameter")
	case len(i.RedirectURIs) == 0:
		return fmt.Errorf("missing redirect URIs parameter")
	default:
		return nil
	}
}

// Register performs Dynamic Client Registration
func (i Config) Register() (RegisterResponse, error) {
	if err := i.validateRegisterConfig(); err != nil {
		return RegisterResponse{}, fmt.Errorf("invalid parameters for client registration: %v", err)
	}

	// Build redirect URIs JSON array
	redirectURIsJSON := "["
	for idx, uri := range i.RedirectURIs {
		if idx > 0 {
			redirectURIsJSON += ","
		}
		redirectURIsJSON += fmt.Sprintf(`"%s"`, uri)
	}
	redirectURIsJSON += "]"

	payload := strings.NewReader(fmt.Sprintf(`{
  "client_name": "%s",
  "redirect_uris": %s
}`, i.ClientName, redirectURIsJSON))

	endpoint := strings.TrimRight(i.URL, "/") + "/ims/register"

	req, err := http.NewRequest("POST", endpoint, payload)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("error making registration request: %v", err)
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("error reading response body: %v", err)
	}

	return RegisterResponse{
		StatusCode: res.StatusCode,
		Body:       string(body),
	}, nil
}
