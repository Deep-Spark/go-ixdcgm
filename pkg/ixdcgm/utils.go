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
	"unsafe"
)

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
