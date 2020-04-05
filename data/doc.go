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

// Package data contains methods to authenticate data.
// There are two basic methods: normal authentification and timed authentification.
// An authentification consists of an id. The id can be shown publicly and can not be derivated from the data other than through brute force.
// The verification of the data soley depends on the id, so no state is required. This also means that you can use the functions parallel without problems or side effects.
// There are two caviats, though:
// * A hidden value is used to make predictions impossible. This means that whenever you restart the program, old ids are no longer valid.
// * One data / id combination is always valid (as long as the hidden value is the same).
package data
