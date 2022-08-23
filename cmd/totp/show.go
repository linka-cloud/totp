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
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"go.linka.cloud/totp"
)

var (
	showOut = ""

	showCmd = &cobra.Command{
		Use:               "show [name]",
		Short:             "Create QRCode for TOTP account",
		Aliases:           []string{"qr", "qrcode"},
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completeAccounts,
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			if dumpOut == "" {
				dumpOut = filepath.Join(os.TempDir(), "totp")
			}
			if err := os.MkdirAll(dumpOut, 0700); err != nil {
				fmt.Println("create directory: ", err)
				os.Exit(1)
			}
			as, err := store.Load()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			var a *totp.OTPAccount
			for _, v := range as {
				if strings.EqualFold(v.GetName(), name) {
					a = v
					break
				}
			}
			if a == nil {
				fmt.Printf("%s: account not found\n", name)
				os.Exit(1)
			}
			saveAsQRCode(filepath.Join(dumpOut, a.GetName()+".png"), a)
		},
	}
)

func init() {
	showCmd.Flags().StringVarP(&showOut, "out", "o", "", "The qrcode image output directory")
	rootCmd.AddCommand(showCmd)
}
