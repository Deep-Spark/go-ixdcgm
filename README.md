# Go-IXDCGM

## Introduction

IXDCGM is a tool provided for monitoring and managing IX GPUs, offering a rich set of APIs to retrieve information about GPU status, performance, power consumption, and more. Go-IXDCGM is a wrapper library for IXDCGM written in Go language, providing a simple set of functions that facilitate the easy invocation of IXDCGM's APIs.

## Install

The installation of Go-IXDCGM is very simple, just execute the following command in the command line：

```bash
go get gitee.com/deep-spark/go-ixdcgm
```

## Samples

An example of go-ixdcgm for device-info is under:

```go
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
PowerLimit (W)         : {{or .PowerLimit "N/A"}}
---------------------------------------------------------------------
`
)

var (
	connectAddr = flag.String("connectAddr", "0.0.0.0:5777", "DCGM connect address")
	isSocket    = flag.String("socket", "0", "Connect to Unix socket")
)

func main() {
	// choose ixdcgm hostengine running mode
	// 1. ixdcgm.Embedded
	// 2. ixdcgm.Standalone -connect "addr", -socket "isSocket"
	// 3. ixdcgm.StartHostengine
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
```

## License

Copyright (c) 2024 Iluvatar CoreX. All rights reserved. This project has an Apache-2.0 license, as
found in the [LICENSE](LICENSE) file.
