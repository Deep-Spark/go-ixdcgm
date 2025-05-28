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
		PowerPolicyEnabled:     true,
		PowerPolicyThreshold:   250, // W
		XidPolicyEnabled:       true,
	}

	// Monitor policy violations for all GPUs
	// Note: if you want to monitor policy violations for special GPUs (e.g., gpuId0 and gpuId1),
	// use the api: ixdcgm.ListenForPolicyViolationsForGPUs(ctx, params, gpuId0, gpuId1)
	ch, err := ixdcgm.ListenForPolicyViolationsForAllGPUs(ctx, params)
	if err != nil {
		fmt.Printf("Failed to monitor policy violations, err: %v", err)
		return
	}

	// Read the policy violations from the channel as soon as possible.
	for {
		select {
		case pe := <-ch:
			fmt.Printf("PolicyViolation : %v\nTimestamp       : %v\nData            : %+v\n",
				pe.Condition, pe.Timestamp, pe.Data)
		case <-ctx.Done():
			// Sleep to ensure the ixdcgm policy is unregistered before cleanup.
			time.Sleep(1 * time.Second)
			return
		}
	}
}
