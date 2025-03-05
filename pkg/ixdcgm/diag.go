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

package ixdcgm

/*
#include "include/dcgm_agent.h"
#include "include/dcgm_structs.h"
*/
import "C"
import (
	"context"
	"fmt"
	"time"
	"unsafe"
)

const (
	DCGM_PER_GPU_TEST_COUNT_V8 = 13
	DIAG_RESULT_STRING_SIZE    = 1024

	UNDEFINED_SWTEST = "UNDEFINED_SWTEST"
	UNUSED_GPUTEST   = "UNUSED_GPUTEST"
)

type DiagType int

const (
	DiagQuick    DiagType = 1 // run a very basic health check on the system
	DiagMedium   DiagType = 2 // run a medium-length diagnostic (a few minutes)
	DiagLong     DiagType = 3 // run a extensive diagnostic (several minutes)
	DiagExtended DiagType = 4 // run a very extensive diagnostic (many minutes)
)

type DiagResult struct {
	Status       string
	TestName     string
	TestOutput   string
	ErrorCode    uint
	ErrorMessage string
}

type GpuResult struct {
	GPU         uint
	RC          uint
	DiagResults []DiagResult
}

type DiagResults struct {
	Software []DiagResult
	PerGpu   []GpuResult
	gpuCount uint
}

func diagResultString(r int) string {
	switch r {
	case C.DCGM_DIAG_RESULT_PASS:
		return "pass"
	case C.DCGM_DIAG_RESULT_SKIP:
		return "skipped"
	case C.DCGM_DIAG_RESULT_WARN:
		return "warn"
	case C.DCGM_DIAG_RESULT_FAIL:
		return "fail"
	case C.DCGM_DIAG_RESULT_NOT_RUN:
		return "notrun"
	}
	return ""
}

func swTestName(t int) string {
	switch t {
	case C.DCGM_SWTEST_DENYLIST:
		return "presence of drivers on the denylist (e.g. nouveau)"
	case C.DCGM_SWTEST_NVML_LIBRARY:
		return "presence (and version) of NVML lib"
	case C.DCGM_SWTEST_CUDA_MAIN_LIBRARY:
		return "presence (and version) of CUDA lib"
	case C.DCGM_SWTEST_CUDA_RUNTIME_LIBRARY:
		return "presence (and version) of CUDA RT lib"
	case C.DCGM_SWTEST_PERMISSIONS:
		return "character device permissions"
	case C.DCGM_SWTEST_PERSISTENCE_MODE:
		return "persistence mode enabled"
	case C.DCGM_SWTEST_ENVIRONMENT:
		return "CUDA environment vars that may slow tests"
	case C.DCGM_SWTEST_PAGE_RETIREMENT:
		return "pending frame buffer page retirement"
	case C.DCGM_SWTEST_GRAPHICS_PROCESSES:
		return "graphics processes running"
	case C.DCGM_SWTEST_INFOROM:
		return "inforom corruption"
	}
	return UNDEFINED_SWTEST
}

func gpuTestName(t int) string {
	switch t {
	case C.DCGM_MEMORY_INDEX:
		return "Memory"
	case C.DCGM_DIAGNOSTIC_INDEX:
		return "Diagnostic"
	case C.DCGM_PCI_INDEX:
		return "PCIe"
	case C.DCGM_SM_STRESS_INDEX:
		return "SM Stress"
	case C.DCGM_TARGETED_STRESS_INDEX:
		return "Targeted Stress"
	case C.DCGM_TARGETED_POWER_INDEX:
		return "Targeted Power"
	case C.DCGM_MEMORY_BANDWIDTH_INDEX:
		return "Memory Bandwidth"
	case C.DCGM_MEMTEST_INDEX:
		return "Memtest"
	case C.DCGM_PULSE_TEST_INDEX:
		return "Pulse Test"
	case C.DCGM_EUD_TEST_INDEX:
		return "EUD Test"
	case C.DCGM_UNUSED3_TEST_INDEX:
		return "CPU EUD Test"
	case C.DCGM_SOFTWARE_INDEX:
		return "Software"
	case C.DCGM_CONTEXT_CREATE_INDEX:
		return "Context Create"
	}
	return UNUSED_GPUTEST
}

