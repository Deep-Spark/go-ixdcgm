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
#cgo LDFLAGS: -ldl
#include "include/dcgm_agent.h"
#include "include/ixdcgmStructs.h"
#include "include/ixdcgmApiExport.h"
*/
import "C"
import (
	"fmt"
	"math/rand"
	"sync"
	"unsafe"
)

type IxlinkErrorStatus struct {
	Id     uint
	Values map[string]string
}

type ixlinkErrorMetricDef struct {
	Name    string
	FieldID Short
}

// ixlinkErrorMetricDefs currently covers lane-level error counters for link l0-l5 only.
// This matches the currently supported IXLink topology on deployed GPUs.
// Note: IXDCGM_MAX_LINKS_PER_GPU can be 18 in headers, but l6-l17 error field IDs are
// not exposed in the current IXLink error metric set yet.
var ixlinkErrorMetricDefs = []ixlinkErrorMetricDef{
	{"ixlink_bad_tlp_error_count_l0", DCGM_FI_IXLINK_BAD_TLP_ERROR_COUNT_L0},
	{"ixlink_bad_tlp_error_count_l1", DCGM_FI_IXLINK_BAD_TLP_ERROR_COUNT_L1},
	{"ixlink_bad_tlp_error_count_l2", DCGM_FI_IXLINK_BAD_TLP_ERROR_COUNT_L2},
	{"ixlink_bad_tlp_error_count_l3", DCGM_FI_IXLINK_BAD_TLP_ERROR_COUNT_L3},
	{"ixlink_bad_tlp_error_count_l4", DCGM_FI_IXLINK_BAD_TLP_ERROR_COUNT_L4},
	{"ixlink_bad_tlp_error_count_l5", DCGM_FI_IXLINK_BAD_TLP_ERROR_COUNT_L5},
	{"ixlink_lcrc_error_count_l0", DCGM_FI_IXLINK_LCRC_ERROR_COUNT_L0},
	{"ixlink_lcrc_error_count_l1", DCGM_FI_IXLINK_LCRC_ERROR_COUNT_L1},
	{"ixlink_lcrc_error_count_l2", DCGM_FI_IXLINK_LCRC_ERROR_COUNT_L2},
	{"ixlink_lcrc_error_count_l3", DCGM_FI_IXLINK_LCRC_ERROR_COUNT_L3},
	{"ixlink_lcrc_error_count_l4", DCGM_FI_IXLINK_LCRC_ERROR_COUNT_L4},
	{"ixlink_lcrc_error_count_l5", DCGM_FI_IXLINK_LCRC_ERROR_COUNT_L5},
	{"ixlink_bad_dllp_error_count_l0", DCGM_FI_IXLINK_BAD_DLLP_ERROR_COUNT_L0},
	{"ixlink_bad_dllp_error_count_l1", DCGM_FI_IXLINK_BAD_DLLP_ERROR_COUNT_L1},
	{"ixlink_bad_dllp_error_count_l2", DCGM_FI_IXLINK_BAD_DLLP_ERROR_COUNT_L2},
	{"ixlink_bad_dllp_error_count_l3", DCGM_FI_IXLINK_BAD_DLLP_ERROR_COUNT_L3},
	{"ixlink_bad_dllp_error_count_l4", DCGM_FI_IXLINK_BAD_DLLP_ERROR_COUNT_L4},
	{"ixlink_bad_dllp_error_count_l5", DCGM_FI_IXLINK_BAD_DLLP_ERROR_COUNT_L5},
	{"ixlink_replay_timeout_error_count_l0", DCGM_FI_IXLINK_REPLAY_TIMEOUT_ERROR_COUNT_L0},
	{"ixlink_replay_timeout_error_count_l1", DCGM_FI_IXLINK_REPLAY_TIMEOUT_ERROR_COUNT_L1},
	{"ixlink_replay_timeout_error_count_l2", DCGM_FI_IXLINK_REPLAY_TIMEOUT_ERROR_COUNT_L2},
	{"ixlink_replay_timeout_error_count_l3", DCGM_FI_IXLINK_REPLAY_TIMEOUT_ERROR_COUNT_L3},
	{"ixlink_replay_timeout_error_count_l4", DCGM_FI_IXLINK_REPLAY_TIMEOUT_ERROR_COUNT_L4},
	{"ixlink_replay_timeout_error_count_l5", DCGM_FI_IXLINK_REPLAY_TIMEOUT_ERROR_COUNT_L5},
	{"ixlink_retry_tlp_error_count_l0", DCGM_FI_IXLINK_RETRY_TLP_ERROR_COUNT_L0},
	{"ixlink_retry_tlp_error_count_l1", DCGM_FI_IXLINK_RETRY_TLP_ERROR_COUNT_L1},
	{"ixlink_retry_tlp_error_count_l2", DCGM_FI_IXLINK_RETRY_TLP_ERROR_COUNT_L2},
	{"ixlink_retry_tlp_error_count_l3", DCGM_FI_IXLINK_RETRY_TLP_ERROR_COUNT_L3},
	{"ixlink_retry_tlp_error_count_l4", DCGM_FI_IXLINK_RETRY_TLP_ERROR_COUNT_L4},
	{"ixlink_retry_tlp_error_count_l5", DCGM_FI_IXLINK_RETRY_TLP_ERROR_COUNT_L5},
	{"ixlink_rx_recovery_error_count_l0", DCGM_FI_IXLINK_RX_RECOVERY_ERROR_COUNT_L0},
	{"ixlink_rx_recovery_error_count_l1", DCGM_FI_IXLINK_RX_RECOVERY_ERROR_COUNT_L1},
	{"ixlink_rx_recovery_error_count_l2", DCGM_FI_IXLINK_RX_RECOVERY_ERROR_COUNT_L2},
	{"ixlink_rx_recovery_error_count_l3", DCGM_FI_IXLINK_RX_RECOVERY_ERROR_COUNT_L3},
	{"ixlink_rx_recovery_error_count_l4", DCGM_FI_IXLINK_RX_RECOVERY_ERROR_COUNT_L4},
	{"ixlink_rx_recovery_error_count_l5", DCGM_FI_IXLINK_RX_RECOVERY_ERROR_COUNT_L5},
	{"ixlink_ecrc_error_count_l0", DCGM_FI_IXLINK_ECRC_ERROR_COUNT_L0},
	{"ixlink_ecrc_error_count_l1", DCGM_FI_IXLINK_ECRC_ERROR_COUNT_L1},
	{"ixlink_ecrc_error_count_l2", DCGM_FI_IXLINK_ECRC_ERROR_COUNT_L2},
	{"ixlink_ecrc_error_count_l3", DCGM_FI_IXLINK_ECRC_ERROR_COUNT_L3},
	{"ixlink_ecrc_error_count_l4", DCGM_FI_IXLINK_ECRC_ERROR_COUNT_L4},
	{"ixlink_ecrc_error_count_l5", DCGM_FI_IXLINK_ECRC_ERROR_COUNT_L5},
	{"ixlink_completion_timeout_error_count_l0", DCGM_FI_IXLINK_COMPLETION_TIMEOUT_ERROR_COUNT_L0},
	{"ixlink_completion_timeout_error_count_l1", DCGM_FI_IXLINK_COMPLETION_TIMEOUT_ERROR_COUNT_L1},
	{"ixlink_completion_timeout_error_count_l2", DCGM_FI_IXLINK_COMPLETION_TIMEOUT_ERROR_COUNT_L2},
	{"ixlink_completion_timeout_error_count_l3", DCGM_FI_IXLINK_COMPLETION_TIMEOUT_ERROR_COUNT_L3},
	{"ixlink_completion_timeout_error_count_l4", DCGM_FI_IXLINK_COMPLETION_TIMEOUT_ERROR_COUNT_L4},
	{"ixlink_completion_timeout_error_count_l5", DCGM_FI_IXLINK_COMPLETION_TIMEOUT_ERROR_COUNT_L5},
}

