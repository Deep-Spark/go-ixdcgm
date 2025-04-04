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

import "C"
import (
	"context"
	"fmt"
	"sync"

	_ "gitee.com/deep-spark/go-ixdcgm/pkg/ixdcgm/include"
)

var (
	mux    sync.Mutex
	handle DcgmHandle
)

// Init starts IXDCGM, based on the user selected mode
// IXDCGM can be started in 3 differengt modes:
// 1. Embedded: Start hostengine within this process
// 2. Standalone: Connect to an already running ix-hostengine at the specified address
// Connection address can be passed as command line args: -connect "IP:PORT/Socket" -socket "isSocket"
// 3. StartHostengine: Open an Unix socket to start and connect to the ix-hostengine and terminate before exiting
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
