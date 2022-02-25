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
	"github.com/spf13/cobra/doc"
)

var (
	mdOut = "docs"
	mdCmd = &cobra.Command{
		Use:    "doc",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			if err := os.MkdirAll(mdOut, os.ModePerm); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if err := doc.GenMarkdownTree(cmd.Root(), mdOut); err != nil {
				fmt.Println(err)
			}
		},
	}
)

func init() {
	mdCmd.Flags().StringVarP(&mdOut, "out", "o", "docs", "Markdown docs output directory")
	rootCmd.AddCommand(mdCmd)
}
