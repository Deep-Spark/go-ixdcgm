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

	"github.com/bits-and-blooms/bitset"
)

type DeviceIdentifier struct {
	ProductName   string
	DeviceName    string
	Serial        string
	DriverVersion string
}

type PciInfo struct {
	BusId string

	Bandwidth int64 // MB/s
}

type MemoryUsageInfo struct {
	Version uint
	BAR1    uint // MB
	Total   uint // MB
	Used    uint // MB
	Free    uint // MB
}

type DeviceInfo struct {
	GPUId           uint
	IxDCGMSupported string
	Uuid            string
	PowerLimit      uint // W
	PCI             PciInfo
	MemoryUsage     MemoryUsageInfo
	Identifiers     DeviceIdentifier
	Topology        []P2PLink
	CPUAffinity     string
	NUMAAffinity    string
}

func getAllDeviceCount() (gpuCount uint, err error) {
	var gpuIdList [C.DCGM_MAX_NUM_DEVICES]C.uint
	var count C.int

	r := C.dcgmGetAllDevices(C.ulong(handle.handle), &gpuIdList[0], &count)
	if err = errorString(r); err != nil {
		return gpuCount, err
	}

	gpuCount = uint(count)
	return
}

func getPciBandwidth(gpuId uint) (int64, error) {
	const (
		maxLinkGen int = iota
		maxLinkWidth
		fieldsCount
	)

	pciFields := []Short{
		C.DCGM_FI_DEV_PCIE_MAX_LINK_GEN,
		C.DCGM_FI_DEV_PCIE_MAX_LINK_WIDTH,
	}

	fieldName := fmt.Sprintf("pciBandwidthFields%d", gpuId)
	fieldsId, err := FieldGroupCreate(fieldName, pciFields)
	if err != nil {
		return 0, err
	}

	groupName := fmt.Sprintf("pciBandwidth%d", gpuId)
	groupId, err := WatchFields([]uint{gpuId}, fieldsId, groupName)
	if err != nil {
		FieldGroupDestroy(fieldsId)
		return 0, err
	}

	values, err := GetLatestValuesForFields(gpuId, pciFields)
	if err != nil {
		FieldGroupDestroy(fieldsId)
		DestroyGroup(groupId)
		return 0, fmt.Errorf("failed to get pcie bandwidgth: %s", err)
	}

	var genMap = map[int64]int64{
		1: 250,
		2: 500,
		3: 985,
		4: 1969,
	}

	FieldGroupDestroy(fieldsId)
	DestroyGroup(groupId)

	gen := values[maxLinkGen].Int64()
	width := values[maxLinkWidth].Int64()

	bandwidth := genMap[gen] * width
	return bandwidth, nil
}

func getDeviceInfo(gpuId uint) (DeviceInfo, error) {
	var dcgmAttr C.dcgmDeviceAttributes_t
	dcgmAttr.version = C.uint(makeVersion3(unsafe.Sizeof(dcgmAttr)))

	res := C.dcgmGetDeviceAttributes(C.ulong(handle.handle), C.uint(gpuId), &dcgmAttr)
	if err := errorString(res); err != nil {
		return DeviceInfo{}, err
	}

	// check if the given GPU is IxDCGM supported
	gpus, err := getSupportedDevices()
	if err != nil {
		return DeviceInfo{}, err
	}

	supported := "N"
	for _, gpu := range gpus {
		if gpuId == gpu {
			supported = "Y"
			break
		}
	}

	cpuAffinity, err := getAffinity(gpuId, "CPU")
	if err != nil {
		fmt.Printf("Error getting cpu affinity, set CPU Affinity to N/A, err: %s", err)
	}
	numaAffinity, err := getAffinity(gpuId, "NUMA")
	if err != nil {
		fmt.Printf("Error getting numa affinity, set NUMA Affinity to N/A, err: %s", err)
	}

	var topology []P2PLink
	var bandwidth int64
	if supported == "Y" {
		topology, err = getDeviceTopology(gpuId)
		if err != nil {
			return DeviceInfo{}, err
		}
		bandwidth, err = getPciBandwidth(gpuId)
		if err != nil {
			return DeviceInfo{}, err
		}
	}

	uuid := cChar2String(&dcgmAttr.identifiers.uuid[0])
	powerLimit := uint(dcgmAttr.powerLimits.defaultPowerLimit)
	busId := cChar2String(&dcgmAttr.identifiers.pciBusId[0])

	pci := PciInfo{
		BusId:     busId,
		Bandwidth: bandwidth,
	}

	id := DeviceIdentifier{
		ProductName:   cChar2String(&dcgmAttr.identifiers.deviceName[0]),
		Serial:        cChar2String(&dcgmAttr.identifiers.serial[0]),
		DriverVersion: cChar2String(&dcgmAttr.identifiers.driverVersion[0]),
	}

	memInfo := MemoryUsageInfo{
		Total: uint(dcgmAttr.memoryUsage.fbTotal),
		Used:  uint(dcgmAttr.memoryUsage.fbUsed),
		Free:  uint(dcgmAttr.memoryUsage.fbFree),
	}

	return DeviceInfo{
		GPUId:           gpuId,
		IxDCGMSupported: supported,
		Uuid:            uuid,
		PowerLimit:      powerLimit,
		PCI:             pci,
		MemoryUsage:     memInfo,
		Identifiers:     id,
		Topology:        topology,
		CPUAffinity:     cpuAffinity,
		NUMAAffinity:    numaAffinity,
	}, nil
}

