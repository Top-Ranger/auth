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

// This file contains the basic generator.
// All other generators schould use this as a basis (or an other generator).

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"sync"
	"time"
)

const (
	// RandomSizeDefault contains the suggested default size for random data.
	RandomSizeDefault = 6
)

var (
	randomData               = []byte{}
	initialisationRandomData = sync.Once{}
	hashGenerator            = sha256.New
	hashSize                 = hashGenerator().Size()
)

// setRandomData sets the hidden random data. It should be called before generating the first captcha, and only once (since resetting makes all older captchas invalid).
// The generator functions do this automatically, so there is no need to call it manually.
func setRandomData() error {
	b := make([]byte, hashSize*2)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}

	randomData = b

	return nil
}

// Get returns one random id / captcha combination.
// The number of bytes is determined by randomSize.
//
// Can be used concurrent.
func Get(randomSize int) (id, captcha []byte, err error) {
	initialisationRandomData.Do(func() {
		setRandomData()
	})

	if randomSize < 1 {
		err = errors.New("randomSize must be positive")
		return
	}

	b := make([]byte, randomSize)
	_, err = rand.Read(b)
	if err != nil {
		return
	}
	captcha = b[:]
	hash := hmac.New(hashGenerator, randomData)
	hash.Write(captcha)
	id = hash.Sum(nil)
	return
}

// Verify validates whether an id / captia combination is valid.
// randomSize musst correspond to the captcha size and must be the same as at the generation.
//
// Since Verify does not check if an id is already used, the same id / captcha combination is always valid.
//
// Can be used concurrent.
func Verify(id, captcha []byte, randomSize int) bool {
	initialisationRandomData.Do(func() {
		setRandomData()
	})

	if randomSize < 1 {
		return false
	}

	if len(id) != hashSize {
		return false
	}
	if len(captcha) != randomSize {
		return false
	}

	hash := hmac.New(hashGenerator, randomData)
	hash.Write(captcha)
	checksum := hash.Sum(nil)
	return subtle.ConstantTimeCompare(checksum, id) == 1
}

// GetTimed returns one timed random id / captcha combination.
// The number of bytes is determined by randomSize.
// start determines the time from which the captcha is valid.
//
// You do not need to remember the time since it is encoded in the id (and can not be tampered with without invalidating the captcha).
//
// Can be used concurrent.
func GetTimed(start time.Time, randomSize int) (id, captcha []byte, err error) {
	initialisationRandomData.Do(func() {
		setRandomData()
	})

	if randomSize < 1 {
		err = errors.New("randomSize must be positive")
		return
	}

	b := make([]byte, randomSize)
	_, err = rand.Read(b)
	if err != nil {
		return
	}
	timeEncoded, err := start.GobEncode()
	if err != nil {
		return
	}
	captcha = b[:]
	hash := hmac.New(hashGenerator, randomData)
	hash.Write(captcha)
	hash.Write(timeEncoded)
	id = hash.Sum(timeEncoded)
	return
}

// VerifyTimed validates whether an id / captia combination is valid and in date.
// randomSize musst correspond to the captcha size and must be the same as at the generation.
// Duration determines how long a captcha should be seen as valid.
//
// Since VerifyTimed does not check if an id is already used, the same id / captcha combination is always valid (in the given time period).
//
// Can be used concurrent.
func VerifyTimed(id, captcha []byte, now time.Time, validDuration time.Duration, randomSize int) bool {
	initialisationRandomData.Do(func() {
		setRandomData()
	})

	if randomSize < 1 {
		return false
	}

	if len(id) <= hashSize {
		return false
	}
	if len(captcha) != randomSize {
		return false
	}

	timeEncoded := make([]byte, len(id)-hashSize)
	copy(timeEncoded, id[:len(id)-hashSize]) // We need a true copy here, or else subtle.ConstantTimeCompare returns always true

	hash := hmac.New(hashGenerator, randomData)
	hash.Write(captcha)
	hash.Write(timeEncoded)
	checksum := hash.Sum(timeEncoded)
	if subtle.ConstantTimeCompare(checksum, id) == 0 {
		return false
	}
	var t time.Time
	err := t.GobDecode(timeEncoded)
	if err != nil {
		return false
	}
	if now.Before(t) {
		return false
	}
	if now.Sub(t) > validDuration {
		return false
	}
	return true
}
