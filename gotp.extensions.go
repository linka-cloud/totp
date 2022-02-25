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
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"image"
	"net/url"
	"strconv"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"
	"google.golang.org/protobuf/proto"
)

const (
	Period = 30
)

type options struct {
	alg    *Algorithm
	size   uint
	digit  *Digit
	period *uint
}

type Option func(o *options)

func WithDigits(d Digit) Option {
	return func(o *options) {
		o.digit = &d
	}
}

func WithAlgorithm(al Algorithm) Option {
	return func(o *options) {
		switch al {
		case AlgorithmSHA256, AlgorithmSHA512, AlgorithmMD5:
			o.alg = &al
		default:
			v := AlgorithmSHA1
			o.alg = &v
		}
	}
}

func WithSecretSize(size uint) Option {
	return func(o *options) {
		o.size = size
	}
}

func NewOTPAccount(issuer, name string, opts ...Option) (*OTPAccount, error) {
	o := &options{}
	for _, fn := range opts {
		fn(o)
	}
	if issuer == "" {
		return nil, otp.ErrGenerateMissingIssuer
	}
	if name == "" {
		return nil, otp.ErrGenerateMissingAccountName
	}

	if o.size == 0 {
		o.size = 20
	}
	secret := make([]byte, o.size)
	if _, err := rand.Reader.Read(secret); err != nil {
		return nil, err
	}
	a := &OTPAccount{
		Name:      proto.String(name),
		Issuer:    proto.String(issuer),
		Algorithm: o.alg,
		Secret:    secret,
		Digits:    o.digit,
	}
	a.Default()
	return a, nil
}

func (x *OTPAccount) URL() string {
	v := url.Values{}
	v.Set("secret", base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(x.Secret))
	v.Set("issuer", x.GetIssuer())
	v.Set("period", strconv.FormatUint(uint64(Period), 10))
	v.Set("algorithm", x.GetAlgorithm().String())
	switch x.GetDigits() {
	case DigitEight:
		v.Set("digits", "8")
	default:
		v.Set("digits", "6")
	}
	u := url.URL{
		Scheme:   "otpauth",
		Host:     "totp",
		Path:     "/" + x.GetIssuer() + ":" + x.GetName(),
		RawQuery: v.Encode(),
	}
	return u.String()
}

func (x *OTPAccount) Image(width int, height int) (image.Image, error) {
	b, err := qr.Encode(x.URL(), qr.M, qr.Auto)
	if err != nil {
		return nil, err
	}
	b, err = barcode.Scale(b, width, height)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (x *OTPAccount) opts() totp.ValidateOpts {
	opts := totp.ValidateOpts{}
	switch x.GetDigits() {
	case DigitEight:
		opts.Digits = otp.DigitsEight
	default:
		opts.Digits = otp.DigitsSix
	}
	switch x.GetAlgorithm() {
	case AlgorithmSHA256:
		opts.Algorithm = otp.AlgorithmSHA256
	case AlgorithmSHA512:
		opts.Algorithm = otp.AlgorithmSHA512
	case AlgorithmMD5:
		opts.Algorithm = otp.AlgorithmMD5
	default:
		opts.Algorithm = otp.AlgorithmSHA1
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
