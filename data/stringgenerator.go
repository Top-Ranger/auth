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

package data

import (
	"encoding/base64"
	"time"
)

// GetStrings returns a string representation of the id for verification. Please note: You have no access on the original id.
// See Get for more information.
//
// Can be used concurrent.
func GetStrings(data string) (id string, err error) {
	i, e := Get([]byte(data))
	if e != nil {
		err = e
		return
	}
	id = base64.StdEncoding.EncodeToString(i)
	return
}

// VerifyStrings verifies a an id / data combination.
// See Verify for more information about captchas.
//
// Can be used concurrent.
func VerifyStrings(id, data string) bool {
	i, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return false
	}
	return Verify(i, []byte(data))
}

// GetStringsTimed returns a timed string representation of the id for verification. Please note: You have no access on the original id.
// See Get for more information about captchas.
//
// Can be used concurrent.
func GetStringsTimed(start time.Time, data string) (id string, err error) {
	i, e := GetTimed(start, []byte(data))
	if e != nil {
		err = e
		return
	}
	id = base64.StdEncoding.EncodeToString(i)
	return
}

// VerifyStringsTimed verifies a timed id / data combination.
// See Verify for more information about captchas.
//
// Can be used concurrent.
func VerifyStringsTimed(id, data string, now time.Time, validDuration time.Duration) bool {
	i, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return false
	}
	return VerifyTimed(i, []byte(data), now, validDuration)
}
