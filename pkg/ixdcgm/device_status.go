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
	"math/rand"
	"time"
)

type PerfState uint

const (
	PerfStateMax     PerfState = 0
	PerfStateMin     PerfState = 15
	PerfStateUnknown PerfState = 32
)

func (p PerfState) String() string {
	if p >= PerfStateMax && p <= PerfStateMin {
		return fmt.Sprintf("P%d", p)
	}
	return "Unknown"
}

type UtilizationInfo struct {
	Gpu int64 // %
	Mem int64 // %
}

type ClockInfo struct {
	Sm  int64 // MHz
	Mem int64 // MHz
}

type PCIStatusInfo struct {
	Rx            int64 // KB/s
	Tx            int64 // KB/s
	ReplayCounter int64 // Counter
}

type MemoryUsage struct {
	Total int64 // Total Memory (Frame Buffer) of the GPU in MB
	Used  int64 // Used Memory (Frame Buffer) in MB
	Free  int64 // Free Memory (Frame Buffer) in MB
}

type DeviceStatus struct {
	Id          uint
	Power       string // "N/A" or float64 str, W
	Temperature string // "N/A" or int64 str, °C

	Utilization UtilizationInfo
	Clocks      ClockInfo
	PCI         PCIStatusInfo
	Performance PerfState
	MemUsage    MemoryUsage

	FanSpeed     string // "N/A" or int64 str, %
	EccSbeVolDev string // "N/A" or int64 str, 1 for errors occurred, 0 for no errors
	EccDbeVolDev string // "N/A" or int64 str, 1 for errors occurred, 0 for no errors
}

type DeviceProfStatus struct {
	SmActive    string // "N/A" or float64 str, %
	SmOccupancy string // "N/A" or float64 str, %
	DramActive  string // "N/A" or float64 str, %
}

func getDeviceStatus(gpuId uint) (status DeviceStatus, err error) {
	const (
		IdxPower int = iota
		IdxGpuTemp
		IdxGpuUtil
		IdxMemUtil
		IdxSmClock
		IdxMemClock
		IdxPcieRxThroughput
		IdxPcieTxThroughput
		IdxPcieReplayCounter
		IdxFanSpeed
		IdxEccSbeVolDev
		IdxEccDbeVolDev
		IdxMemTotal
		IdxMemUsed
		IdxMemFree
	)

	fields := []Short{
		DCGM_FI_DEV_POWER_USAGE,
		DCGM_FI_DEV_GPU_TEMP,
		DCGM_FI_DEV_GPU_UTIL,
		DCGM_FI_DEV_MEM_COPY_UTIL,
		DCGM_FI_DEV_SM_CLOCK,
		DCGM_FI_DEV_MEM_CLOCK,
		DCGM_FI_DEV_PCIE_RX_THROUGHPUT,
		DCGM_FI_DEV_PCIE_TX_THROUGHPUT,
		DCGM_FI_DEV_PCIE_REPLAY_COUNTER,
		DCGM_FI_DEV_FAN_SPEED,
		DCGM_FI_DEV_ECC_SBE_VOL_DEV,
		DCGM_FI_DEV_ECC_DBE_VOL_DEV,
		DCGM_FI_DEV_FB_TOTAL,
		DCGM_FI_DEV_FB_USED,
		DCGM_FI_DEV_FB_FREE,
	}

	fieldGrpName := fmt.Sprintf("devStatusFields%d", rand.Uint64())
	fieldGrp, err := FieldGroupCreate(fieldGrpName, fields)
	if err != nil {
		return
	}

	gpuGrpName := fmt.Sprintf("devStatusGrp%d", rand.Uint64())
	gpuGrpHdl, err := WatchFields([]uint{gpuId}, fieldGrp, gpuGrpName)
	if err != nil {
		_ = FieldGroupDestroy(fieldGrp)
		return
	}

	values, err := GetLatestValuesForFields(gpuId, fields)
	if err != nil {
		_ = FieldGroupDestroy(fieldGrp)
		_ = DestroyGroup(gpuGrpHdl)
		return status, err
	}

	clocks := ClockInfo{
		Sm:  values[IdxSmClock].Int64(),
		Mem: values[IdxMemClock].Int64(),
	}

	utilInfo := UtilizationInfo{
		Gpu: values[IdxGpuUtil].Int64(),
		Mem: values[IdxMemUtil].Int64(),
	}

	pciInfo := PCIStatusInfo{
		Rx:            values[IdxPcieRxThroughput].Int64(),
		Tx:            values[IdxPcieTxThroughput].Int64(),
		ReplayCounter: values[IdxPcieReplayCounter].Int64(),
	}

	memUsage := MemoryUsage{
		Total: values[IdxMemTotal].Int64(),
		Free:  values[IdxMemFree].Int64(),
		Used:  values[IdxMemUsed].Int64(),
	}

	status = DeviceStatus{
		Id:           gpuId,
		Power:        GetFieldValueStr(values[IdxPower], "float64"),
		Temperature:  GetFieldValueStr(values[IdxGpuTemp], "int64"),
		Utilization:  utilInfo,
		Clocks:       clocks,
		PCI:          pciInfo,
		MemUsage:     memUsage,
		FanSpeed:     GetFieldValueStr(values[IdxFanSpeed], "int64"),
		EccSbeVolDev: GetFieldValueStr(values[IdxEccSbeVolDev], "int64"),
		EccDbeVolDev: GetFieldValueStr(values[IdxEccDbeVolDev], "int64"),
	}

	_ = FieldGroupDestroy(fieldGrp)
	_ = DestroyGroup(gpuGrpHdl)
	return
}

func getDeviceProfStatus(gpuId uint) (status DeviceProfStatus, err error) {
	const (
		IdxSmActive int = iota
		IdxSmOccupancy
		IdxDramActive
	)

	fields := []Short{
		DCGM_FI_PROF_SM_ACTIVE,
		DCGM_FI_PROF_SM_OCCUPANCY,
		DCGM_FI_PROF_DRAM_ACTIVE,
	}

	fieldGrpName := fmt.Sprintf("devProfStatusFields%d", rand.Uint64())
	fieldGrp, err := FieldGroupCreate(fieldGrpName, fields)
	if err != nil {
		return
	}

	grpName := fmt.Sprintf("devProfStatusGrp%d", rand.Uint64())
	grpId, err := WatchFields([]uint{gpuId}, fieldGrp, grpName)
	if err != nil {
		_ = FieldGroupDestroy(fieldGrp)
		return
	}

	time.Sleep(2000 * time.Millisecond)
	values, err := GetLatestValuesForFields(gpuId, fields)
	if err != nil {
		_ = FieldGroupDestroy(fieldGrp)
		_ = DestroyGroup(grpId)
		return status, err
	}

	status = DeviceProfStatus{
		SmActive:    GetFieldValueStr(values[IdxSmActive], "float64"),
		SmOccupancy: GetFieldValueStr(values[IdxSmOccupancy], "float64"),
		DramActive:  GetFieldValueStr(values[IdxDramActive], "float64"),
	}

	_ = FieldGroupDestroy(fieldGrp)
	_ = DestroyGroup(grpId)
	return

}