var (
	ixlinkSupportOnce   sync.Once
	ixlinkSupportCached bool
)

type IxlinkErrorGather struct {
	gpuIds    []uint
	fields    []Short
	fieldGrp  FieldGrpHandle
	gpuGrpHdl GroupHandle
	destroyed bool
}

func newIxlinkErrorGather(gpuId uint) (*IxlinkErrorGather, bool, error) {
	return newIxlinkErrorGatherForGPUs([]uint{gpuId})
}

func newIxlinkErrorGatherForGPUs(gpuIds []uint) (*IxlinkErrorGather, bool, error) {
	if len(gpuIds) == 0 {
		return nil, false, fmt.Errorf("gpuIds must not be empty")
	}

	if !deviceIxlinkSupported() {
		return nil, false, nil
	}

	fields := make([]Short, len(ixlinkErrorMetricDefs))
	for i, def := range ixlinkErrorMetricDefs {
		fields[i] = def.FieldID
	}

	fieldGrpName := fmt.Sprintf("devIxlinkErrorStatusFields%d", rand.Uint64())
	fieldGrp, err := FieldGroupCreate(fieldGrpName, fields)
	if err != nil {
		return nil, true, err
	}

	gpuGrpName := fmt.Sprintf("devIxlinkErrorStatusGrp%d", rand.Uint64())
	gpuGrpHdl, err := WatchFields(gpuIds, fieldGrp, gpuGrpName)
	if err != nil {
		_ = FieldGroupDestroy(fieldGrp)
		return nil, true, err
	}

	return &IxlinkErrorGather{
		gpuIds:    append([]uint(nil), gpuIds...),
		fields:    fields,
		fieldGrp:  fieldGrp,
		gpuGrpHdl: gpuGrpHdl,
	}, true, nil
}

