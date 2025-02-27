package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitee.com/deep-spark/go-ixdcgm/pkg/ixdcgm"
)

// Based on ixdcgmi policy commands:
// - Create group: ixdcgmi group -c <groupName>
// - Set violation policy: ixdcgmi policy -g GROUPID --set 0,0 -e -p -T 60
// - Register for policy updates: ixdcgmi policy -g GROUPID --reg
func main() {
	// Choose ixdcgm hostengine running mode
	// 1. ixdcgm.Embedded
	// 2. ixdcgm.Standalone -connect "addr", -socket "isSocket"
	// 3. ixdcgm.StartHostengine
	cleanup, err := ixdcgm.Init(ixdcgm.Embedded)
	if err != nil {
		log.Panicln(err)
	}
	defer func() {
		cleanup()
	}()

	ctx, done := context.WithCancel(context.Background())
	// Handle SIGINT (Ctrl+C) and SIGTERM (termination signal)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("Received termination signal, exiting...")
		done()
	}()

	// Create policy condition parameters to register violation callback.
	// Note: at least one policy must be enabled.
	params := &ixdcgm.PolicyConditionParams{
		DbePolicyEnabled:       true,
		PCIePolicyEnabled:      true,
		ThermalPolicyEnabled:   true,
		ThermalPolicyThreshold: 60, // °C
	}

	// Monitor policy violations for all GPUs
	ch, err := ixdcgm.ListenForPolicyViolationsForAllGPUs(ctx, params)

	// If you want to monitor policy violations for particular GPUs (e.g., gpuId0 and gpuId1),
	// use the following code:
	// ch, err := ixdcgm.ListenForPolicyViolationsForGPUs(ctx, params, 0, 1)

	if err != nil {
		fmt.Printf("Failed to monitor policy violations, err: %v", err)
		return
	}

	for {
		select {
		case pe := <-ch:
			fmt.Printf("PolicyViolation : %v\nTimestamp       : %v\nData            : %v\n",
				pe.Condition, pe.Timestamp, pe.Data)
		case <-ctx.Done():
			// Sleep to ensure the ixdcgm policy is unregistered before cleanup.
			time.Sleep(1 * time.Second)
			return
		}
	}
}
