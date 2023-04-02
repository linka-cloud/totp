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

package store

import (
	"github.com/zalando/go-keyring"
	"google.golang.org/protobuf/encoding/prototext"

	"go.linka.cloud/totp"
)

const (
	key = "accounts"
)

type keyRing struct {
	name string
}

func NewKeyRing(name string) (Store, error) {
	return &keyRing{name: name}, nil
}

func (k *keyRing) Load() ([]*totp.OTPAccount, error) {
	v, err := keyring.Get(k.name, key)
	if err != nil {
		return nil, err
	}
	var data totp.OTPData
	if err := prototext.Unmarshal([]byte(v), &data); err != nil {
		return nil, err
	}
	return data.OTPAccounts, nil
}

func (k *keyRing) Save(accounts []*totp.OTPAccount) error {
	b, err := prototext.Marshal(&totp.OTPData{OTPAccounts: accounts})
	if err != nil {
		return err
	}
	return keyring.Set(k.name, key, string(b))
}

func (k *keyRing) Import(v []byte) error {
	var data totp.OTPData
	if err := prototext.Unmarshal(v, &data); err != nil {
		return err
	}
	return keyring.Set(k.name, key, string(v))
}

func (k *keyRing) Dump() (string, error) {
	return keyring.Get(k.name, key)
}
