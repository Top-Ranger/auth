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
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Top-Ranger/auth/data"
)

func main() {
	err := os.MkdirAll("../data/corpus/", os.ModePerm)
	if err != nil {
		panic(err)
	}
	t := data.FuzzTime.Add(0)
	for i := 0; i < 100; i++ {
		t = t.Add(-1)
		f, err := os.Create(filepath.Join("../data/corpus/", strconv.Itoa(i)))
		if err != nil {
			panic(err)
		}

		b, err := data.GetTimed(data.FuzzTime, []byte("test"))
		if err != nil {
			panic(err)
		}

		if !data.VerifyTimed(b, []byte("test"), data.FuzzTime, 24*time.Hour) {
			panic("Invalid ID")
		}

		_, err = f.Write(b)
		if err != nil {
			panic(err)
		}

		f.Close()
	}
}
