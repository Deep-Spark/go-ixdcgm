/*
Copyright (c) 2024, Shanghai Iluvatar CoreX Semiconductor Co., Ltd.
All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ixdcgm

/*
#cgo linux LDFLAGS: -ldl -Wl,--export-dynamic -Wl,--unresolved-symbols=ignore-in-object-files

#include <dlfcn.h>
#include <stdlib.h>
#include "include/ixdcgmStructs.h"
#include "include/ixdcgmApiExport.h"
#include "include/dcgm_agent.h"
#include "include/dcgm_structs.h"
*/
import "C"
import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"
)

func uintPtr(c C.uint) *uint {
	i := uint(c)
	return &i
}

func stringPtr(c *C.char) *string {
	s := C.GoString(c)
	return &s
}

type DcgmError struct {
	msg  string         // description of error
	Code C.dcgmReturn_t // dcgmReturn_t value of error
}

func (e *DcgmError) Error() string { return e.msg }

func makeVersion1(struct_type uintptr) C.uint {
	version := C.uint(struct_type | 1<<24)
	return version
}

func makeVersion2(struct_type uintptr) C.uint {
	version := C.uint(struct_type | 2<<24)
	return version
}

func makeVersion3(struct_type uintptr) C.uint {
	version := C.uint(struct_type | 3<<24)
	return version
}

func makeVersion4(struct_type uintptr) C.uint {
	version := C.uint(struct_type | 4<<24)
	return version
}

func makeVersion10(struct_type uintptr) C.uint {
	version := C.uint(struct_type | 10<<24)
	return version
}

func errorString(result C.dcgmReturn_t) error {
	if result == C.DCGM_ST_OK {
		return nil
	}
	err := C.GoString(C.errorString(result))
	return fmt.Errorf("%v", err)
}

func ixdcgmErrorString(result C.ixdcgmReturn_t) error {
	if result == C.IXDCGM_RET_OK {
		return nil
	}
	err := C.GoString(C.ixdcgmErrorString(result))
	return fmt.Errorf("%v", err)
}

func string2Char(c string) *C.char {
	return C.CString(c)
}

func freeCString(c *C.char) {
	C.free(unsafe.Pointer(c))
}

func cChar2String(c *C.char) string {
	return C.GoString(c)
}

func removeBytesSpaces(originalBytes []byte) string {
	lastNonZeroIndex := len(originalBytes) - 1
	for ; lastNonZeroIndex >= 0; lastNonZeroIndex-- {
		if originalBytes[lastNonZeroIndex] != 0 {
			break
		}
	}
	cleanedBytes := originalBytes[:lastNonZeroIndex+1]

	return string(cleanedBytes)
}

// convertBitsetStr converts a set of numbers in string format to a range representation.
// input sample: "{0,1,2,3,6,10,11,12,13}"
// output sample: "0-3,6,10-13"
func convertBitsetStr(input string) (output string) {
	input = strings.Trim(input, "{}")
	numStrs := strings.Split(input, ",")
	nums := make([]int, len(numStrs))

	// Convert string numbers to integers
	for i, numStr := range numStrs {
		num, err := strconv.Atoi(strings.TrimSpace(numStr))
		if err != nil {
			panic(err)
		}
		nums[i] = num
	}

	// Sort the numbers (assuming they are not sorted)
	// If the input is always sorted, you can skip this step
	for i := 0; i < len(nums)-1; i++ {
		for j := 0; j < len(nums)-1-i; j++ {
			if nums[j] > nums[j+1] {
				nums[j], nums[j+1] = nums[j+1], nums[j]
			}
		}
	}

	// Process the numbers to create ranges
	var result []string
	start := nums[0]
	end := nums[0]

	for i := 1; i < len(nums); i++ {
		if nums[i] == end+1 {
			end = nums[i]
		} else {
			if start == end {
				result = append(result, strconv.Itoa(start))
			} else {
				result = append(result, fmt.Sprintf("%d-%d", start, end))
			}
			start = nums[i]
			end = nums[i]
		}
	}

	// Handle the last range
	if start == end {
		result = append(result, strconv.Itoa(start))
	} else {
		result = append(result, fmt.Sprintf("%d-%d", start, end))
	}

	// Join the result into a single string
	output = strings.Join(result, ",")
	return
}

func parseDirPath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("Error to parse dir path %s, err: %v", path, err)
	}
	return absPath, nil
}
