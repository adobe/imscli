// Copyright 2026 Adobe. All rights reserved.
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
	"fmt"
	"os"

	"github.com/pkg/browser"
)

// openBrowser opens the given URL in the system default browser. Temporarily
// mutes browser.Stdout to suppress the "Opening in existing browser session"
// messages that some chromium-based browsers emit; the CLI's token output goes
// to stdout, so stray browser chatter would corrupt piped or scripted output.
// On failure, prints a fallback instruction to stderr and returns — callers
// continue (the user can open the URL manually).
func openBrowser(url string) {
	origStdout := browser.Stdout
	browser.Stdout = nil
	err := browser.OpenURL(url)
	browser.Stdout = origStdout
	if err != nil {
		fmt.Fprintf(os.Stderr, "error launching the browser, open it and visit %s\n", url)
	}
}
