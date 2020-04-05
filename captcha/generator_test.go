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
	"bytes"
	"crypto/rand"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	i, c, err := Get(RandomSizeDefault)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}
	if len(c) != RandomSizeDefault {
		t.Errorf("c has wrong size (is: %d, should: %d)", len(c), RandomSizeDefault)
	}
	if len(i) != hashSize {
		t.Errorf("i has wrong size (is: %d, should: %d)", len(i), hashSize)
	}

	// simple random test
	for x := 0; x < 100; x++ {
		i2, c2, err := Get(RandomSizeDefault)

		if err != nil {
			t.Logf("error occured: %s", err.Error())
			t.FailNow()
		}
		if bytes.Compare(i, i2) == 0 {
			t.Errorf("simple random error failed: %s == %s (i)", i, i2)
			break
		}
		if bytes.Compare(c, c2) == 0 {
			t.Errorf("simple random error failed: %s == %s (c)", c, c2)
			break
		}
	}

	// Test different random size
	i, c, err = Get(1000)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}
	if len(c) != 1000 {
		t.Errorf("c has wrong size (is: %d, should: %d)", len(c), 1000)
	}
	if len(i) != hashSize {
		t.Errorf("i has wrong size (is: %d, should: %d)", len(i), hashSize)
	}

	// Test negative size
	_, _, err = Get(-1000)
	if err == nil {
		t.Errorf("Generating negative size does not show an error")
	}

	_, _, err = Get(0)
	if err == nil {
		t.Errorf("Generating zero size does not show an error")
	}
}

func TestVerify(t *testing.T) {
	i, c, err := Get(RandomSizeDefault)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}

	// Correct captcha
	r := Verify(i, c, RandomSizeDefault)
	if r == false {
		t.Error("capcha verification failed")
	}

	// Wrong capcha - most likely
	rand.Read(c)
	r = Verify(i, c, RandomSizeDefault)
	if r == true {
		t.Error("capcha verification succeeded for wrong captcha")
	}

	rand.Read(c)
	rand.Read(i)
	r = Verify(i, c, RandomSizeDefault)
	if r == true {
		t.Error("capcha verification succeeded for wrong captcha (new id)")
	}

	// Test wrong sizes
	iNew := make([]byte, 55)
	rand.Read(iNew)
	r = Verify(iNew, c, RandomSizeDefault)
	if r == true {
		t.Error("capcha verification succeeded for wrong captcha (i wrong size)")
	}

	cNew := make([]byte, 55)
	rand.Read(cNew)
	r = Verify(i, cNew, RandomSizeDefault)
	if r == true {
		t.Error("capcha verification succeeded for wrong captcha (c wrong size)")
	}

	// Mixed timed and not timed
	it, ct, err := GetTimed(time.Now(), RandomSizeDefault)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}

	r = Verify(it, ct, RandomSizeDefault)
	if r == true {
		t.Error("capcha verification succeeded for wrong captcha (captcha is timed)")
	}

	// Test different size
	i, c, err = Get(1000)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}

	// Correct captcha
	r = Verify(i, c, 1000)
	if r == false {
		t.Error("capcha verification failed (size 1000)")
	}

	r = Verify(i, c, RandomSizeDefault)
	if r == true {
		t.Error("capcha verification failed (size 1000, tested default)")
	}

	// Test negative size
	r = Verify(i, c, -1000)
	if r == true {
		t.Error("capcha verification worked with negative size")
	}

	r = Verify(i, c, 0)
	if r == true {
		t.Error("capcha verification worked with zero size")
	}
}

func TestGetTimed(t *testing.T) {
	testtime := time.Now()
	i, c, err := GetTimed(testtime, RandomSizeDefault)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}
	if len(c) != RandomSizeDefault {
		t.Errorf("c has wrong size (is: %d, should: %d)", len(c), RandomSizeDefault)
	}
	if len(i) <= hashSize {
		t.Errorf("i has wrong size (is: %d, should be larger than: %d)", len(i), hashSize)
	}

	// simple random test
	for x := 0; x < 100; x++ {
		i2, c2, err := GetTimed(testtime, RandomSizeDefault)

		if err != nil {
			t.Logf("error occured: %s", err.Error())
			t.FailNow()
		}
		if bytes.Compare(i, i2) == 0 {
			t.Errorf("simple random error failed: %s == %s (i)", i, i2)
			break
		}
		if bytes.Compare(c, c2) == 0 {
			t.Errorf("simple random error failed: %s == %s (c)", c, c2)
			break
		}
	}

	// Test different size
	i, c, err = GetTimed(testtime, 1000)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}
	if len(c) != 1000 {
		t.Errorf("c has wrong size (is: %d, should: %d)", len(c), 1000)
	}
	if len(i) <= hashSize {
		t.Errorf("i has wrong size (is: %d, should be larger than: %d)", len(i), hashSize)
	}

	// Test negative size
	_, _, err = GetTimed(testtime, -1000)
	if err == nil {
		t.Errorf("Generating negative size does not show an error")
	}

	_, _, err = GetTimed(testtime, 0)
	if err == nil {
		t.Errorf("Generating zero size does not show an error")
	}
}

