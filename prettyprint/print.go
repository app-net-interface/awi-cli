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
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	minimumSpace = 3
)

type Display struct {
	Name    string
	Display string
}

const (
	jsonFormat string = "json"
)

func PrintData[T any](data []T, displays []Display, format string) {
	if format == jsonFormat {
		fmt.Println(getJsonFormat(data))
	} else {
		fmt.Println(getPrettyFormat(data, displays))
	}
}

func PrintConvertedData[T, C any](data []T, convertedData []C, displays []Display, format string) {
	if format == jsonFormat {
		fmt.Println(getJsonFormat(data))
	} else {
		fmt.Println(getPrettyFormat(convertedData, displays))
	}
}

func getJsonFormat(data any) string {
	result, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Printf("could not parse data: %v", err)
	}
	return string(result)
}

func getPrettyFormat[T any](data []T, displays []Display) string {
	format := getFormat(data, displays)
	formatArguments := make([]any, 0, len(displays))
	for _, display := range displays {
		formatArguments = append(formatArguments, display.Display)
	}
	var buffer strings.Builder
	buffer.WriteString(fmt.Sprintf(format, formatArguments...))
	for _, d := range data {
		formatArguments = make([]any, 0, len(displays))
		for _, display := range displays {
			res := stringValueOf(d, display.Name)
			formatArguments = append(formatArguments, res)
		}
		buffer.WriteString(fmt.Sprintf(format, formatArguments...))
	}
	return buffer.String()
}

func getFormat[T any](data []T, displays []Display) string {
	longestFields := make([]int, len(displays))
	for i, display := range displays {
		longestFields[i] = minimumSpace + len(display.Display)
		for _, d := range data {
			res := stringValueOf(d, display.Name)
			longestFields[i] = max(longestFields[i], len(res)+minimumSpace)
		}
	}
	var formatBuffer strings.Builder
	for i := 0; i < len(longestFields)-1; i++ {
		formatBuffer.WriteString(fmt.Sprintf("%%-%ds", longestFields[i]))
	}
	formatBuffer.WriteString("%s\n")
	return formatBuffer.String()
}

func stringValueOf(d any, fieldName string) string {
	v := reflect.ValueOf(d)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	val := v.FieldByName(fieldName)

	switch val.Kind() {
	case reflect.Ptr:
		if val.IsZero() {
			return "-"
		}
		val = val.Elem()
	case reflect.Int:
		return strconv.FormatInt(val.Int(), 10)
	case reflect.Int64:
		return strconv.FormatInt(val.Int(), 10)
	case reflect.Bool:
		return strconv.FormatBool(val.Bool())
	case reflect.Map:
		str := ""
		if val.Len() > 0 {
			keys := val.MapKeys()
			for _, key := range keys {
				value := val.MapIndex(key)
				str += fmt.Sprintf(key.String() + "=" + value.String() + ",")
			}
			return str[:len(str)-1]
		}
		return str
	}
	return val.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
