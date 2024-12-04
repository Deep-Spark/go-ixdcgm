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

const (
	Embedded int = iota
	Standalone
	StartHostengine
)

type DcgmHandle struct {
	handle C.dcgmHandle_t
}

type Interface interface {
	Shutdown() (err error)
	Start(args ...string) (DcgmHandle, error)
}

func New(m int) (Interface, error) {
	switch m {
	case Embedded:
		return &embedded{}, nil
	case Standalone:
		return &standalone{}, nil
	case StartHostengine:
		return &startHostengine{}, nil
	default:
		return nil, fmt.Errorf("unknown mode: %d", m)
	}
}

var _ Interface = (*embedded)(nil)
var _ Interface = (*standalone)(nil)
var _ Interface = (*startHostengine)(nil)
