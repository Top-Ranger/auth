// SPDX-License-Identifier: Apache-2.0
// Copyright 2020 Marcus Soll
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package captcha

import (
	"encoding/base64"
	"time"
)

// GetStrings returns a string representation of a new captcha (with default size). Please note: You have no access on the original id.
// See Get for more information about captchas.
//
// Can be used concurrent.
func GetStrings() (id, captcha string, err error) {
	i, c, e := Get(RandomSizeDefault)
	if e != nil {
		err = e
		return
	}
	id = base64.StdEncoding.EncodeToString(i)
	captcha = base64.StdEncoding.EncodeToString(c)
	return
}

// VerifyStrings verifies a string representation of a new captcha (with default size).
// See Verify for more information about captchas.
//
// Can be used concurrent.
func VerifyStrings(id, captcha string) bool {
	i, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return false
	}
	c, err := base64.StdEncoding.DecodeString(captcha)
	if err != nil {
		return false
	}
	return Verify(i, c, RandomSizeDefault)
}

// GetStringsTimed returns a string representation of a new timed captcha (with default size). Please note: You have no access on the original id.
// See Get for more information about captchas.
//
// Can be used concurrent.
func GetStringsTimed(start time.Time) (id, captcha string, err error) {
	i, c, e := GetTimed(start, RandomSizeDefault)
	if e != nil {
		err = e
		return
	}
	id = base64.StdEncoding.EncodeToString(i)
	captcha = base64.StdEncoding.EncodeToString(c)
	return
}

// VerifyStringsTimed verifies a string representation of a new timed captcha (with default size).
// See Verify for more information about captchas.
//
// Can be used concurrent.
func VerifyStringsTimed(id, captcha string, now time.Time, validDuration time.Duration) bool {
	i, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return false
	}
	c, err := base64.StdEncoding.DecodeString(captcha)
	if err != nil {
		return false
	}
	return VerifyTimed(i, c, now, validDuration, RandomSizeDefault)
}
