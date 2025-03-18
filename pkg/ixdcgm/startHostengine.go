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
import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

var (
	hostengineAsChildPid int
	startHostengineDir   = "/tmp"
)

type startHostengine struct {
}

func (s *startHostengine) Shutdown() (err error) {
	if err = s.disconnect(); err != nil {
		return
	}

	// terminate ix-hostengine
	cmd := exec.Command("ix-hostengine", "--term")
	cmd.Env = append(os.Environ(),
		"PATH="+os.Getenv("PATH=")+":"+ixdcgmBinDir,
		"LD_LIBRARY_PATH="+os.Getenv("LD_LIBRARY_PATH")+":"+ixdcgmLibDir,
	)
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("Error terminating ix-hostengine: %s", err)
	}
	fmt.Println("Successfully terminated ix-hostengine.")

	return syscall.Kill(hostengineAsChildPid, syscall.SIGKILL)
}

func (s *startHostengine) disconnect() (err error) {
	result := C.dcgmDisconnect(handle.handle)
	if err = errorString(result); err != nil {
		return fmt.Errorf("Error disconnecting from ix-hostengine: %s", err)
	}

	result = C.dcgmShutdown()
	if err = errorString(result); err != nil {
		return fmt.Errorf("Error shutting down IXDCGM: %s", err)
	}
	return
}

func (s startHostengine) Start(args ...string) (DcgmHandle, error) {
	fmt.Println("Start ixdcgm based on StartHostengine mode.")

	os.Setenv("PATH", os.Getenv("PATH=")+":"+ixdcgmBinDir)
	bin, err := exec.LookPath("ix-hostengine")
	if err != nil {
		return DcgmHandle{}, fmt.Errorf("Error finding ix-hostengine: %s", err)
	}
	var procAttr syscall.ProcAttr
	procAttr.Files = []uintptr{
		uintptr(syscall.Stdin),
		uintptr(syscall.Stdout),
		uintptr(syscall.Stderr)}
	procAttr.Sys = &syscall.SysProcAttr{Setpgid: true}
	procAttr.Env = []string{"LD_LIBRARY_PATH=" + os.Getenv("LD_LIBRARY_PATH") + ":" + ixdcgmLibDir}

	dir := startHostengineDir
	socketFile, err := os.CreateTemp(dir, "ixdcgm")
	if err != nil {
		return DcgmHandle{}, fmt.Errorf("Error creating socket file in %s directory: %s", dir, err)
	}

	socketPath := socketFile.Name()
	defer os.Remove(socketPath)
	connectArg := "--domain-socket"
	hostengineAsChildPid, err = syscall.ForkExec(bin, []string{bin, connectArg, socketPath}, &procAttr)
	if err != nil {
		return DcgmHandle{}, fmt.Errorf("Error fork-execing ix-hostengine: %s", err)
	}
	result := C.dcgmInit()
	if err = errorString(result); err != nil {
		return DcgmHandle{}, fmt.Errorf("Error initializing IXDCGM: %s", err)
	}

	var cHandle C.dcgmHandle_t
	var connectParams C.dcgmConnectV2Params_v2
	connectParams.version = makeVersion2(unsafe.Sizeof(connectParams))
	connectParams.addressIsUnixSocket = C.uint(1)
	cSockPath := C.CString(socketPath)
	defer freeCString(cSockPath)

	result = C.dcgmConnect_v2(cSockPath, &connectParams, &cHandle)
	if err = errorString(result); err != nil {
		return DcgmHandle{}, fmt.Errorf("Error connecting to ix-hostengine: %s", err)
	}

	return DcgmHandle{handle: cHandle}, nil
}

func SetStartHostengineDir(dir string) error {
	uptDirMu.Lock()
	defer uptDirMu.Unlock()

	path, err := parseDirPath(dir)
	if err != nil {
		return err
	}
	startHostengineDir = path
	return nil
}
