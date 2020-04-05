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
	"testing"
	"time"
)

func TestGetStrings(t *testing.T) {
	i, c, err := GetStrings()
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}
	if c == "" {
		t.Error("c ist empty string")
	}
	if i == "" {
		t.Error("i ist empty string")
	}
}

func TestVerifyStrings(t *testing.T) {
	i, c, err := GetStrings()
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}
	if !VerifyStrings(i, c) {
		t.Error("Verification failed")
	}

	if VerifyStrings("äää", c) {
		t.Error("Verification failed (invalid id)")
	}

	if VerifyStrings(i, "äää") {
		t.Error("Verification failed (invalid captcha)")
	}
}

func TestGetStringsTimed(t *testing.T) {
	testtime := time.Now()
	i, c, err := GetStringsTimed(testtime)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}
	if c == "" {
		t.Error("c ist empty string")
	}
	if i == "" {
		t.Error("i ist empty string")
	}
}

func TestVerifyStringsTimed(t *testing.T) {
	testtime := time.Now()
	i, c, err := GetStringsTimed(testtime)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}
	if !VerifyStringsTimed(i, c, testtime.Add(2*time.Second), 1*time.Minute) {
		t.Error("Verification failed")
	}

	if VerifyStringsTimed("äää", c, testtime.Add(2*time.Second), 1*time.Minute) {
		t.Error("Verification failed (invalid id)")
	}

	if VerifyStringsTimed(i, "äää", testtime.Add(2*time.Second), 1*time.Minute) {
		t.Error("Verification failed (invalid captcha)")
	}
}
