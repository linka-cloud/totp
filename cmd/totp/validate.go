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
)

var (
	validateCmd = &cobra.Command{
		Use:     "validate",
		Short:   "Validates configured TOTP accounts",
		Aliases: []string{"v", "check"},
		Run: func(cmd *cobra.Command, args []string) {
			as, err := loadFile(configPath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			code := 0
			for _, v := range as {
				if err := v.ValidateConfig(); err != nil {
					fmt.Printf("%s: %v\n", v.GetName(), err)
					code = 1
				} else if !quiet {
					fmt.Printf("%s: valid\n", v.GetName())
				}
			}
			os.Exit(code)
		},
	}
)

func init() {
	rootCmd.AddCommand(validateCmd)
}
