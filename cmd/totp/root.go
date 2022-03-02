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
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"

	"go.linka.cloud/totp"
	store2 "go.linka.cloud/totp/pkg/store"
)

const (
	keyRingName = "totp"

	// defaultConfigPath = "~/.config/totp/data"
	defaultConfigPath = ""
)

var (
	configPath string
	useKeyRing = true

	store store2.Store

	rootCmd = &cobra.Command{
		Use: "totp",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(configPath) == "" {
				return initKeyRing()
			}
			if !strings.HasPrefix(configPath, "~/") {
				return initFileStore()
			}
			h, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("home: %w", err)
			}
			configPath = filepath.Join(h, strings.TrimPrefix(configPath, "~/"))
			return initFileStore()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", envd("TOTP_CONFIG", defaultConfigPath), "The path to the TOTP accounts configuration [$TOTP_CONFIG]")
}

func envd(name, v string) string {
	if e := os.Getenv(name); e != "" {
		return e
	}
	return v
}

func initKeyRing() (err error) {
	store, err = store2.NewKeyRing(keyRingName)
	return
}

func initFileStore() (err error) {
	store, err = store2.NewFileStore(configPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
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

func saveAsQRCode(path string, a *totp.OTPAccount) {
	img, err := a.Image(200, 200)
	if err != nil {
		fmt.Println("failed to generate qrcode: ", err)
		os.Exit(1)
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		fmt.Println("failed to encode qrcode: ", err)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(path, buf.Bytes(), 0700); err != nil {
		fmt.Println("failed to write qrcode: ", err)
		os.Exit(1)
	}
}
