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
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	removeCmd = &cobra.Command{
		Use:               "remove [name]",
		Short:             "Remove a TOTP accounts",
		Aliases:           []string{"rm", "delete", "del"},
		ValidArgsFunction: completeAccounts,
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			as, err := loadFile(configPath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			index := -1
			for i, v := range as {
				if strings.EqualFold(v.GetName(), name) {
					index = i
				}
			}
			if index == -1 {
				fmt.Printf("%s: account not found\n", name)
				os.Exit(1)
			}
			as = append(as[:index], as[index+1:]...)
			save(as)
		},
	}
)

func init() {
	rootCmd.AddCommand(removeCmd)
}
