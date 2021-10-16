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

func (i Config) validateDecodeTokenConfig() error {

	if i.Token == "" {
		return fmt.Errorf("missing token parameter")
	}
	return nil
}

func (i Config) DecodeToken() ([]string, error) {
	err := i.validateDecodeTokenConfig()
	if err != nil {
		return nil, fmt.Errorf("incomplete parameters for token decodification: %v", err)
	}
	parts := strings.Split(i.Token, ".")

	if len(parts) != 3 {
		return nil, fmt.Errorf("the JWT is not composed by 3 parts")
	}

	// i<2 to not decode the signature since it is not encoded
	for i:=0; i<2; i++  {
		decodedPart, err := base64.RawURLEncoding.DecodeString(parts[i])
		if err != nil {
			return nil, fmt.Errorf("error decoding token, part %d: %v", i, err)
		}
		parts[i] = string(decodedPart)
	}
	return parts, nil
}