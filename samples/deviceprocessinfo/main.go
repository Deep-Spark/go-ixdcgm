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

	cleanup, err := ixdcgm.Init(ixdcgm.Embedded, "LogInfo")
	if err != nil {
		log.Panicln(err)
	}
	defer cleanup()

	gpuIds, err := ixdcgm.GetSupportedDevices()
	if err != nil {
		log.Panicln(err)
	}

	for _, gpuId := range gpuIds {
		fmt.Printf("Get the running process infos of gpu %d\n", gpuId)
		infos, err := ixdcgm.GetDeviceRunningProcesses(gpuId)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		for _, info := range infos {
			fmt.Printf("> Pid: %d, Name: %s, UsedGpuMemory(MiB): %d\n", info.Pid, info.Name, info.UsedGpuMemory)
		}
		fmt.Println("---------------------------------------------------------------------")
	}

}