func getSupportedDevices() (gpus []uint, err error) {
	var gpuIdList [C.DCGM_MAX_NUM_DEVICES]C.uint
	var count C.int

	r := C.dcgmGetAllSupportedDevices(C.ulong(handle.handle), &gpuIdList[0], &count)
	if err = errorString(r); err != nil {
		return gpus, err
	}
	numGpus := uint(count)
	gpus = make([]uint, numGpus)
	for i := uint(0); i < numGpus; i++ {
		gpus[i] = uint(gpuIdList[i])
	}
	return
}

// if err is not nil, return "N/A" as result
func getAffinity(gpuId uint, typ string) (result string, err error) {
	const (
		affinity0 int = iota
		affinity1
		affinity2
		affinity3
		fieldsCount
	)

	affFields := make([]Short, fieldsCount)
	switch typ {
	case "CPU":
		affFields[affinity0] = DCGM_FI_DEV_CPU_AFFINITY_0
		affFields[affinity1] = DCGM_FI_DEV_CPU_AFFINITY_1
		affFields[affinity2] = DCGM_FI_DEV_CPU_AFFINITY_2
		affFields[affinity3] = DCGM_FI_DEV_CPU_AFFINITY_3
	case "NUMA":
		affFields[affinity0] = DCGM_FI_DEV_MEM_AFFINITY_0
		affFields[affinity1] = DCGM_FI_DEV_MEM_AFFINITY_1
		affFields[affinity2] = DCGM_FI_DEV_MEM_AFFINITY_2
		affFields[affinity3] = DCGM_FI_DEV_MEM_AFFINITY_3
	default:
		return "N/A", fmt.Errorf("not supported affinity type: %s", typ)
	}

	fieldGrpName := fmt.Sprintf("%sAffFields%d", typ, gpuId)
	fieldGrpHdl, err := FieldGroupCreate(fieldGrpName, affFields)
	if err != nil {
		return "N/A", err
	}
	defer FieldGroupDestroy(fieldGrpHdl)

	gpuGrpName := fmt.Sprintf("%sAff%d", typ, gpuId)
	gpuGrpHdl, err := WatchFields([]uint{gpuId}, fieldGrpHdl, gpuGrpName)
	if err != nil {
		return "N/A", err
	}
	defer DestroyGroup(gpuGrpHdl)

	values, err := GetLatestValuesForFields(gpuId, affFields)
	if err != nil {
		return "N/A", fmt.Errorf("Error getting %s affinity: %s", typ, err)
	}

	bits := make([]int64, 4)
	bits[0] = values[affinity0].Int64()
	bits[1] = values[affinity1].Int64()
	bits[2] = values[affinity2].Int64()
	bits[3] = values[affinity3].Int64()
	for _, bit := range bits {
		if bit >= DCGM_FT_INT64_BLANK {
			// Retrieved affinity value is invalid.
			return "N/A", nil
		}
	}

	ubits := make([]uint64, len(bits))
	for i, val := range bits {
		ubits[i] = uint64(val)
	}
	b := bitset.From(ubits)
	str := convertBitsetStr(b.String())
	return str, nil
}
