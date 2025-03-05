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
	"os"
	"unsafe"
)

const (
	defaultUpdateFreq     = 1000000 // usec
	defaultMaxKeepAge     = 0       // sec
	defaultMaxKeepSamples = 1       // Keep one sample by default since we only ask for latest
)

type FieldGrpHandle struct{ handle C.dcgmFieldGrp_t }

func FieldGroupCreate(groupName string, fields []Short) (fgId FieldGrpHandle, err error) {
	var fieldsGroup C.dcgmFieldGrp_t
	cfields := *(*[]C.ushort)(unsafe.Pointer(&fields))

	gn := string2Char(groupName)
	defer freeCString(gn)

	res := C.dcgmFieldGroupCreate(handle.handle, C.int(len(fields)), &cfields[0], gn, &fieldsGroup)
	if err = errorString(res); err != nil {
		return fgId, fmt.Errorf("error creating DCGM fields group: %s", err)
	}

	fgId = FieldGrpHandle{
		handle: fieldsGroup,
	}
	return
}

func FieldGroupDestroy(fieldGroup FieldGrpHandle) (err error) {
	res := C.dcgmFieldGroupDestroy(handle.handle, fieldGroup.handle)
	if err = errorString(res); err != nil {
		return fmt.Errorf("error destroying DCGM fields group: %s", err)
	}
	return nil
}

func WatchFields(gpuIds []uint, fieldGrp FieldGrpHandle, groupName string) (GroupHandle, error) {
	group, err := CreateGroup(groupName)
	if err != nil {
		return GroupHandle{}, err
	}
	for _, gpuId := range gpuIds {
		err = AddToGroup(group, gpuId)
		if err != nil {
			return GroupHandle{}, err
		}
	}

	res := C.dcgmWatchFields(handle.handle, group.handle, fieldGrp.handle,
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
	return group, nil
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

func GetFieldValueStr(fv FieldValue_v1, typ string) string {
	st := fv.Status
	if st != C.DCGM_ST_OK {
		return "N/A"
	}

	switch typ {
	case "int64":
		value := *(*int64)(unsafe.Pointer(&fv.Value[0]))
		if value >= DCGM_FT_INT64_BLANK {
			return "N/A" // indicate the field is not supported
		}
		return fmt.Sprintf("%d", value)

	case "float64":
		value := *(*float64)(unsafe.Pointer(&fv.Value[0]))
		if value >= DCGM_FT_FP64_BLANK {
			return "N/A" // indicate the field is not supported
		}
		// sync the precision with the display of ixdcgmi
		return fmt.Sprintf("%.3f", value)

	case "string":
		// remove redundant spaces of string converted from C bytes
		return removeBytesSpaces(fv.Value[:])

	default:
		fmt.Printf("Not Supported Type: %s\n", typ)
		os.Exit(1)
		return "N/A"
	}
}

type Field_Entity_Group uint

const (
	FE_NONE Field_Entity_Group = iota
	FE_GPU
	FE_VGPU
	FE_SWITCH
	FE_GPU_I
	FE_GPU_CI
	FE_LINK
	FE_CPU
	FE_CPU_CORE
	FE_COUNT
)

type GroupEntityPair struct {
	EntityGroupId Field_Entity_Group
	EntityId      uint
}

func (e Field_Entity_Group) String() string {
	switch e {
	case FE_GPU:
		return "GPU"
	case FE_VGPU:
		return "vGPU"
	case FE_SWITCH:
		return "Switch"
	case FE_GPU_I:
		return "GPU Instance"
	case FE_GPU_CI:
		return "GPU Compute Instance"
	case FE_LINK:
		return "Link"
	case FE_CPU:
		return "CPU"
	case FE_CPU_CORE:
		return "CPU Core"
	}
	return "unknown"
}

func toFieldValue(cfields []C.dcgmFieldValue_v1) []FieldValue_v1 {
	fields := make([]FieldValue_v1, len(cfields))
	for i, f := range cfields {
		fields[i] = FieldValue_v1{
			Version:   uint(f.version),
			FieldId:   uint(f.fieldId),
			FieldType: uint(f.fieldType),
			Status:    int(f.status),
			Ts:        int64(f.ts),
			Value:     f.value,
		}
	}

	return fields
}

func (fv FieldValue_v1) Int64() int64 {
	return *(*int64)(unsafe.Pointer(&fv.Value[0]))
}

func (fv FieldValue_v1) Float64() float64 {
	return *(*float64)(unsafe.Pointer(&fv.Value[0]))
}

func (fv FieldValue_v1) String() string {
	return C.GoString((*C.char)(unsafe.Pointer(&fv.Value[0])))
}

func (fv FieldValue_v1) Blob() [4096]byte {
	return fv.Value
}
