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
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	recoverCmd = &cobra.Command{
		Use:   "recover",
		Short: "Recover a TOTP Keyring content",
		Long:  "Recovers is a set of commands to manually recover a TOTP Keyring content.",
	}

	recoverDumpCmd = &cobra.Command{
		Use:     "dump",
		Short:   "Dump the content of the keyring to stdout",
		Long:    "Dump the content of the keyring to stdout. This is useful to recover the content of the keyring in case of corrupted keyring content.",
		Aliases: []string{"export"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			s, err := store.Dump()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
			fmt.Println(s)
		},
	}

	recoverLoadCmd = &cobra.Command{
		Use:     "load [file]",
		Short:   "Load the content of the keyring from file or stdin",
		Long:    "Load the content of the keyring from file or stdin. This is useful to recover the content of the keyring in case of corrupted keyring content.",
		Aliases: []string{"import"},
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				b   []byte
				err error
			)
			if len(args) == 0 {
				b, _, err = bufio.NewReader(os.Stdin).ReadLine()
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: %v\n", err)
					os.Exit(1)
				}
			} else {
				b, err = os.ReadFile(args[0])
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: %v\n", err)
					os.Exit(1)
				}
			}
			if err := store.Import(b); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	recoverCmd.AddCommand(recoverDumpCmd, recoverLoadCmd)
	rootCmd.AddCommand(recoverCmd)
}
