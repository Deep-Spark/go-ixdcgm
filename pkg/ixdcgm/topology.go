/*
Copyright (c) 2024, NVIDIA CORPORATION.
Copyright (c) 2024, Shanghai Iluvatar CoreX Semiconductor Co., Ltd.
All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

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

type P2PLinkType uint

const (
	P2PLinkUnknown P2PLinkType = iota // N/A

	P2PLinkCrossCPU     // SYS  = Connection traversing PCIe as well as the SMP interconnect between NUMA nodes (e.g., QPI/UPI)
	P2PLinkSameCPU      // NODE = Connection traversing PCIe as well as the interconnect between PCIe Host Bridges within a NUMA node
	P2PLinkHostBridge   // PHB  = Connection traversing PCIe as well as a PCIe Host Bridge (typically the CPU)
	P2PLinkMultiSwitch  // PXB  = Connection traversing multiple PCIe bridges (without traversing the PCIe Host Bridge)
	P2PLinkSingleSwitch // PIX  = Connection traversing at most a single PCIe bridge
	P2PLinkSameBoard    // INTE = Connection traversing at most a single on-board PCIe bridge
	P2PLinkIXLINK1      // IX1  = Connection traversing a single IXLink
	P2PLinkIXLINK2      // IX2  = Connection traversing two IXLinks
	P2PLinkIXLINK3      // IX3  = Connection traversing three IXLinks
	P2PLinkIXLINK4      // IX4  = Connection traversing four IXLinks
	P2PLinkIXLINK5      // IX5  = Connection traversing five IXLinks
	P2PLinkIXLINK6      // IX6  = Connection traversing six IXLinks
)

func (l P2PLinkType) PCIPaths() string {
	switch l {
	case P2PLinkSameBoard:
		return "INTE"
	case P2PLinkSingleSwitch:
		return "PIX"
	case P2PLinkMultiSwitch:
		return "PXB"
	case P2PLinkHostBridge:
		return "PHB"
	case P2PLinkSameCPU:
		return "NODE"
	case P2PLinkCrossCPU:
		return "SYS"
	case P2PLinkIXLINK1:
		return "IX1"
	case P2PLinkIXLINK2:
		return "IX2"
	case P2PLinkIXLINK3:
		return "IX3"
	case P2PLinkIXLINK4:
		return "IX4"
	case P2PLinkIXLINK5:
		return "IX5"
	case P2PLinkIXLINK6:
		return "IX6"
	case P2PLinkUnknown:
	}
	return "N/A"
}

type P2PLink struct {
	GPU   uint
	BusID string
	Link  P2PLinkType
}

func getP2PLink(path uint) P2PLinkType {
	switch path {
	case C.DCGM_TOPOLOGY_BOARD:
		return P2PLinkSameBoard
	case C.DCGM_TOPOLOGY_SINGLE:
		return P2PLinkSingleSwitch
	case C.DCGM_TOPOLOGY_MULTIPLE:
		return P2PLinkMultiSwitch
	case C.DCGM_TOPOLOGY_HOSTBRIDGE:
		return P2PLinkHostBridge
	case C.DCGM_TOPOLOGY_CPU:
		return P2PLinkSameCPU
	case C.DCGM_TOPOLOGY_SYSTEM:
		return P2PLinkCrossCPU
	case C.DCGM_TOPOLOGY_NVLINK1:
		return P2PLinkIXLINK1
	case C.DCGM_TOPOLOGY_NVLINK2:
		return P2PLinkIXLINK2
	case C.DCGM_TOPOLOGY_NVLINK3:
		return P2PLinkIXLINK3
	case C.DCGM_TOPOLOGY_NVLINK4:
		return P2PLinkIXLINK4
	case C.DCGM_TOPOLOGY_NVLINK5:
		return P2PLinkIXLINK5
	case C.DCGM_TOPOLOGY_NVLINK6:
		return P2PLinkIXLINK6
	}
	return P2PLinkUnknown
}

func getBusid(gpuid uint) (string, error) {
	var device C.dcgmDeviceAttributes_v3
	device.version = makeVersion3(unsafe.Sizeof(device))

	result := C.dcgmGetDeviceAttributes(handle.handle, C.uint(gpuid), &device)
	if err := errorString(result); err != nil {
		return "", fmt.Errorf("Error getting device busid: %s", err)
	}
	return *stringPtr(&device.identifiers.pciBusId[0]), nil
}

func getDeviceTopology(gpuid uint) (links []P2PLink, err error) {
	var topology C.dcgmDeviceTopology_v1
	topology.version = makeVersion1(unsafe.Sizeof(topology))

	result := C.dcgmGetDeviceTopology(handle.handle, C.uint(gpuid), &topology)
	if result == C.DCGM_ST_NOT_SUPPORTED {
		return links, fmt.Errorf("DcgmGetDeviceTopology is not supported")
	}
	if result != C.DCGM_ST_OK {
		return links, &DcgmError{msg: C.GoString(C.errorString(result)), Code: result}
	}

	busid, err := getBusid(gpuid)
	if err != nil {
		return
	}

	for i := uint(0); i < uint(topology.numGpus); i++ {
		gpu := topology.gpuPaths[i].gpuId
		p2pLink := P2PLink{
			GPU:   uint(gpu),
			BusID: busid,
			Link:  getP2PLink(uint(topology.gpuPaths[i].path)),
		}
		links = append(links, p2pLink)
	}
	return
}
