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

	store2 "go.linka.cloud/totp/pkg/store"
)

var (
	migrateCmd = &cobra.Command{
		Use:     "migrate",
		Short:   "Move configured TOTP accounts to keyring",
		Aliases: []string{"mig", "move"},
		Run: func(cmd *cobra.Command, args []string) {
			as, err := store.Load()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			s, err := store2.NewKeyRing(keyRingName)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if err := s.Save(as); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}
