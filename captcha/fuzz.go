// +build gofuzz

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
	"encoding/gob"
	"time"
)

var FuzzTime = time.Date(2020, 01, 01, 20, 0, 0, 0, time.UTC)

type FuzzData struct {
	Id, C []byte
	Size  int
}

func Fuzz(data []byte) int {
	buf := bytes.NewBuffer(data)
	var fd FuzzData
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&fd)
	if err != nil {
		return -1
	}

	if VerifyTimed(fd.Id, fd.C, FuzzTime, 24*time.Hour, fd.Size) {
		return 1
	}
	return 0
}
