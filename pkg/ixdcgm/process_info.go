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
import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

type DeviceProcessInfo struct {
	Pid           uint64
	Name          string
	UsedGpuMemory uint64 // MiB
}

func getDeviceRunningProcesses(gpuId uint) ([]DeviceProcessInfo, error) {
	cnt, pids, usedMemoryBytes, err := ixdcgmGetDeviceRunningProcesses(gpuId)
	if err != nil {
		return nil, err
	}
	InfoCount := int(uint32(cnt))
	infos := make([]DeviceProcessInfo, InfoCount)
	for i := 0; i < InfoCount; i++ {
		infos[i].Pid = uint64(pids[i])
		infos[i].Name = getPidName(uint64(pids[i]))
		infos[i].UsedGpuMemory = uint64(usedMemoryBytes[i]) / 1024 / 1024
	}
	return infos, nil
}

func ixdcgmGetDeviceRunningProcesses(gpuId uint) (cnt C.uint32_t, pids []C.uint64_t, usedMemoryBytes []C.uint64_t, err error) {
	cnt = 1
	for i := 0; i < 2; i++ {
		pids = make([]C.uint64_t, cnt)
		usedMemoryBytes = make([]C.uint64_t, cnt)
		ret := C.ixdcgmGetDeviceRunningProcesses(C.ulong(handle.handle), C.uint(gpuId), &cnt, &pids[0], &usedMemoryBytes[0])
		if ret == C.IXDCGM_RET_OK {
			// fmt.Printf("the number of valid pids/usedMemoryBytes info is %d\n", uint32(cnt))
			err = nil
			return
		} else if ret == C.IXDCGM_RET_INSUFFICIENT_SIZE {
			// fmt.Printf("INSUFFICIENT_SIZE Warnnig: the needed buffer size is %d\n", uint32(cnt))
			continue
		} else if ret == C.IXDCGM_RET_BADPARAM {
			err = fmt.Errorf("bad parameter")
			return
		}
	}
	err = fmt.Errorf("failed to call ixdcgm api with the needed buffer size %d", uint32(cnt))
	return
}

func getPidName(pid uint64) string {
	cmdlinePath := fmt.Sprintf("/proc/%d/cmdline", pid)
	data, err := os.ReadFile(cmdlinePath)
	if err != nil {
		return ""
	}
	data = bytes.ReplaceAll(data, []byte{0}, []byte{' '})
	return strings.TrimSuffix(string(data), "\x00")
}
