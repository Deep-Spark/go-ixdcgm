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

#include <dlfcn.h>
#include <stdlib.h>
#include "include/dcgm_agent.h"
#include "include/dcgm_structs.h"
*/
import "C"
import (
	"context"
	"fmt"
	"sync"
	"unsafe"

	_ "gitee.com/deep-spark/go-ixdcgm/pkg/ixdcgm/include"
)

var (
	ixdcgmLibHandler  unsafe.Pointer
	ixdcgmInitCounter int
	mux               sync.Mutex
	connection        Interface
	handle            DcgmHandle
)

// dynamic library path
const (
	ixdcgmLib = "libixdcgm.so"
)

func initIxDcgm(m int) (err error) {
	lib := string2Char(ixdcgmLib)
	defer freeCString(lib)

	ixdcgmLibHandler = C.dlopen(lib, C.RTLD_LAZY|C.RTLD_GLOBAL)
	if ixdcgmLibHandler == nil {
		return fmt.Errorf("failed to load %s", ixdcgmLib)
	}

	connection, err = New(m)
	if err != nil {
		return err
	}
	return nil
}

func Init(m int, args ...string) (cleanup func(), err error) {
	mux.Lock()
	defer mux.Unlock()
	if ixdcgmInitCounter < 0 {
		return nil, fmt.Errorf("ixdcgm already initialized %d", ixdcgmInitCounter)
	}
	if ixdcgmInitCounter == 0 {
		err = initIxDcgm(m)
		if err != nil {
			return nil, err
		}

		handle, err = connection.Start(args...)
		if err != nil {
			return nil, err
		}
		cleanup = func() {
			shutdown()
		}
	}

	ixdcgmInitCounter += 1

	return cleanup, err
}

func shutdown() (err error) {
	mux.Lock()
	defer mux.Unlock()

	if ixdcgmInitCounter <= 0 {
		return fmt.Errorf("ixdcgm already shutdown")
	}

	if ixdcgmInitCounter == 1 {
		err = connection.Shutdown()
		if err != nil {
			return err
		}
	}

	C.dlclose(ixdcgmLibHandler)
	ixdcgmInitCounter -= 1
	return nil
}

func GetAllDeviceCount() (uint, error) {
	return getAllDeviceCount()
}

func GetSupportedDevices() ([]uint, error) {
	return getSupportedDevices()
}

// GetDeviceInfo describes the given device
func GetDeviceInfo(gpuId uint) (DeviceInfo, error) {
	return getDeviceInfo(gpuId)
}

// GetDeviceStatus monitors GPU status including its power, memory and GPU utilization
func GetDeviceStatus(gpuId uint) (DeviceStatus, error) {
	return getDeviceStatus(gpuId)
}

// GetDeviceProfStatus monitors GPM info including SM_ACTIVE, SM_OCCUPANCY and DRAM_ACTIVE
func GetDeviceProfStatus(gpuId uint) (DeviceProfStatus, error) {
	return getDeviceProfStatus(gpuId)
}

// GetDeviceRunningProcess get the running process infos for the given gpu id
func GetDeviceRunningProcesses(gpuId uint) ([]DeviceProcessInfo, error) {
	return getDeviceRunningProcesses(gpuId)
}

// GetDeviceRunning checks whether the two GPUs are on the same board
func GetDeviceOnSameBoard(gpuId1, gpuId2 uint) (bool, error) {
	return getDeviceOnSameBoard(gpuId1, gpuId2)
}

// HealthCheckByGpuId monitors GPU health for any errors/failures/warnings
func HealthCheckByGpuId(gpuId uint) (DeviceHealth, error) {
	return healthCheckByGpuId(gpuId)
}

// GetDeviceTopology returns device topology corresponding to the gpuId
func GetDeviceTopology(gpuId uint) ([]P2PLink, error) {
	return getDeviceTopology(gpuId)
}

// ListenForPolicyViolationsForAllGPUs sets GPU usage and error policies and notifies in case of any violations on all GPUs
func ListenForPolicyViolationsForAllGPUs(ctx context.Context, params *PolicyConditionParams) (<-chan PolicyViolation, error) {
	groupId := GroupAllGPUs()
	return registerPolicy(ctx, groupId, params)
}

// ListenForPolicyViolationsForGPUs sets GPU usage and error policies and notifies in case of any violations on special GPUs
func ListenForPolicyViolationsForGPUs(ctx context.Context, params *PolicyConditionParams, gpuIds ...uint) (<-chan PolicyViolation, error) {
	return registerPolicyForGpus(ctx, params, gpuIds...)
}
