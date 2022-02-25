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

	"github.com/spf13/cobra"
)

var (
	dumpOut = ""

	dumpCmd = &cobra.Command{
		Use:   "dump",
		Short: "Dump configured TOTP accounts to qrcode images",
		Run: func(cmd *cobra.Command, args []string) {
			if dumpOut == "" {
				dumpOut = filepath.Join(os.TempDir(), "totp")
			}
			if err := os.MkdirAll(dumpOut, 0700); err != nil {
				fmt.Println("create directory: ", err)
				os.Exit(1)
			}
			as, err := loadFile(configPath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			for _, v := range as {
				saveAsQRCode(filepath.Join(dumpOut, v.GetName()+".png"), v)
			}
		},
	}
)

func init() {
	dumpCmd.Flags().StringVarP(&dumpOut, "out", "o", ".", "The qrcode images output directory")
	rootCmd.AddCommand(dumpCmd)
}
