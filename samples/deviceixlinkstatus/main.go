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

package main

import (
	"fmt"
	"log"
	"sort"

	"gitee.com/deep-spark/go-ixdcgm/pkg/ixdcgm"
)

func main() {
	// Choose ixdcgm hostengine running mode
	// 1. ixdcgm.Embedded
	// 2. ixdcgm.Standalone -connect "addr", -socket "isSocket"
	// 3. ixdcgm.StartHostengine
	cleanup, err := ixdcgm.Init(ixdcgm.Embedded)
	if err != nil {
		log.Fatalf("failed to initialize ixdcgm: %v", err)
	}
	defer cleanup()

	gpuIds, err := ixdcgm.GetSupportedDevices()
	if err != nil {
		log.Fatalf("failed to get supported devices: %v", err)
	}

	// Test case 1: get global IXLink link status for all GPUs and NvSwitches
	linkStatus, err := ixdcgm.GetIxlinkStatus()
	if err != nil {
		log.Fatalf("failed to get global ixlink status: %v", err)
	}
	fmt.Printf("Global IXLink status: gpu_entries=%d nvswitch_entries=%d\n",
		len(linkStatus.GPUs), len(linkStatus.IxSwitches))

	// Test case 2: check if IXLink is supported on current system
	fmt.Printf("DeviceIxlinkSupported  : %t\n", ixdcgm.DeviceIxlinkSupported())
	fmt.Println("===========================================")

	for _, gpuId := range gpuIds {
		// Test case 3: get per-GPU IXLink link status
		gpuLinkStatus, found, err := ixdcgm.GetDeviceIxlinkStatus(gpuId)
		if err != nil {
			log.Fatalf("failed to get ixlink status for gpu %d: %v", gpuId, err)
		}
		if found {
			var up, down, disabled, unsupported int
			for _, state := range gpuLinkStatus.LinkState {
				switch state {
				case ixdcgm.IxlinkLinkStateUp:
					up++
				case ixdcgm.IxlinkLinkStateDown:
					down++
				case ixdcgm.IxlinkLinkStateDisabled:
					disabled++
				case ixdcgm.IxlinkLinkStateNotSupported:
					unsupported++
				}
			}
			fmt.Printf("GPUId                  : %d\n", gpuId)
			fmt.Printf("IXLINK LinkState       : up=%d down=%d disabled=%d unsupported=%d\n",
				up, down, disabled, unsupported)
		} else {
			fmt.Printf("GPUId                  : %d\n", gpuId)
			fmt.Println("IXLINK LinkState       : not found in GetIxlinkLinkStatus")
		}

		ixlinkErrorStatus, supported, err := ixdcgm.GetDeviceIxlinkErrorStatus(gpuId)
		if err != nil {
			log.Fatalf("failed to get ixlink error status for gpu %d: %v", gpuId, err)
		}

		fmt.Printf("IXLINK Supported       : %t\n", supported)
		if !supported {
			continue
		}

		metricNames := make([]string, 0, len(ixlinkErrorStatus.Values))
		for metricName := range ixlinkErrorStatus.Values {
			metricNames = append(metricNames, metricName)
		}
		sort.Strings(metricNames)

		for _, metricName := range metricNames {
			fmt.Printf("%-34s: %s\n", metricName, ixlinkErrorStatus.Values[metricName])
		}
		fmt.Println("-------------------------------------------")
	}
}
