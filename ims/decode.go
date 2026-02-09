// Copyright 2021 Adobe. All rights reserved.
// This file is licensed to you under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License. You may obtain a copy
// of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under
// the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR REPRESENTATIONS
// OF ANY KIND, either express or implied. See the License for the specific language
// governing permissions and limitations under the License.

package ims

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// DecodedToken represents the decoded parts of a JWT token.
type DecodedToken struct {
	Header    string
	Payload   string
	Signature string
}

func (i Config) validateDecodeTokenConfig() error {

	if i.Token == "" {
		return fmt.Errorf("missing token parameter")
	}
	return nil
}

func (i Config) DecodeToken() (*DecodedToken, error) {
	err := i.validateDecodeTokenConfig()
	if err != nil {
		return nil, fmt.Errorf("incomplete parameters for token decodification: %v", err)
	}
	parts := strings.Split(i.Token, ".")

	if len(parts) != 3 {
		return nil, fmt.Errorf("the JWT is not composed by 3 parts")
	}

	// Decode header and payload (not signature since it's binary)
	decoded := &DecodedToken{
		Signature: parts[2],
	}

	// Decode and prettify header
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("error decoding token header: %v", err)
	}
	decoded.Header, err = prettyJSON(headerBytes)
	if err != nil {
		return nil, fmt.Errorf("error formatting token header: %v", err)
	}

	// Decode and prettify payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("error decoding token payload: %v", err)
	}
	decoded.Payload, err = prettyJSON(payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("error formatting token payload: %v", err)
	}

	return decoded, nil
}

// prettyJSON formats JSON bytes with indentation.
func prettyJSON(data []byte) (string, error) {
	var prettyBuf bytes.Buffer
	err := json.Indent(&prettyBuf, data, "", "  ")
	if err != nil {
		return "", err
	}
	return prettyBuf.String(), nil
}
