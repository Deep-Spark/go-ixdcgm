/*
Copyright (c) 2024, NVIDIA CORPORATION.
Copyright (c) 2024, Shanghai Iluvatar CoreX Semiconductor Co., Ltd.
All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

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
import "unsafe"

type IxDcgmStatus struct {
	Memory int64 // KB
	CPU    float64
}

func introspect() (engine IxDcgmStatus, err error) {
	var memory C.dcgmIntrospectMemory_t
	memory.version = makeVersion1(unsafe.Sizeof(memory))
	waitIfNoData := 1
	result := C.dcgmIntrospectGetHostengineMemoryUsage(handle.handle, &memory, C.int(waitIfNoData))

	if err = errorString(result); err != nil {
		return engine, &DcgmError{msg: C.GoString(C.errorString(result)), Code: result}
	}

	var cpu C.dcgmIntrospectCpuUtil_t

	cpu.version = makeVersion1(unsafe.Sizeof(cpu))
	result = C.dcgmIntrospectGetHostengineCpuUtilization(handle.handle, &cpu, C.int(waitIfNoData))

	if err = errorString(result); err != nil {
		return engine, &DcgmError{msg: C.GoString(C.errorString(result)), Code: result}
	}

	engine = IxDcgmStatus{
		Memory: toInt64(memory.bytesUsed) / 1024,
		CPU:    *dblToFloat(cpu.total) * 100,
	}
	return
}
