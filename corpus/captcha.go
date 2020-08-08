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

package main

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Top-Ranger/auth/captcha"
)

func main() {
	err := os.MkdirAll("../captcha/corpus/", os.ModePerm)
	if err != nil {
		panic(err)
	}
	t := captcha.FuzzTime.Add(0)
	for i := 1; i <= 100; i++ {
		t = t.Add(-1)
		f, err := os.Create(filepath.Join("../captcha/corpus/", strconv.Itoa(i)))
		if err != nil {
			panic(err)
		}

		id, c, err := captcha.GetTimed(t, i)
		if err != nil {
			panic(err)
		}

		if !captcha.VerifyTimed(id, c, captcha.FuzzTime, 24*time.Hour, i) {
			panic("Invalid ID")
		}

		data := captcha.FuzzData{id, c, i}

		enc := gob.NewEncoder(f)
		err = enc.Encode(&data)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}
