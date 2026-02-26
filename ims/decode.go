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
	"encoding/base64"
	"fmt"
	"strings"
)

// DecodedToken represents the decoded parts of a JWT token.
type DecodedToken struct {
	Header  string
	Payload string
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
		return nil, fmt.Errorf("incomplete parameters for token decodification: %w", err)
	}
	parts := strings.Split(i.Token, ".")

	if len(parts) != 3 {
		return nil, fmt.Errorf("the JWT is not composed by 3 parts")
	}

	decoded := &DecodedToken{}

	// Decode header
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("error decoding token header: %w", err)
	}
	decoded.Header = string(headerBytes)

	// Decode payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("error decoding token payload: %w", err)
	}
	decoded.Payload = string(payloadBytes)

	return decoded, nil
}
