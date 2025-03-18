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

import (
	"fmt"
	"log"
	"sync"
	"unsafe"
)

var (
	ixdcgmLibHandler  unsafe.Pointer
	ixdcgmInitCounter int
	connection        Interface
)

var (
	uptDirMu     sync.Mutex
	ixdcgmBinDir = "/usr/local/ixdcgm/bin"
	ixdcgmLibDir = "/usr/local/ixdcgm/lib64"
)

const ixdcgmLib = "libixdcgm.so"

func initIxDcgm(m int) (err error) {
	lib := string2Char(ixdcgmLib)
	defer freeCString(lib)

	ixdcgmLibHandler = C.dlopen(lib, C.RTLD_LAZY|C.RTLD_GLOBAL)
	if ixdcgmLibHandler == nil {
		errMsg := C.GoString(C.dlerror())
		log.Printf("failed to load %s from system library path: %s\ntry to load from %s\n",
			ixdcgmLib, errMsg, ixdcgmLibDir)

		abslib := string2Char(ixdcgmLibDir + "/" + ixdcgmLib)
		defer freeCString(abslib)
		ixdcgmLibHandler = C.dlopen(abslib, C.RTLD_LAZY|C.RTLD_GLOBAL)
	}
	if ixdcgmLibHandler == nil {
		errMsg := C.GoString(C.dlerror())
		return fmt.Errorf("failed to load %s, err: %s", ixdcgmLib, errMsg)
	}

	connection, err = New(m)
	if err != nil {
		return err
	}
	return nil
}

func shutdown() (err error) {
	mux.Lock()
	defer mux.Unlock()
	if ixdcgmInitCounter <= 0 {
		return fmt.Errorf("ixdcgm already shutdown")
	}

	if ixdcgmInitCounter == 1 {
		err = connection.Shutdown()
		if err != nil {
			return err
		}
	}

	C.dlclose(ixdcgmLibHandler)
	ixdcgmInitCounter -= 1
	return nil
}

func SetIxDcgmBinDir(dir string) error {
	uptDirMu.Lock()
	defer uptDirMu.Unlock()

	path, err := parseDirPath(dir)
	if err != nil {
		return err
	}
	ixdcgmBinDir = path
	return nil
}

func SetIxDcgmLibDir(dir string) error {
	uptDirMu.Lock()
	defer uptDirMu.Unlock()

	path, err := parseDirPath(dir)
	if err != nil {
		return err
	}
	ixdcgmLibDir = path
	return nil
}
