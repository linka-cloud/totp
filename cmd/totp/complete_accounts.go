// Copyright 2021 Linka Cloud  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"strings"

	"github.com/spf13/cobra"
)

func completeAccounts(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	as, err := loadFile(configPath)
	if err != nil {
		return nil, cobra.ShellCompDirectiveDefault
	}
	var out []string
	for _, v := range as {
		if strings.HasPrefix(strings.ToLower(v.GetName()), strings.ToLower(toComplete)) {
			out = append(out, v.GetName())
		}
	}
	return out, cobra.ShellCompDirectiveDefault
}
