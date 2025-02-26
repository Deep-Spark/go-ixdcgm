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
	"html/template"
	"log"
	"os"

	"gitee.com/deep-spark/go-ixdcgm/pkg/ixdcgm"
)

const (
	deviceInfo = `Driver Version         : {{.Identifiers.DriverVersion}}
GPUId                  : {{.GPUId}}
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
PowerLimit (W)         : {{or .PowerLimit "N/A"}}
CPUAffinity            : {{or .CPUAffinity "N/A"}}
NUMAAffinity           : {{or .NUMAAffinity "N/A"}}
P2P Available          : {{if not .Topology}}None{{else}}{{range .Topology}}
    GPU{{.GPU}} - (BusID){{.BusID}} - {{.Link.PCIPaths}}{{end}}{{end}}
--------------------------------------------------
`
)

func main() {
	// Choose ixdcgm hostengine running mode
	// 1. ixdcgm.Embedded
	// 2. ixdcgm.Standalone -connect "addr", -socket "isSocket"
	// 3. ixdcgm.StartHostengine
	cleanup, err := ixdcgm.Init(ixdcgm.Embedded)
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
			log.Panicln("Template error: ", err)
		}
	}
}
