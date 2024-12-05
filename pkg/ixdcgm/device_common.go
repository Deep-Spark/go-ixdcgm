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
#cgo LDFLAGS: -ldl
#include "include/ixdcgmFields.h"
#include "include/ixdcgmStructs.h"
#include "include/ixdcgmApiExport.h"
*/
import "C"
import "fmt"

func getDeviceOnSameBoard(gpuId1, gpuId2 uint) (isOnSameBoard bool, err error) {
	var onSameBoard C.int
	r := C.ixdcgmDeviceOnSameBoard(C.ulong(handle.handle), C.uint(gpuId1), C.uint(gpuId2), &onSameBoard)
	fmt.Println(r)
	if err = ixdcgmErrorString(r); err != nil {
		return false, err
	}
	if onSameBoard == 0 {
		isOnSameBoard = false
	} else {
		isOnSameBoard = true
	}

	return isOnSameBoard, nil
}
