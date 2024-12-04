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

#include <dlfcn.h>
#include <stdlib.h>
#include "include/dcgm_agent.h"
#include "include/dcgm_structs.h"
*/
import "C"
import "fmt"

type embedded struct {
	// TODO: implement embeded mode.
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
	result := C.dcgmInit()
	if err := errorString(result); err != nil {
		return DcgmHandle{}, fmt.Errorf("failed to initialize dcgm: %v", err)
	}

	var cHandler C.dcgmHandle_t
	result = C.dcgmStartEmbedded(C.DCGM_OPERATION_MODE_AUTO, &cHandler)
	if err := errorString(result); err != nil {
		return DcgmHandle{}, fmt.Errorf("failed to start embedded dcgm: %v", err)
	}

	return DcgmHandle{handle: cHandler}, nil
}