func newDiagResult(testResult C.dcgmDiagTestResult_v3, testName string) DiagResult {
	msg := C.GoString((*C.char)(unsafe.Pointer(&testResult.error[0].msg)))
	info := C.GoString((*C.char)(unsafe.Pointer(&testResult.info)))

	return DiagResult{
		Status:       diagResultString(int(testResult.status)),
		TestName:     testName,
		TestOutput:   info,
		ErrorCode:    uint(testResult.error[0].code),
		ErrorMessage: msg,
	}
}

func diagLevel(diagType DiagType) C.dcgmDiagnosticLevel_t {
	switch diagType {
	case DiagQuick:
		return C.DCGM_DIAG_LVL_SHORT
	case DiagMedium:
		return C.DCGM_DIAG_LVL_MED
	case DiagLong:
		return C.DCGM_DIAG_LVL_LONG
	case DiagExtended:
		return C.DCGM_DIAG_LVL_XLONG
	}
	return C.DCGM_DIAG_LVL_INVALID
}

// RunDiagWithTimeout executes a diagnostic with a specified timeout.
// If the diagnostic does not complete within the timeout duration, an error is returned.
func RunDiagWithTimeout(diagType DiagType, groupId GroupHandle, t time.Duration) (DiagResults, error) {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	// Channels to receive the diagnostic results or error
	resultChan := make(chan DiagResults, 1)
	errChan := make(chan error, 1)

	// Run the diagnostic in a separate goroutine
	go func() {
		result, err := RunDiag(diagType, groupId)
		if err != nil {
			errChan <- err
		} else {
			resultChan <- result
		}
	}()

	// Wait for the diagnostic to complete or the timeout to occur
	select {
	case <-ctx.Done():
		return DiagResults{}, fmt.Errorf("Error: diagnostic execution timed out after %v", t)
	case err := <-errChan:
		return DiagResults{}, err
	case result := <-resultChan:
		return result, nil
	}
}

func RunDiag(diagType DiagType, groupId GroupHandle) (DiagResults, error) {
	var diagResults C.dcgmDiagResponse_v10
	diagResults.version = makeVersion10(unsafe.Sizeof(diagResults))

	result := C.dcgmRunDiagnostic(handle.handle, groupId.handle, diagLevel(diagType), (*C.dcgmDiagResponse_v10)(unsafe.Pointer(&diagResults)))
	if err := errorString(result); err != nil {
		return DiagResults{}, &DcgmError{msg: C.GoString(C.errorString(result)), Code: result}
	}
	defer C.dcgmStopDiagnostic(handle.handle)

	var diagRun DiagResults
	diagRun.gpuCount = uint(diagResults.gpuCount)

	for i := 0; i < int(diagResults.levelOneTestCount); i++ {
		testName := swTestName(i)
		if testName == UNDEFINED_SWTEST {
			continue
		}
		dr := newDiagResult(diagResults.levelOneResults[i], testName)
		diagRun.Software = append(diagRun.Software, dr)
	}

	for i := uint(0); i < uint(diagResults.gpuCount); i++ {
		r := diagResults.perGpuResponses[i]
		gr := GpuResult{GPU: uint(r.gpuId), RC: uint(r.hwDiagnosticReturn)}
		for j := 0; j < DCGM_PER_GPU_TEST_COUNT_V8; j++ {
			testName := gpuTestName(j)
			if testName == UNUSED_GPUTEST {
				continue
			}
			dr := newDiagResult(r.results[j], testName)
			gr.DiagResults = append(gr.DiagResults, dr)
		}
		diagRun.PerGpu = append(diagRun.PerGpu, gr)
	}

	return diagRun, nil
}
