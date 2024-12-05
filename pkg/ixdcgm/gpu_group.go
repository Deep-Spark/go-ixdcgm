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

type GroupHandle struct {
	handle C.dcgmGpuGrp_t
}

func CreateGroup(groupName string) (GroupHandle, error) {
	var cGroupId C.dcgmGpuGrp_t
	cgn := string2Char(groupName)
	defer freeCString(cgn)

	res := C.dcgmGroupCreate(handle.handle, C.DCGM_GROUP_EMPTY, cgn, &cGroupId)
	if err := errorString(res); err != nil {
		return GroupHandle{}, err
	}

	return GroupHandle{cGroupId}, nil
}

func AddDevice(groupId GroupHandle, gpuId uint) error {
	res := C.dcgmGroupAddDevice(handle.handle, groupId.handle, C.uint(gpuId))
	if err := errorString(res); err != nil {
		return err
	}
	return nil
}

func DestroyGroup(groupId GroupHandle) error {
	res := C.dcgmGroupDestroy(handle.handle, groupId.handle)
	if err := errorString(res); err != nil {
		return err
	}
	return nil
}
