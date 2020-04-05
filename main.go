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

// This contains a small test wrapper
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/toqueteos/webbrowser"
	"github.com/Top-Ranger/auth/captcha"
	"github.com/Top-Ranger/auth/data"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		i, c, err := captcha.GetStrings()
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		rw.Write([]byte("Normal:"))
		rw.Write([]byte("\n   id: "))
		rw.Write([]byte(i))
		rw.Write([]byte("\n   captcha: "))
		rw.Write([]byte(c))

		i, c, err = captcha.GetStringsTimed(time.Now())
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write([]byte("\nTimed:"))
		rw.Write([]byte("\n   id: "))
		rw.Write([]byte(i))
		rw.Write([]byte("\n   captcha: "))
		rw.Write([]byte(c))

		i, err = data.GetStrings("data")
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write([]byte("\ndata:"))
		rw.Write([]byte("\n   id: "))
		rw.Write([]byte(i))
		rw.Write([]byte("\n   data: "))
		rw.Write([]byte("data"))

		i, err = data.GetStringsTimed(time.Now(), "data")
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write([]byte("\nTimed data:"))
		rw.Write([]byte("\n   id: "))
		rw.Write([]byte(i))
		rw.Write([]byte("\n   data: "))
		rw.Write([]byte("data"))
	})

	go func() {
		time.Sleep(1 * time.Second)
		webbrowser.Open("http://localhost:8080")
	}()
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
