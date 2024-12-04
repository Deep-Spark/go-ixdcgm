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
#include "include/dcgm_agent.h"
#include "include/dcgm_structs.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func getDeviceTopology(deviceId uint) (err error) {
	var topology C.dcgmDeviceTopology_v1
	topology.version = makeVersion1(unsafe.Sizeof(topology))

	res := C.dcgmGetDeviceTopology(handle.handle, C.uint(deviceId), &topology)
	if res == C.DCGM_ST_NOT_SUPPORTED {
		fmt.Println("not supported")
		return nil
	}
	if res != C.DCGM_ST_OK {
		return fmt.Errorf("error getting device topology %s", C.GoString(C.errorString(res)))
	}

	return nil
}
