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
)

type PerfState uint

const (
	PerfStateMax     = 0
	PerfStateMin     = 15
	PerfStateUnknown = 32
)

func (p PerfState) String() string {
	if p >= PerfStateMax && p <= PerfStateMin {
		return fmt.Sprintf("P%d", p)
	}
	return "Unknown"
}

type UtilizationInfo struct {
	GPU    int64 // %
	Memory int64 // %
}

type ClockInfo struct {
	Cores  int64 // MHz
	Memory int64 // MHz
}

type PCIStatusInfo struct {
	Rx      int64 // MB/s
	Tx      int64 // MB/s
	Replays int64
}

type DeviceStatus struct {
	Power       float64 // W
	Temperature int64   // °C
	Utilization UtilizationInfo
	Clocks      ClockInfo
	PCI         PCIStatusInfo
	Performance PerfState
	FanSpeed    int64 // %
}

func getDeviceStatus(gpuId uint) (status DeviceStatus, err error) {
	const (
		pwr int = iota
		temp
		sm
		mem
		smClock
		memClock
		pcieRxThroughput
		pcieTxThroughput
		pcieReplay
		fanSpeed
	)

	deviceFields := []Short{
		C.DCGM_FI_DEV_POWER_USAGE,
		C.DCGM_FI_DEV_GPU_TEMP,
		C.DCGM_FI_DEV_GPU_UTIL,
		C.DCGM_FI_DEV_MEM_COPY_UTIL,
		C.DCGM_FI_DEV_SM_CLOCK,
		C.DCGM_FI_DEV_MEM_CLOCK,
		C.DCGM_FI_DEV_PCIE_RX_THROUGHPUT,
		C.DCGM_FI_DEV_PCIE_TX_THROUGHPUT,
		C.DCGM_FI_DEV_PCIE_REPLAY_COUNTER,
		C.DCGM_FI_DEV_FAN_SPEED,
	}

	fieldsName := fmt.Sprintf("devStatusFields%d", rand.Uint64())
	fieldsId, err := FieldGroupCreate(fieldsName, deviceFields)
	if err != nil {
		return
	}

	groupName := fmt.Sprintf("devStatus%d", rand.Uint64())
	groupId, err := WatchFields(gpuId, fieldsId, groupName)
	if err != nil {
		_ = FieldGroupDestroy(fieldsId)
		return
	}

	values, err := GetLatestValuesForFields(gpuId, deviceFields)
	if err != nil {
		_ = FieldGroupDestroy(fieldsId)
		_ = DestroyGroup(groupId)
		return status, err
	}

	power := values[pwr].Float64()

	clocks := ClockInfo{
		Cores:  values[smClock].Int64(),
		Memory: values[memClock].Int64(),
	}

	gpuUtil := UtilizationInfo{
		GPU:    values[sm].Int64(),
		Memory: values[mem].Int64(),
	}

	pci := PCIStatusInfo{
		Rx:      values[pcieRxThroughput].Int64(),
		Tx:      values[pcieTxThroughput].Int64(),
		Replays: values[pcieReplay].Int64(),
	}
	status = DeviceStatus{
		Power:       power,
		Temperature: values[temp].Int64(),
		Utilization: gpuUtil,
		Clocks:      clocks,
		PCI:         pci,
		FanSpeed:    values[fanSpeed].Int64(),
	}

	_ = FieldGroupDestroy(fieldsId)
	_ = DestroyGroup(groupId)
	return
}
