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

	"github.com/spf13/cobra"

	totp "go.linka.cloud/totp"
)

var (
	detailedCode bool

	codeCmd = &cobra.Command{
		Use:               "code [account]",
		Short:             "Generates a TOTP code for the account",
		Aliases:           []string{"get", "code"},
		ValidArgsFunction: completeAccounts,
		Args:              cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			as, err := store.Load()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			var a *totp.OTPAccount
			for _, v := range as {
				if v.GetName() == args[0] {
					a = v
					break
				}
			}
			if a == nil {
				fmt.Printf("account '%s' not found\n", args[0])
				os.Exit(1)
			}
			c, err := a.Generate()
			if err != nil {
				fmt.Printf("failed to generate code: %v\n", err)
			}
			if !detailedCode {
				fmt.Printf(c)
				return
			}
			fmt.Println(c, " ", a.ValidFor())
		},
	}
)

func init() {
	codeCmd.Flags().BoolVarP(&detailedCode, "details", "d", false, "Show code with validity duration")
	rootCmd.AddCommand(codeCmd)
}
