package main

import (
	"fmt"
	"log"

	"gitee.com/deep-spark/go-ixdcgm/pkg/ixdcgm"
)

// ixdcgmi introspect -s -H
func main() {
	cleanup, err := ixdcgm.Init(ixdcgm.Embedded)
	if err != nil {
		log.Panicln(err)
	}
	defer cleanup()

	st, err := ixdcgm.Introspect()
	if err != nil {
		log.Panicln(err)
	}

	fmt.Printf("Memory %2s %v KB\nCPU %5s %.2f %s\n", ":", st.Memory, ":", st.CPU, "%")
}
