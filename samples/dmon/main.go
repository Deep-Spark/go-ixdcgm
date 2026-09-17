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

const (
	header = `# gpu   pwr  temp    sm   mem  mclk  pclk
# Idx     W     C     %     %   MHz   MHz`
)

// modelled on ixsmi dmon
// ixdcgmi dmon -e 155,150,203,204,100,101
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

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	fmt.Println(header)
	for {
		select {
		case <-ticker.C:
			for _, gpu := range gpus {
				st, err := ixdcgm.GetDeviceStatus(gpu)
				if err != nil {
					log.Panicln(err)
				}

				fmt.Printf("%5d %.5s %5s %5d %5d %5d %5d\n",
					gpu, st.Power, st.GpuTemperature, st.Utilization.Gpu, st.Utilization.Mem,
					st.Clocks.Mem, st.Clocks.Sm)
			}

		case <-sigs:
			return
		}
	}
}
