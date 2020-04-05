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

// Package captcha contains methods to generate captchas.
// There are two basic methods: normal captchas and timed captchas.
// A captcha consists of a captcha and an id. The id can be shown publicly, the captcha should be guessed by humans. A captcha can not be derivated from an id other than through brute force.
// The verification of a captcha soley depends on the id, so no state is required. This also means that you can use the functions parallel without problems or side effects.
// There are two caviats, though:
// * A hidden value is used to make predictions impossible. This means that whenever you restart the program, old captchas are no longer valid.
// * No session management is implemented. One captcha / id combination is always valid (as long as the hidden value is the same).
//
// Normal only captchas consist of a id / captcha combination. Therefore, you are strongly advised to use some sort of session management.
// Timed captchas are valid for a specified amount of time. Therefore, a session management might not be needed (but you might use one, too).
package captcha