func (gather *IxlinkErrorGather) GetIxlinkErrorStatus(gpuId uint) (IxlinkErrorStatus, error) {
	status := defaultIxlinkErrorStatus(gpuId)
	values, err := GetLatestValuesForFields(gpuId, gather.fields)
	if err != nil {
		return status, err
	}

	status = parseIxlinkErrorStatus(gpuId, values)
	return status, nil
}

func (gather *IxlinkErrorGather) GetAllIxlinkErrorStatus() (map[uint]IxlinkErrorStatus, error) {
	statuses := make(map[uint]IxlinkErrorStatus, len(gather.gpuIds))

	for _, gpuId := range gather.gpuIds {
		status, err := gather.GetIxlinkErrorStatus(gpuId)
		if err != nil {
			return nil, err
		}
		statuses[gpuId] = status
	}
	return statuses, nil
}

func (gather *IxlinkErrorGather) Destroy() error {
	if gather == nil || gather.destroyed {
		return nil
	}

	var fieldGrpErr error
	var gpuGrpErr error

	fieldGrpErr = FieldGroupDestroy(gather.fieldGrp)
	gpuGrpErr = DestroyGroup(gather.gpuGrpHdl)
	gather.destroyed = true

	if fieldGrpErr != nil && gpuGrpErr != nil {
		return fmt.Errorf("failed to destroy ixlink gather resources: field group: %v, gpu group: %v", fieldGrpErr, gpuGrpErr)
	}
	if fieldGrpErr != nil {
		return fieldGrpErr
	}
	return gpuGrpErr
}

func defaultIxlinkErrorStatus(gpuId uint) IxlinkErrorStatus {
	values := make(map[string]string, len(ixlinkErrorMetricDefs))
	for _, def := range ixlinkErrorMetricDefs {
		values[def.Name] = "N/A"
	}
	return IxlinkErrorStatus{
		Id:     gpuId,
		Values: values,
	}
}

func parseIxlinkErrorStatus(gpuId uint, values []FieldValue_v1) IxlinkErrorStatus {
	ixlink := defaultIxlinkErrorStatus(gpuId)
	if len(values) < len(ixlinkErrorMetricDefs) {
		return ixlink
	}

	for i, def := range ixlinkErrorMetricDefs {
		ixlink.Values[def.Name] = GetFieldValueStr(values[i], "int64")
	}

	return ixlink
}

func getDeviceIxlinkErrorStatus(gpuId uint) (status IxlinkErrorStatus, supported bool, err error) {
	status = defaultIxlinkErrorStatus(gpuId)
	gather, supported, err := newIxlinkErrorGather(gpuId)
	if err != nil || !supported {
		return status, supported, err
	}
	defer func() {
		_ = gather.Destroy()
	}()

	status, err = gather.GetIxlinkErrorStatus(gpuId)
	if err != nil {
		return defaultIxlinkErrorStatus(gpuId), true, err
	}
	return status, true, nil
}

