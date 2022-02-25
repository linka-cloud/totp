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
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"

	"go.linka.cloud/totp"
)

var (
	configPath string
	quiet      bool

	rootCmd = &cobra.Command{
		Use: "totp",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !strings.HasPrefix(configPath, "~/") {
				return nil
			}
			h, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("home: %w", err)
			}
			configPath = filepath.Join(h, strings.TrimPrefix(configPath, "~/"))
			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", envd("TOTP_CONFIG", "~/.config/totp/data"), "The path to the TOTP accounts configuration [$TOTP_CONFIG]")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Display only the code")
}

func envd(name, v string) string {
	if e := os.Getenv(name); e != "" {
		return e
	}
	return v
}

func loadFile(path string) ([]*totp.OTPAccount, error) {
	b, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return load(b)
}

func loadOrCreate(path string) ([]*totp.OTPAccount, error) {
	as, err := loadFile(path)
	if err == nil {
		return as, nil
	}
	if !os.IsNotExist(err) {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("create config file: %w", err)
	}
	return nil, f.Close()
}

func fromGoogleAuthenticatorMigration(data string) ([]*totp.OTPAccount, error) {
	q, err := url.QueryUnescape(data)
	if err != nil {
		return nil, fmt.Errorf("unescape data query string: %w", err)
	}
	b, err := base64.StdEncoding.DecodeString(q)
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %w", err)
	}
	return load(b)
}

func load(b []byte) ([]*totp.OTPAccount, error) {
	p := &totp.OTPData{}
	if err := proto.Unmarshal(b, p); err != nil {
		return nil, fmt.Errorf("proto decode: %w", err)
	}
	return p.OTPAccounts, nil
}

func save(path string, as []*totp.OTPAccount) {
	data := &totp.OTPData{
		OTPAccounts: as,
	}
	b, err := proto.Marshal(data)
	if err != nil {
		fmt.Printf("encode failed: %v\n", err)
		os.Exit(1)
	}
	tmp := path + ".tmp"
	if err := ioutil.WriteFile(tmp, b, 0700); err != nil {
		fmt.Printf("write config failed: %v\n", err)
		os.Exit(1)
	}
	if err := os.Rename(tmp, configPath); err != nil {
		fmt.Printf("write config failed: %v\n", err)
		os.Exit(1)
	}
}
