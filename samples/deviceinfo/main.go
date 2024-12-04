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
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"

	"gitee.com/deep-spark/go-ixdcgm/pkg/ixdcgm"
)

const (
	deviceInfo = `Driver Version         : {{.Identifiers.DriverVersion}}
GPUId		       : {{.GPUId}}
IxDCGMSupported        : {{.IxDCGMSupported}}
Uuid                   : {{.Uuid}}
Product Name           : {{.Identifiers.ProductName}}
Serial Number          : {{.Identifiers.Serial}}
Bus ID                 : {{.PCI.BusId}}
BAR1 (MB)              : {{or .MemoryUsage.BAR1 "N/A"}}
Total Memory (MB):     : {{or .MemoryUsage.Total "N/A"}}
Used Memory (MB):      : {{or .MemoryUsage.Used "N/A"}}
Free Memory (MB):      : {{or .MemoryUsage.Free "N/A"}}
Bandwidth (MB/s)       : {{or .PCI.Bandwidth "N/A"}}
Power (W)              : {{or .Power "N/A"}}
---------------------------------------------------------------------
`
)

var (
	connectAddr = flag.String("connectAddr", "0.0.0.0:5777", "DCGM connect address")
	isSocket    = flag.String("socket", "0", "Connect to Unix socket")
)

func main() {
	// choose dcgm hostengine running mode
	// 1. dcgm.Embedded
	// 2. dcgm.Standalone -connect "addr", -socket "isSocket"
	// 3. dcgm.StartHostengine
	flag.Parse()
	cleanup, err := ixdcgm.Init(ixdcgm.Standalone, *connectAddr, *isSocket)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	gpuCount, err := ixdcgm.GetAllDeviceCount()
	if err != nil {
		panic(err)
	}
	fmt.Println("GPU Count:", gpuCount)
	t := template.Must(template.New("DeviceInfo").Parse(deviceInfo))

	for i := uint(0); i < gpuCount; i++ {
		d, err := ixdcgm.GetDeviceInfo(i)
		if err != nil {
			panic(err)
		}

		if err = t.Execute(os.Stdout, d); err != nil {
			log.Panicln("Template error:", err)
		}
	}
}
