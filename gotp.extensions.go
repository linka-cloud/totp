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

package totp

import (
	"encoding/base32"
	"fmt"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"
)

const (
	Period = 30
)

func (x *OTPAccount) opts() totp.ValidateOpts {
	opts := totp.ValidateOpts{}
	switch x.GetDigits() {
	case DigitSix:
		opts.Digits = otp.DigitsSix
	case DigitEight:
		opts.Digits = otp.DigitsEight
	}
	switch x.GetAlgorithm() {
	case AlgorithmSHA1:
		opts.Algorithm = otp.AlgorithmSHA1
	case AlgorithmSHA256:
		opts.Algorithm = otp.AlgorithmSHA256
	case AlgorithmSHA512:
		opts.Algorithm = otp.AlgorithmSHA512
	case AlgorithmMD5:
		opts.Algorithm = otp.AlgorithmMD5
	}
	return opts
}

func (x *OTPAccount) Generate() (string, error) {
	opts := x.opts()
	switch x.GetType() {
	case OTPTypeHOTP:
		return hotp.GenerateCodeCustom(x.secret(), uint64(x.GetCounter()), hotp.ValidateOpts{
			Digits:    opts.Digits,
			Algorithm: opts.Algorithm,
		})
	default:
		return totp.GenerateCodeCustom(x.secret(), time.Now(), opts)
	}
}

func (x *OTPAccount) Validate(code string) bool {
	switch x.GetType() {
	case OTPTypeHOTP:
		return hotp.Validate(code, uint64(x.GetCounter()), x.secret())
	default:
		return totp.Validate(code, x.secret())
	}
}

func (x *OTPAccount) ValidFor() time.Duration {
	return time.Duration(Period-time.Now().Second()%Period) * time.Second
}

func (x *OTPAccount) secret() string {
	return base32.StdEncoding.EncodeToString(x.Secret)
}

func (x *OTPAccount) ValidateConfig() error {
	c, err := totp.GenerateCodeCustom(x.secret(), time.Now(), x.opts())
	if err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	if !totp.Validate(c, x.secret()) {
		return fmt.Errorf("totp validation failed")
	}
	return nil
}