type IxlinkLinkState uint

const (
	IxlinkLinkStateNotSupported IxlinkLinkState = iota // Link is unsupported by this GPU.
	IxlinkLinkStateDisabled                            // Link is supported by this GPU but disabled.
	IxlinkLinkStateDown                                // Link is down (inactive).
	IxlinkLinkStateUp                                  // Link is up (active).
)

// ixlinkFixedPortsPerGPU represents the number of links per GPU.
// This matches the currently supported IXLink topology on deployed GPUs.
// Note: IXDCGM_MAX_LINKS_PER_GPU can be 18 in headers, but l6-l17 link states are
// not exposed in the current IXLink link status metric set yet.
const ixlinkFixedPortsPerGPU = 6

type IxlinkGpuLinkStatus struct {
	EntityID  uint
	LinkState [ixlinkFixedPortsPerGPU]IxlinkLinkState
}

type IxlinkSwitchStatus struct {
	EntityID  uint
	LinkState []IxlinkLinkState
}

type IxlinkStatus struct {
	GPUs       []IxlinkGpuLinkStatus
	IxSwitches []IxlinkSwitchStatus
}

func (status IxlinkStatus) SupportedIxLink() bool {
	for _, gpu := range status.GPUs {
		for _, linkState := range gpu.LinkState {
			if linkState != IxlinkLinkStateNotSupported {
				return true
			}
		}
	}
	return false
}

func deviceIxlinkSupported() bool {
	ixlinkSupportOnce.Do(func() {
		linkStatus, err := getIxlinkStatus()
		if err != nil {
			ixlinkSupportCached = false
			return
		}
		ixlinkSupportCached = linkStatus.SupportedIxLink()
	})
	return ixlinkSupportCached
}

func getIxlinkStatus() (IxlinkStatus, error) {
	var linkStatus C.ixdcgmLinkStatus_v3
	linkStatus.version = C.uint(makeVersion3(unsafe.Sizeof(linkStatus)))

	ret := C.ixdcgmGetLinkStatus(C.ulong(handle.handle), &linkStatus)
	if err := ixdcgmErrorString(ret); err != nil {
		return IxlinkStatus{}, fmt.Errorf("error getting ixlink link status: %w", err)
	}
	status := IxlinkStatus{}

	for i := 0; i < int(linkStatus.numGpus); i++ {
		cGpu := linkStatus.gpus[i]
		gpu := IxlinkGpuLinkStatus{
			EntityID: uint(cGpu.entityId),
		}
		for linkID := 0; linkID < ixlinkFixedPortsPerGPU; linkID++ {
			gpu.LinkState[linkID] = IxlinkLinkState(cGpu.linkState[linkID])
		}
		status.GPUs = append(status.GPUs, gpu)
	}

	// linkStatus.numNvSwitches is the number of IxSwitches
	for i := 0; i < int(linkStatus.numNvSwitches); i++ {
		// linkStatus.nvSwitches is the array of IxSwitches
		cSwitch := linkStatus.nvSwitches[i]
		ixSwitch := IxlinkSwitchStatus{
			EntityID: uint(cSwitch.entityId),
		}
		for _, linkState := range cSwitch.linkState {
			state := IxlinkLinkState(linkState)
			if state != IxlinkLinkStateNotSupported {
				ixSwitch.LinkState = append(ixSwitch.LinkState, state)
			}
		}
		status.IxSwitches = append(status.IxSwitches, ixSwitch)
	}

	return status, nil
}

func getDeviceIxlinkGpuStatus(gpuId uint) (status IxlinkGpuLinkStatus, found bool, err error) {
	linkStatus, err := getIxlinkStatus()
	if err != nil {
		return IxlinkGpuLinkStatus{}, false, err
	}

	for _, gpuLinkStatus := range linkStatus.GPUs {
		if gpuLinkStatus.EntityID == gpuId {
			return gpuLinkStatus, true, nil
		}
	}
	return IxlinkGpuLinkStatus{}, false, nil
}
