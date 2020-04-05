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

// This file contains the basic generator.

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"sync"
	"time"
)

var (
	randomData               = []byte{}
	initialisationRandomData = sync.Once{}
	hashGenerator            = sha256.New
	hashSize                 = hashGenerator().Size()
)

// setRandomData sets the hidden random data. It should be called before generating the first id, and only once (since resetting makes all older captchas invalid).
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
//
// Can be used concurrent.
func Get(data []byte) (id []byte, err error) {
	initialisationRandomData.Do(func() {
		setRandomData()
	})

	hash := hmac.New(hashGenerator, randomData)
	hash.Write(data)
	id = hash.Sum(nil)
	return
}

// Verify validates whether an id / data combination is valid.
//
// Can be used concurrent.
func Verify(id, data []byte) bool {
	initialisationRandomData.Do(func() {
		setRandomData()
	})

	if len(id) != hashSize {
		return false
	}

	hash := hmac.New(hashGenerator, randomData)
	hash.Write(data)
	checksum := hash.Sum(nil)
	return subtle.ConstantTimeCompare(checksum, id) == 1
}

// GetTimed returns one timed random id / data combination.
// start determines the time from which the authentification is valid.
//
// You do not need to remember the time since it is encoded in the id (and can not be tampered with without invalidating the authentification).
//
// Can be used concurrent.
func GetTimed(start time.Time, data []byte) (id []byte, err error) {
	initialisationRandomData.Do(func() {
		setRandomData()
	})

	timeEncoded, err := start.GobEncode()
	if err != nil {
		return
	}
	hash := hmac.New(hashGenerator, randomData)
	hash.Write(data)
	hash.Write(timeEncoded)
	id = hash.Sum(timeEncoded)
	return
}

// VerifyTimed validates whether an id / data combination is valid and in date.
// Duration determines how long a captcha should be seen as valid.
//
// Can be used concurrent.
func VerifyTimed(id, data []byte, now time.Time, validDuration time.Duration) bool {
	initialisationRandomData.Do(func() {
		setRandomData()
	})

	if len(id) <= hashSize {
		return false
	}

	timeEncoded := make([]byte, len(id)-hashSize)
	copy(timeEncoded, id[:len(id)-hashSize]) // We need a true copy here, or else subtle.ConstantTimeCompare returns always true

	hash := hmac.New(hashGenerator, randomData)
	hash.Write(data)
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
