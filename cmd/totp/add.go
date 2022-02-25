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
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"

	totp "go.linka.cloud/totp"
)

var (
	addCmd = &cobra.Command{
		Use:     "add [name] [secret]",
		Short:   "Add a TOTP accounts",
		Aliases: []string{"a", "new", "n", "create", "c"},
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			name, secret := strings.TrimSpace(args[0]), strings.TrimSpace(args[1])
			if name == "" || secret == "" {
				fmt.Println("invalid name or secret")
				os.Exit(1)
			}
			as, err := loadFile(configPath)
			if os.IsNotExist(err) {
				var f *os.File
				if f, err = os.Create(configPath); err == nil {
					err = f.Close()
				}
			}
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			for _, v := range as {
				if strings.EqualFold(v.GetName(), name) {
					fmt.Printf("%s: already exists\n", secret)
					os.Exit(1)
				}
			}
			b, err := base32.StdEncoding.DecodeString(secret)
			if err != nil {
				fmt.Printf("decode secret failed: %v\n", err)
				os.Exit(1)
			}
			typ := totp.OTPTypeTOTP
			a := &totp.OTPAccount{
				Name:   proto.String(name),
				Secret: b,
				Type:   &typ,
			}
			if err := a.ValidateConfig(); err != nil {
				fmt.Println("validate config failed: ", err)
				os.Exit(1)
			}
			save(append(as, a))
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
}
