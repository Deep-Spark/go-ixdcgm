# Go-IXDCGM

## Introduction

IXDCGM is a tool provided for monitoring and managing IX GPUs, offering a rich set of APIs to retrieve information about GPU status, performance, power consumption, and more.   
Go-IXDCGM is a wrapper library for IXDCGM written in Go language, providing a simple set of functions that facilitate the easy invocation of IXDCGM's APIs.

**Note:** 
- The runtime environment requires the library of **libixdcgm.so**, please install IXDCGM SDK firstly.
- The current version of Go-IXDCGM is compatible with IX driver version **4.2.0**.

## Install

The installation of go-ixdcgm is very simple, just execute the following command in the command line：

```bash
go get gitee.com/deep-spark/go-ixdcgm
```

## Sample

A simple example of go-ixdcgm for getting device info is under:

```go
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
---------------------------------------------------------------------
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
			log.Panicln("Template error:", err)
		}
	}
}
```

## IXDCGM running modes
IXDCGM can be run in three different ways.

#### Embedded Mode
In embedded mode, hostengine is started as part of the running process and is loaded as a shared library. In this mode, metrics are also updated and collected automatically. This mode is recommended for users who wants to avoid managing an autonomous hostengine.

#### Standalone Mode
This mode allows you to connect to a running hostengine using a specified TCP/IP or Unix socket address. It is recommended for remote connections to the hostengine. By default, IXDCGM assumes a TCP connection and attempts to connect to localhost:5777, unless specified otherwise.

If the hostengine is running at a different address, pass it to `-connect`:
- "IP" - A valid IP address for the remote hostengine, at port 5777.
- "IP:PORT" - A valid IP address and port.

The `-socket` parameter identifies whether the passed `-connect` address is a Unix socket filename (1) or a TCP/IP address (0):
- "0" - The given address is a TCP/IP address.
- "1" - The given address is a Unix socket filename.

For example:
```
go run samples/devicecommon/main.go -connect "0.0.0.0:5777" -socket "0"
```

#### StartHostengine Mode
This is an add-on mode which opens an Unix socket for starting and connecting with hostengine. The hostengine is started as a child process of the running process and automatically terminated on exit. When operating in this mode, make sure to stop an already running hostengine to avoid any connection address conflicts. This mode is recommended for safely integrating IXDCGM in an already existing setup.

## More Samples

The `samples` folder contains more simple examples of how to use go-ixdcgm to call the IXDCGM API.

To get device information, run the following command:
```
$ go run samples/deviceinfo/main.go

# sample output

Driver Version         : 4.2.0
GPUId                  : 1
IxDCGMSupported        : Y
Uuid                   : GPU-6d2ec5fa-f293-57a3-9f2c-335f78120578
Product Name           : Iluvatar BI-V150S
Serial Number          : 24120026944896
Bus ID                 : 00000000:8A:00.0
BAR1 (MB)              : N/A
Total Memory (MB):     : 32768
Used Memory (MB):      : 25500
Free Memory (MB):      : 7268
Bandwidth (MB/s)       : 31504
PowerLimit (W)         : 205
CPUAffinity            : 20-39,60-79
NUMAAffinity           : 1
P2P Available          :
    GPU0 - (BusID)00000000:8A:00.0 - SYS
    GPU2 - (BusID)00000000:8A:00.0 - INTE
--------------------------------------------------	
```

To get device status, run the following command:
```
$ go run samples/devicestatus/main.go

# sample output

GPUId                  : 1
Power Usage (W)        : 150.000
Temperature (°C)       : 68
FanSpeed (%)           : N/A
Utilization.GPU (%)    : 85
Utilization.Mem (%)    : 78
Clocks.Cores (MHz)     : 1750
Clocks.Mem (MHz)       : 1600
EccSdbVolDev           : 0
EccDdbVolDev           : 0
PCI.Tx (MB/s)          : 107
PCI.Rx (MB/s)          : 92544
PCI.ReplayCounter      : 0
Total Memory (MB)      : 32768
Used Memory (MB)       : 25500
Free Memory (MB)       : 7268
SmActive               : 0.792
SmOccupancy            : 0.222
DramActive             : 0.622
-------------------------------------------
```

To get running process information of device, run the following command:
```
$ go run samples/deviceprocessinfo/main.go

# sample output

Get the running process infos of gpu 1
> Pid: 4009629
  Name: ./gemm_perf --i 1 --d 0 --m 1024 --l 2000
  UsedGpuMemory(MiB): 128
```

To monitor device health iteratively, run the following command:
```
$ go run samples/health/main.go

# sample output

GPU                : 0
Status             : Healthy

GPU                : 1
Status             : Healthy

GPU                : 2
Status             : Healthy

...
```
`Note`: Press Ctrl+C to stop the iteration output.  

To find the topology of GPUs on the system, run the following command:
```
$ go run samples/topology/main.go

# sample output

        GPU0    GPU1    GPU2    CPU Affinity    NUMA Affinity
GPU0     X      SYS     SYS     0-19,40-59      0
GPU1    SYS      X      INTE    20-39,60-79     1
GPU2    SYS     INTE     X      20-39,60-79     1

Legend:
  X    = Self
  SYS  = Connection traversing PCIe as well as the SMP interconnect between NUMA nodes (e.g., QPI/UPI)
  NODE = Connection traversing PCIe as well as the interconnect between PCIe Host Bridges within a NUMA node
  PHB  = Connection traversing PCIe as well as a PCIe Host Bridge (typically the CPU)
  PXB  = Connection traversing multiple PCIe bridges (without traversing the PCIe Host Bridge)
  PIX  = Connection traversing at most a single PCIe bridge
  INTE = Connection traversing at most a single on-board PCIe bridge
  IX#  = Connection traversing a bonded set of # IXLinks
```

To set violation policy and monitor policy violations iteratively, run the following command:
```
$ go run samples/policy/main.go

# sample output

2025/02/25 17:05:22 Policy successfully set.
2025/02/25 17:05:22 Listening for violations...
PolicyViolation : Thermal Limit
Timestamp       : 2025-02-25 17:05:42 +0800 CST
Data            : {61}
PolicyViolation : Thermal Limit
Timestamp       : 2025-02-25 17:05:42 +0800 CST
Data            : {61}
...
```
`Note`: Press Ctrl+C to stop the iteration output.


To execute a diagnostic on the system and devices, run the following command:
```
$ go run samples/diag/main.go

# sample output

Software:

  presence of drivers on the denylist (e.g. nouveau) pass
  presence (and version) of NVML lib                 pass
  presence (and version) of CUDA lib                 pass
  presence (and version) of CUDA RT lib              pass
  character device permissions                       pass

GPU: 0

  Memory               notrun
  Diagnostic           notrun
  PCIe                 notrun
  SM Stress            notrun
  Targeted Stress      notrun
  Targeted Power       notrun
  Memory Bandwidth     notrun
  Memtest              notrun
  Pulse Test           notrun
  EUD Test             notrun
  CPU EUD Test         notrun

```

## License

Copyright (c) 2024 Iluvatar CoreX. All rights reserved. This project has an Apache-2.0 license, as
found in the [LICENSE](LICENSE) file.
