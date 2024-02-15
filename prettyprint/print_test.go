// Copyright (c) 2023 Cisco Systems, Inc. and its affiliates
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http:www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package prettyprint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrettyPrint(t *testing.T) {
	type Data struct {
		One   string  `json:"one"`
		Two   string  `json:"two"`
		Three *string `json:"three"`
		Four  int     `json:"four"`
		Five  string  `json:"five"` // ignored
		Six   *int64  `json:"six"`
		Seven bool    `json:"seven"`
	}
	a := "a"
	ab := "ab"
	abc := "abc"
	num := int64(123456789012345)
	data := []Data{
		{
			One:   "a",
			Two:   "ab",
			Three: &abc,
			Four:  1,
			Six:   &num,
			Seven: true,
		},
		{
			One:   "abcdefg",
			Two:   "abc",
			Three: &a,
			Four:  12,
			Six:   &num,
			Seven: false,
		},
		{
			One:   "abc",
			Two:   "abcde",
			Three: &ab,
			Four:  123,
			Six:   &num,
			Seven: true,
		},
	}
	names := []Display{
		{"One", "A"},
		{"Two", "B"},
		{"Three", "C"},
		{"Six", "E"},
		{"Four", "D"},
		{"Seven", "F"},
	}
	expected := `A         B       C     E                 D     F
a         ab      abc   123456789012345   1     true
abcdefg   abc     a     123456789012345   12    false
abc       abcde   ab    123456789012345   123   true
`

	result := getPrettyFormat(data, names)
	require.Equal(t, expected, result)

	names = []Display{
		{"One", "A"},
		{"Two", "LONG_NAME"},
		{"Three", "C"},
	}
	data = append(data, Data{"abcdefghij", "a", nil, 3, "ignored", &num, true})
	expected = `A            LONG_NAME   C
a            ab          abc
abcdefg      abc         a
abc          abcde       ab
abcdefghij   a           -
`

	result = getPrettyFormat(data, names)
	require.Equal(t, expected, result)
}
