package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"

	"gitee.com/deep-spark/go-ixdcgm/pkg/ixdcgm"
)

const (
	HealthStatus = `GPU                : {{.GPU}}
Status             : {{.Status}}
{{range .Watches}}
Type               : {{.Type}}
Status             : {{.Status}}
Error              : {{.Error}}
{{end}}
`
)

// Based on ixdcgmi health commands:
// - Create group: ixdcgmi group -c <groupName>
// - Enable all watches: ixdcgmi health -g GROUPID -s a
// - Check: ixdcgmi health -g GROUPID -c
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	gpuIds, err := ixdcgm.GetSupportedDevices()
	if err != nil {
		log.Panicln(err)
	}

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	t := template.Must(template.New("HealthStatus").Parse(HealthStatus))
	for {
		select {
		case <-ticker.C:
			for _, gpuId := range gpuIds {
				h, err := ixdcgm.HealthCheckByGpuId(gpuId)
				if err != nil {
					log.Panicln(err)
				}

				if err = t.Execute(os.Stdout, h); err != nil {
					log.Panicln("Template error: ", err)
				}
			}
		case <-sigs:
			return
		}
	}

}
