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
	"encoding/base32"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"go.linka.cloud/totp"
)

var (
	generateAsURL      bool
	generateQRCodePath string

	generateCmd = &cobra.Command{
		Use:     "generate [issuer] [account]",
		Short:   "Generate a new TOTP account",
		Aliases: []string{"new", "make", "gen"},
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			a, err := totp.NewOTPAccount(args[0], args[1])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			switch {
			case generateQRCodePath != "":
				fmt.Println(a.URL())
				saveAsQRCode(generateQRCodePath, a)
			case generateAsURL:
				fmt.Println(a.URL())
			default:
				fmt.Println(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(a.Secret))
			}
		},
	}
)

func init() {
	generateCmd.Flags().StringVar(&generateQRCodePath, "qrcode", "", "Image path for the generated QRCode Image")
	generateCmd.Flags().BoolVarP(&generateAsURL, "url", "u", false, "Generate account as URL")
	rootCmd.AddCommand(generateCmd)
}
