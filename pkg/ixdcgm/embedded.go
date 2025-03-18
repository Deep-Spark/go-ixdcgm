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
#cgo linux LDFLAGS: -ldl -Wl,--export-dynamic -Wl,--unresolved-symbols=ignore-in-object-files
#include "include/dcgm_agent.h"
#include "include/dcgm_structs.h"
*/
import "C"
import "fmt"

type embedded struct {
}

func (e *embedded) Shutdown() error {
	result := C.dcgmStopEmbedded(handle.handle)
	if err := errorString(result); err != nil {
		return fmt.Errorf("failed to stop embedded dcgm: %v", err)
	}

	result = C.dcgmShutdown()
	if err := errorString(result); err != nil {
		return fmt.Errorf("failed to shutdown dcgm: %v", err)
	}
	return nil
}

func (e *embedded) Start(args ...string) (DcgmHandle, error) {
	fmt.Println("Start ixdcgm based on Embedded mode.")

	result := C.dcgmInit()
	if err := errorString(result); err != nil {
		return DcgmHandle{}, fmt.Errorf("failed to initialize dcgm: %v", err)
	}

	logLevel := C.DcgmLoggingSeverityNone
	if len(args) > 0 {
		logLevelStr := args[0]
		switch logLevelStr {
		case "LogNone":
			logLevel = C.DcgmLoggingSeverityNone
		case "LogFatal":
			logLevel = C.DcgmLoggingSeverityFatal
		case "LogError":
			logLevel = C.DcgmLoggingSeverityError
		case "LogWarn":
			logLevel = C.DcgmLoggingSeverityWarning
		case "LogInfo":
			logLevel = C.DcgmLoggingSeverityInfo
		case "LogDebug":
			logLevel = C.DcgmLoggingSeverityDebug
		case "LogVerb":
			logLevel = C.DcgmLoggingSeverityVerbose
		default:
			errMsg := fmt.Sprintf("Invalid log level: %s", logLevelStr)
			fmt.Println(errMsg)
			fmt.Println("The following log levels are supported: LogNone, LogFatal, LogError, LogWarn, LogInfo, LogDebug, LogVerb.")
			fmt.Println(" - LogNone  : No logging")
			fmt.Println(" - LogFatal : Fatal errors")
			fmt.Println(" - LogError : Errors")
			fmt.Println(" - LogWarn  : Warnings")
			fmt.Println(" - LogInfo  : Informative, will generate medium logs")
			fmt.Println(" - LogDebug : Debug infomation, will generate large logs")
			fmt.Println(" - LogVerb  : Verbose debugging information, will generate more large logs")
			fmt.Println()
			return DcgmHandle{}, fmt.Errorf("%v", errMsg)
		}
	}

	params := C.dcgmStartEmbeddedV2Params_v1{
		version:       C.dcgmStartEmbeddedV2Params_version1,
		opMode:        C.dcgmOperationMode_t(C.DCGM_OPERATION_MODE_AUTO),
		dcgmHandle:    C.dcgmHandle_t(0),
		logFile:       nil, // use default log file
		severity:      C.DcgmLoggingSeverity_t(logLevel),
		denyListCount: 0, // no deny list
		denyList:      [C.DcgmModuleIdCount]C.uint{0},
	}

	// Use dcgmStartEmbedded_v2 but dcgmStartEmbedded which using verbose log
	result = C.dcgmStartEmbedded_v2(&params)
	if err := errorString(result); err != nil {
		return DcgmHandle{}, fmt.Errorf("failed to start embedded dcgm: %v", err)
	}

	var cHandler C.dcgmHandle_t = params.dcgmHandle
	return DcgmHandle{handle: cHandler}, nil
}
