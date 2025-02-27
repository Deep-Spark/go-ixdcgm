package main

import (
	"fmt"
	"log"
	"strconv"

	"gitee.com/deep-spark/go-ixdcgm/pkg/ixdcgm"
)

const (
	legend = `
Legend:
  X    = Self
  SYS  = Connection traversing PCIe as well as the SMP interconnect between NUMA nodes (e.g., QPI/UPI)
  NODE = Connection traversing PCIe as well as the interconnect between PCIe Host Bridges within a NUMA node
  PHB  = Connection traversing PCIe as well as a PCIe Host Bridge (typically the CPU)
  PXB  = Connection traversing multiple PCIe bridges (without traversing the PCIe Host Bridge)
  PIX  = Connection traversing at most a single PCIe bridge
  INTE = Connection traversing at most a single on-board PCIe bridge
  IX#  = Connection traversing a bonded set of # IXLinks`
)

// Based on topo commands of ixdcgmi and ixsmi
func main() {
	// Choose ixdcgm hostengine running mode
	// 1. ixdcgm.Embedded
	// 2. ixdcgm.Standalone -connect "addr", -socket "isSocket"
	// 3. ixdcgm.StartHostengine
	cleanup, err := ixdcgm.Init(ixdcgm.Embedded)
	if err != nil {
		log.Panicln(err)
	}
	defer cleanup()

	gpus, err := ixdcgm.GetSupportedDevices()
	if err != nil {
		log.Panicln(err)
	}

	fmt.Printf("%-8s", "")
	for _, gpu := range gpus {
		fmt.Printf("%-8s", "GPU"+strconv.Itoa(int(gpu)))
	}
	fmt.Printf("%-16s", "CPU Affinity")
	fmt.Printf("%-16s\n", "NUMA Affinity")

	numGpus := len(gpus)
	gpuTopo := make([]string, numGpus)
	for i := 0; i < numGpus; i++ {
		topo, err := ixdcgm.GetDeviceTopology(gpus[i])
		if err != nil {
			log.Panicln(err)
		}

		fmt.Printf("%-8s", "GPU"+strconv.Itoa(int(gpus[i])))
		for j := 0; j < len(topo); j++ {
			// skip current GPU
			gpuTopo[topo[j].GPU] = topo[j].Link.PCIPaths()
		}
		gpuTopo[i] = " X "
		for j := 0; j < numGpus; j++ {
			fmt.Printf("%-8s", gpuTopo[j])
		}
		deviceInfo, err := ixdcgm.GetDeviceInfo(gpus[i])
		if err != nil {
			log.Panicln(err)
		}
		fmt.Printf("%-16s", deviceInfo.CPUAffinity)
		fmt.Printf("%-16s\n", deviceInfo.NUMAAffinity)
	}
	fmt.Println(legend)
}
