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
	"os"
	"os/signal"
	"syscall"

	"gitee.com/deep-spark/go-ixdcgm/pkg/ixdcgm"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	cleanup, err := ixdcgm.Init(ixdcgm.Embedded)
	if err != nil {
		log.Panicln(err)
	}
	defer cleanup()

	gpuIds, err := ixdcgm.GetSupportedDevices()
	if err != nil {
		log.Panicln(err)
	}

	for _, gpuId := range gpuIds {
		st, err := ixdcgm.GetDeviceStatus(gpuId)
		if err != nil {
			log.Panicln(err)
		}
		pst, err := ixdcgm.GetDeviceProfStatus(gpuId)
		if err != nil {
			log.Panicln(err)
		}

		fmt.Printf("GPUId                  : %d\n", st.Id)
		fmt.Printf("Power Usage (W)        : %s\n", st.Power)
		fmt.Printf("Temperature (°C)       : %s\n", st.Temperature)
		fmt.Printf("FanSpeed (%%)           : %s\n", st.FanSpeed)
		fmt.Printf("Utilization.GPU (%%)    : %d\n", st.Utilization.Gpu)
		fmt.Printf("Utilization.Mem (%%)    : %d\n", st.Utilization.Mem)
		fmt.Printf("Clocks.Cores (MHz)     : %d\n", st.Clocks.Sm)
		fmt.Printf("Clocks.Mem (MHz)       : %d\n", st.Clocks.Mem)
		fmt.Printf("EccSdbVolDev           : %s\n", st.EccSbeVolDev)
		fmt.Printf("EccDdbVolDev           : %s\n", st.EccDbeVolDev)
		fmt.Printf("PCI.Tx (MB/s)          : %d\n", st.PCI.Tx)
		fmt.Printf("PCI.Rx (MB/s)          : %d\n", st.PCI.Rx)
		fmt.Printf("PCI.ReplayCounter      : %d\n", st.PCI.ReplayCounter)
		fmt.Printf("Total Memory (MB)      : %d\n", st.MemUsage.Total)
		fmt.Printf("Used Memory (MB)       : %d\n", st.MemUsage.Used)
		fmt.Printf("Free Memory (MB)       : %d\n", st.MemUsage.Free)
		fmt.Printf("SmActive               : %s\n", pst.SmActive)
		fmt.Printf("SmOccupancy            : %s\n", pst.SmOccupancy)
		fmt.Printf("DramActive             : %s\n", pst.DramActive)
		fmt.Println("-------------------------------------------")
	}

}
