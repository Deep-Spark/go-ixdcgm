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
	"time"

	"gitee.com/deep-spark/go-ixdcgm/pkg/ixdcgm"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cleanup, err := ixdcgm.Init(ixdcgm.Embedded)
	if err != nil {
		log.Panicln(err)
	}
	defer cleanup()

	gpus, err := ixdcgm.GetSupportedDevices()
	if err != nil {
		log.Panicln(err)
	}

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, gpu := range gpus {
				st, err := ixdcgm.GetDeviceStatus(gpu)
				if err != nil {
					log.Panicln(err)
				}
				fmt.Printf("st.Power Usage %f\n", st.Power)
				fmt.Printf("st.Temperature %d\n", st.Temperature)
				fmt.Printf("st.Utilization.GPU %d\n", st.Utilization.GPU)
				fmt.Printf("st.Utilization.Memory %d\n", st.Utilization.Memory)
				fmt.Printf("st.Clocks.Cores %d\n", st.Clocks.Cores)
				fmt.Printf("st.Clocks.Memory %d\n", st.Clocks.Memory)
				fmt.Printf("st.FanSpeed %d\n", st.FanSpeed)
				fmt.Printf("st.st.PCI.Tx %d, st.PCI.Rx %d, st.PCI.Replays %d\n", st.PCI.Tx, st.PCI.Rx, st.PCI.Replays)
			}

		case <-sigs:
			return
		}
	}
}
