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
	"fmt"
	"sync"
	"unsafe"
)

var (
	ixdcgmLibHandler  unsafe.Pointer
	ixdcgmInitCounter int
	mux               sync.Mutex
	connectionsMode   Interface
	handle            DcgmHandle
)

// to do path
const (
	ixdcgmLib = "libixdcgm.so"
)

func initIxDcgm(m int, args ...string) (err error) {
	lib := string2Char(ixdcgmLib)
	defer freeCString(lib)

	ixdcgmLibHandler = C.dlopen(lib, C.RTLD_LAZY|C.RTLD_GLOBAL)
	if ixdcgmLibHandler == nil {
		return fmt.Errorf("failed to load %s", ixdcgmLib)
	}

	connectionsMode, err = New(m)
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
		err = initIxDcgm(m, args...)
		if err != nil {
			return nil, err
		}

		handle, err = connectionsMode.Start(args...)
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
		err = connectionsMode.Shutdown()
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

func GetDeviceInfo(gpuId uint) (DeviceInfo, error) {
	return getDeviceInfo(gpuId)
}

// GetDeviceInfo describes the given device
func GetSupportedDevices() ([]uint, error) {
	return getSupportedDevices()
}

// GetDeviceStatus monitors GPU status including its power, memory and GPU utilization
func GetDeviceStatus(gpuId uint) (DeviceStatus, error) {
	return getDeviceStatus(gpuId)
}

func GetDeviceOnSameBoard(gpuId1, gpuId2 uint) (bool, error) {
	return getDeviceOnSameBoard(gpuId1, gpuId2)
}
