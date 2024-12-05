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

const (
	defaultUpdateFreq     = 30000000 // usec
	defaultMaxKeepAge     = 0        // sec
	defaultMaxKeepSamples = 1        // Keep one sample by default since we only ask for latest
)

type FieldHandle struct{ handle C.dcgmFieldGrp_t }

func FieldGroupCreate(groupName string, fields []Short) (fieldsId FieldHandle, err error) {
	var fieldsGroup C.dcgmFieldGrp_t
	cfields := *(*[]C.ushort)(unsafe.Pointer(&fields))

	gn := string2Char(groupName)
	defer freeCString(gn)

	res := C.dcgmFieldGroupCreate(handle.handle, C.int(len(fields)), &cfields[0], gn, &fieldsGroup)
	if err = errorString(res); err != nil {
		return fieldsId, fmt.Errorf("error creating DCGM fields group: %s", err)
	}

	fieldsId = FieldHandle{
		handle: fieldsGroup,
	}
	return
}

func FieldGroupDestroy(fieldGroup FieldHandle) (err error) {
	res := C.dcgmFieldGroupDestroy(handle.handle, fieldGroup.handle)
	if err = errorString(res); err != nil {
		return fmt.Errorf("error destroying DCGM fields group: %s", err)
	}
	return nil
}

func WatchFields(gpuId uint, fieldsGroup FieldHandle, groupName string) (GroupHandle, error) {
	groups, err := CreateGroup(groupName)
	if err != nil {
		return GroupHandle{}, err
	}

	err = AddDevice(groups, gpuId)
	if err != nil {
		return GroupHandle{}, err
	}

	res := C.dcgmWatchFields(handle.handle, groups.handle, fieldsGroup.handle,
		C.longlong(defaultUpdateFreq),
		C.double(defaultMaxKeepAge),
		C.int(defaultMaxKeepSamples))
	if err = errorString(res); err != nil {
		return GroupHandle{}, fmt.Errorf("error watching DCGM fields: %s", err)
	}

	cWaitForUpdate := C.int(1)
	res = C.dcgmUpdateAllFields(handle.handle, cWaitForUpdate)
	if err = errorString(res); err != nil {
		return GroupHandle{}, fmt.Errorf("error updating DCGM fields: %s", err)
	}
	return groups, nil
}

func GetLatestValuesForFields(gpu uint, fields []Short) ([]FieldValue_v1, error) {
	values := make([]C.dcgmFieldValue_v1, len(fields))
	cFields := *(*[]C.ushort)(unsafe.Pointer(&fields))
	res := C.dcgmGetLatestValuesForFields(handle.handle, C.int(gpu), &cFields[0], C.uint(len(fields)), &values[0])
	if err := errorString(res); err != nil {
		return nil, fmt.Errorf("error getting latest DCGM fields values: %s", err)
	}
	return toFieldValue(values), nil
}

func toFieldValue(values []C.dcgmFieldValue_v1) (fields []FieldValue_v1) {
	fields = make([]FieldValue_v1, len(values))
	for i, v := range values {
		fields[i] = FieldValue_v1{
			Version:   uint(v.version),
			FieldId:   uint(v.fieldId),
			FieldType: uint(v.fieldType),
			Status:    int(v.status),
			Ts:        int64(v.ts),
			Value:     v.value,
		}
	}
	return
}

// func FieldsInit() int {
// 	return int(C.ixdcgmFieldsInit())
// }

// func FieldsTerm() int {
// 	return int(C.ixdcgmFieldsTerm())
// }