func TestVerifyTimed(t *testing.T) {
	testtime := time.Now()
	i, c, err := GetTimed(testtime, RandomSizeDefault)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}

	// Correct captcha
	r := VerifyTimed(i, c, testtime, 1*time.Minute, RandomSizeDefault)
	if r == false {
		t.Error("capcha verification failed")
	}

	// Correct later
	modifiedTime := testtime.Add(2 * time.Second)
	r = VerifyTimed(i, c, modifiedTime, 1*time.Minute, RandomSizeDefault)
	if r == false {
		t.Error("capcha verification failed (modified time)")
	}

	// Too late
	modifiedTime = testtime.Add(2 * time.Minute)
	r = VerifyTimed(i, c, modifiedTime, 1*time.Minute, RandomSizeDefault)
	if r == true {
		t.Error("capcha verification succeeded for wrong captcha (time over duration)")
	}

	modifiedTime = testtime.Add(-2 * time.Second)
	r = VerifyTimed(i, c, modifiedTime, 1*time.Minute, RandomSizeDefault)
	if r == true {
		t.Error("capcha verification succeeded for wrong captcha (now before captcha time)")
	}

	// Wrong capcha - most likely
	rand.Read(c)
	r = VerifyTimed(i, c, testtime, 1*time.Minute, RandomSizeDefault)
	if r == true {
		t.Error("capcha verification succeeded for wrong captcha")
	}

	rand.Read(c)
	rand.Read(i)
	r = VerifyTimed(i, c, testtime, 1*time.Minute, RandomSizeDefault)
	if r == true {
		t.Error("capcha verification succeeded for wrong captcha (new id)")
	}

	// Test wrong sizes
	iNew := make([]byte, 55)
	rand.Read(iNew)
	r = VerifyTimed(iNew, c, testtime, 1*time.Minute, RandomSizeDefault)
	if r == true {
		t.Error("capcha verification succeeded for wrong captcha (i wrong size)")
	}

	cNew := make([]byte, 1)
	rand.Read(cNew)
	r = VerifyTimed(i, cNew, testtime, 1*time.Minute, RandomSizeDefault)
	if r == true {
		t.Error("capcha verification succeeded for wrong captcha (c wrong size)")
	}

	// Mixed timed and not timed
	in, cn, err := Get(RandomSizeDefault)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}

	r = VerifyTimed(in, cn, testtime, 1*time.Minute, RandomSizeDefault)
	if r == true {
		t.Error("capcha verification succeeded for wrong captcha (captcha is timed)")
	}

	// Try to change the time - uses some internal knowledge
	_, realId := i[:len(i)-hashSize], i[len(i)-hashSize:]
	forgedTime, _ := testtime.Add(1 * time.Hour).MarshalBinary()
	forged := append(forgedTime, realId...)
	modifiedTime = testtime.Add(2 * time.Second)
	modifiedTime = testtime.Add(2 * time.Hour)
	r = VerifyTimed(forged, c, modifiedTime, 1*time.Minute, RandomSizeDefault)
	if r == true {
		t.Error("capcha verification succeeded for wrong captcha (Modified timestamp)")
	}

	// invalid time
	brokenTime := []byte{5, 5, 5}
	forged = append(brokenTime, realId...)
	r = VerifyTimed(forged, c, testtime, 1*time.Minute, RandomSizeDefault)
	if r == true {
		t.Error("capcha verification succeeded for wrong captcha (invalid time)")
	}

	// even better forging - should not be able to occure in the wild
	hash := hashGenerator()
	hash.Write(randomData)
	hash.Write(c)
	hash.Write(brokenTime)
	forged = hash.Sum(brokenTime)
	r = VerifyTimed(forged, c, testtime, 1*time.Minute, RandomSizeDefault)
	if r == true {
		t.Error("capcha verification succeeded for wrong captcha (Modified timestamp)")
	}

	// Different size
	i, c, err = GetTimed(testtime, 1000)
	if err != nil {
		t.Logf("error occured: %s", err.Error())
		t.FailNow()
	}

	// Correct captcha
	r = VerifyTimed(i, c, testtime, 1*time.Minute, 1000)
	if r == false {
		t.Error("capcha verification failed (size 1000)")
	}

	modifiedTime = testtime.Add(2 * time.Second)
	r = VerifyTimed(i, c, testtime, 1*time.Minute, 1000)
	if r == false {
		t.Error("capcha verification failed (size 1000, modified time)")
	}

	modifiedTime = testtime.Add(2 * time.Hour)
	r = VerifyTimed(i, c, testtime, 1*time.Minute, 1000)
	if r == false {
		t.Error("capcha verification succeeded for wrong captcha (size 1000, modified time)")
	}

	// Test negative size
	r = VerifyTimed(i, c, testtime, 1*time.Minute, -1000)
	if r == true {
		t.Error("capcha verification with negative size")
	}

	r = VerifyTimed(i, c, testtime, 1*time.Minute, 0)
	if r == true {
		t.Error("capcha verification with zero size")
	}
}
