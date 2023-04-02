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
	"fmt"
	"os"

	"google.golang.org/protobuf/proto"

	"go.linka.cloud/totp"
)

type fileStore struct {
	path string
}

func NewFileStore(path string) (Store, error) {
	return &fileStore{path: path}, nil
}

func (f *fileStore) Load() ([]*totp.OTPAccount, error) {
	b, err := os.ReadFile(f.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", f.path, err)
	}
	return load(b)
}

func (f *fileStore) Save(as []*totp.OTPAccount) error {
	data := &totp.OTPData{
		OTPAccounts: as,
	}
	b, err := proto.Marshal(data)
	if err != nil {
		return fmt.Errorf("encode failed: %v", err)
	}
	tmp := f.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0700); err != nil {
		return fmt.Errorf("write config failed: %v", err)
	}
	if err := os.Rename(tmp, f.path); err != nil {
		return fmt.Errorf("write config failed: %v", err)
	}
	return nil
}

func (f *fileStore) Import(b []byte) error {
	if _, err := load(b); err != nil {
		return err
	}
	tmp := f.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0700); err != nil {
		return fmt.Errorf("write config failed: %v", err)
	}
	if err := os.Rename(tmp, f.path); err != nil {
		return fmt.Errorf("write config failed: %v", err)
	}
	return nil
}

func (f *fileStore) Dump() (string, error) {
	b, err := os.ReadFile(f.path)
	return string(b), err
}

func load(b []byte) ([]*totp.OTPAccount, error) {
	p := &totp.OTPData{}
	if err := proto.Unmarshal(b, p); err != nil {
		return nil, fmt.Errorf("proto decode: %w", err)
	}
	return p.OTPAccounts, nil
}
