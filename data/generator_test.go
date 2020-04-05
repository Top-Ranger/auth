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
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	data := []byte{24, 122, 5, 3}
	i, err := Get(data)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}
	if len(i) != hashSize {
		t.Errorf("i has wrong size (is: %d, should: %d)", len(i), hashSize)
	}

	i, err = Get(nil)
	if err != nil {
		t.Logf("error occured (nil): %s", err.Error())
		t.FailNow()
	}
	if len(i) != hashSize {
		t.Errorf("i has wrong size (is: %d, should: %d)", len(i), hashSize)
	}

	data = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	i, err = Get(data)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}
	if len(i) != hashSize {
		t.Errorf("i has wrong size (is: %d, should: %d)", len(i), hashSize)
	}
}

func TestVerify(t *testing.T) {
	data := []byte{24, 122, 5, 3}
	i, err := Get(data)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}

	// Correct
	r := Verify(i, data)
	if r == false {
		t.Error("verification failed")
	}

	// Nil
	r = Verify(i, nil)
	if r == true {
		t.Error("verification succeeded for nil")
	}

	// Wrong
	data = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	r = Verify(i, data)
	if r == true {
		t.Error("verification succeeded for wrong data")
	}

	r = Verify([]byte{0, 0, 0, 0, 0}, []byte{0, 0, 0, 0, 0})
	if r == true {
		t.Error("verification succeeded for wrong data (new id)")
	}

	// Mixed timed and not timed
	it, err := GetTimed(time.Now(), data)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}

	r = Verify(it, data)
	if r == true {
		t.Error("verification succeeded for wrong data (data is timed)")
	}
}

func TestGetTimed(t *testing.T) {
	data := []byte{24, 122, 5, 3}
	testtime := time.Now()
	i, err := GetTimed(testtime, data)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}
	if len(i) <= hashSize {
		t.Errorf("i has wrong size (is: %d, should be larger than: %d)", len(i), hashSize)
	}

	// Test different size
	data = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	i, err = GetTimed(testtime, data)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}
	if len(i) <= hashSize {
		t.Errorf("i has wrong size (is: %d, should be larger than: %d)", len(i), hashSize)
	}
}

func TestVerifyTimed(t *testing.T) {
	data := []byte{24, 122, 5, 3}
	testtime := time.Now()
	i, err := GetTimed(testtime, data)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}

	// Correct
	r := VerifyTimed(i, data, testtime, 1*time.Minute)
	if r == false {
		t.Error("verification failed")
	}

	// Correct later
	modifiedTime := testtime.Add(2 * time.Second)
	r = VerifyTimed(i, data, modifiedTime, 1*time.Minute)
	if r == false {
		t.Error("verification failed (modified time)")
	}

	// Too late
	modifiedTime = testtime.Add(2 * time.Minute)
	r = VerifyTimed(i, data, modifiedTime, 1*time.Minute)
	if r == true {
		t.Error("verification succeeded for wrong data (time over duration)")
	}

	modifiedTime = testtime.Add(-2 * time.Second)
	r = VerifyTimed(i, data, modifiedTime, 1*time.Minute)
	if r == true {
		t.Error("verification succeeded for wrong data (now before data time)")
	}

	// Wrong
	data = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	r = VerifyTimed(i, data, testtime, 1*time.Minute)
	if r == true {
		t.Error("verification succeeded for wrong data")
	}

	r = VerifyTimed([]byte{0, 0, 0, 0, 0}, []byte{0, 0, 0, 0, 0}, testtime, 1*time.Minute)
	if r == true {
		t.Error("verification succeeded for wrong data (new id)")
	}

	// Mixed timed and not timed
	in, err := Get(data)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}

	r = VerifyTimed(in, data, testtime, 1*time.Minute)
	if r == true {
		t.Error("verification succeeded for wrong data (data is timed)")
	}

	// Try to change the time - uses some internal knowledge
	_, realId := i[:len(i)-hashSize], i[len(i)-hashSize:]
	forgedTime, _ := testtime.Add(1 * time.Hour).MarshalBinary()
	forged := append(forgedTime, realId...)
	modifiedTime = testtime.Add(2 * time.Second)
	modifiedTime = testtime.Add(2 * time.Hour)
	r = VerifyTimed(forged, data, modifiedTime, 1*time.Minute)
	if r == true {
		t.Error("verification succeeded for wrong data (Modified timestamp)")
	}

	// invalid time
	brokenTime := []byte{5, 5, 5}
	forged = append(brokenTime, realId...)
	r = VerifyTimed(forged, data, testtime, 1*time.Minute)
	if r == true {
		t.Error("capcha verification succeeded for wrong data (invalid time)")
	}

	// even better forging - should not be able to occure in the wild
	hash := hashGenerator()
	hash.Write(randomData)
	hash.Write(data)
	hash.Write(brokenTime)
	forged = hash.Sum(brokenTime)
	r = VerifyTimed(forged, data, testtime, 1*time.Minute)
	if r == true {
		t.Error("capcha verification succeeded for wrong data (Modified timestamp)")
	}
}
