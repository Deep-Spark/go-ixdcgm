/*
Copyright (c) 2024, NVIDIA CORPORATION.
Copyright (c) 2024, Shanghai Iluvatar CoreX Semiconductor Co., Ltd.
All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

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

// wrapper for go callback function
extern int violationNotify(void* p);
extern int voidCallback(void* p);
*/
import "C"
import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
	"unsafe"

	"github.com/creasty/defaults"
)

// At least one policy must be enabled.
type PolicyConditionParams struct {
	// DbePolicyEnabled indicates whether the DbePolicy is enabled. Default is false (disabled).
	DbePolicyEnabled bool `default:"false"`

	// PCIePolicyEnabled indicates whether the PCIePolicy is enabled. Default is false (disabled).
	PCIePolicyEnabled bool `default:"false"`

	// MaxRtPgPolicyEnabled indicates whether the MaxRtPgPolicy is enabled. Default is false (disabled).
	MaxRtPgPolicyEnabled bool `default:"false"`

	// MaxRtPgPolicyThreshold specifies the maximum number of retired pages that will trigger a violation.
	// Note that the MaxRtPgPolicyThreshold will be ignored if MaxRtPgPolicy is disabled.
	// Default value is 10.
	MaxRtPgPolicyThreshold uint32 `default:"10"`

	// ThermalPolicyEnabled indicates whether the ThermalPolicy is enabled. Default is false (disabled).
	ThermalPolicyEnabled bool `default:"false"`

	// ThermalPolicyThreshold specifies the maximum temperature a group's GPUs can reach before triggering a violation.
	// Note that the ThermalPolicyThreshold will be ignored if ThermalPolicy is disabled.
	// Default value is 100 and the unit is in degrees Celsius (°C).
	ThermalPolicyThreshold uint32 `default:"100"`

	// PowerPolicyEnabled indicates whether the PowerPolicy is enabled. Default is false (disabled).
	PowerPolicyEnabled bool `default:"false"`

	// PowerPolicyThreshold specifies the maximum power a group's GPUs can reach before triggering a violation.
	// Note that the PowerPolicyThreshold will be ignored if PowerPolicy is fadisabledlse.
	// Default value is 250 and the unit is in watts (W).
	PowerPolicyThreshold uint32 `default:"250"`
}

type policyCondition string

const (
	DbePolicy     = policyCondition("Double-bit ECC Error")
	PCIePolicy    = policyCondition("PCI Error")
	MaxRtPgPolicy = policyCondition("Max Retired Pages Limit")
	ThermalPolicy = policyCondition("Thermal Limit")
	PowerPolicy   = policyCondition("Power Limit")
)

type PolicyViolation struct {
	Condition policyCondition
	Timestamp time.Time
	Data      interface{}
}

type policyIndex int

const (
	dbePolicyIndex policyIndex = iota
	pciePolicyIndex
	maxRtPgPolicyIndex
	thermalPolicyIndex
	powerPolicyIndex
)

type policyConditionParam struct {
	typ   uint32
	value uint32
}

type DbePolicyCondition struct {
	Location  string
	NumErrors uint
}

type PciPolicyCondition struct {
	ReplayCounter uint
}

type RetiredPagesPolicyCondition struct {
	SbePages uint
	DbePages uint
}

type ThermalPolicyCondition struct {
	ThermalViolation uint
}

type PowerPolicyCondition struct {
	PowerViolation uint
}

var (
	policyChanOnce sync.Once
	policyMapOnce  sync.Once

	// callbacks maps PolicyViolation channels with policy
	// captures C callback() value for each violation condition
	callbacks map[string]chan PolicyViolation

	// paramMap maps C.dcgmPolicy_t.parms index and limits
	// to be used in setPolicy() for setting user selected policies
	paramMap map[policyIndex]policyConditionParam

	registerCh = make(chan struct{})
)

func makePolicyChannels() {
	policyChanOnce.Do(func() {
		callbacks = make(map[string]chan PolicyViolation)
		callbacks["dbe"] = make(chan PolicyViolation, 1)
		callbacks["pcie"] = make(chan PolicyViolation, 1)
		callbacks["maxrtpg"] = make(chan PolicyViolation, 1)
		callbacks["thermal"] = make(chan PolicyViolation, 1)
		callbacks["power"] = make(chan PolicyViolation, 1)
	})
}

func makePolicyParamsMap(params *PolicyConditionParams) {
	const (
		policyFieldTypeBool = 0
		policyFieldTypeLong = 1
		policyBoolValue     = 1
	)

	policyMapOnce.Do(func() {
		paramMap = make(map[policyIndex]policyConditionParam)

		paramMap[dbePolicyIndex] = policyConditionParam{
			typ:   policyFieldTypeBool,
			value: policyBoolValue,
		}

		paramMap[pciePolicyIndex] = policyConditionParam{
			typ:   policyFieldTypeBool,
			value: policyBoolValue,
		}

		paramMap[maxRtPgPolicyIndex] = policyConditionParam{
			typ:   policyFieldTypeLong,
			value: params.MaxRtPgPolicyThreshold,
		}

		paramMap[thermalPolicyIndex] = policyConditionParam{
			typ:   policyFieldTypeLong,
			value: params.ThermalPolicyThreshold,
		}

		paramMap[powerPolicyIndex] = policyConditionParam{
			typ:   policyFieldTypeLong,
			value: params.PowerPolicyThreshold,
		}

	})
}

func setPolicy(groupId GroupHandle, condition C.dcgmPolicyCondition_t, paramList []policyIndex) (err error) {
	var policy C.dcgmPolicy_t
	policy.version = makeVersion1(unsafe.Sizeof(policy))
	policy.mode = C.dcgmPolicyMode_t(C.DCGM_OPERATION_MODE_AUTO)
	policy.action = C.DCGM_POLICY_ACTION_NONE
	policy.isolation = C.DCGM_POLICY_ISOLATION_NONE
	policy.validation = C.DCGM_POLICY_VALID_NONE
	policy.condition = condition

	// iterate on paramMap for given policy conditions
	for _, key := range paramList {
		conditionParam, exists := paramMap[policyIndex(key)]
		if !exists {
			return fmt.Errorf("Error: Invalid Policy condition, %v does not exist", key)
		}
		// set policy condition parameters
		// set condition type (bool or longlong)
		policy.parms[key].tag = conditionParam.typ

		// set condition val (violation threshold)
		// policy.parms.val is a C union type
		// cgo docs: Go doesn't have support for C's union type
		// C union types are represented as a Go byte array
		binary.LittleEndian.PutUint32(policy.parms[key].val[:], conditionParam.value)
	}

	var statusHandle C.dcgmStatus_t

	result := C.dcgmPolicySet(handle.handle, groupId.handle, &policy, statusHandle)
	if err = errorString(result); err != nil {
		return fmt.Errorf("Error setting policies: %s", err)
	}

	log.Println("Policy successfully set.")

	return
}

func validatePolicy(p *PolicyConditionParams) error {
	if err := defaults.Set(p); err != nil {
		return err
	}
	if !(p.DbePolicyEnabled || p.PCIePolicyEnabled || p.MaxRtPgPolicyEnabled || p.ThermalPolicyEnabled || p.PowerPolicyEnabled) {
		return fmt.Errorf("bad parameters: at least one policy must be enabled")
	}
	return nil
}

func registerPolicyForGpus(ctx context.Context, params *PolicyConditionParams, gpuIds ...uint) (<-chan PolicyViolation, error) {
	groupId, err := CreateGroup(fmt.Sprintf("PolicyGroup_%d", rand.Uint64()))
	if err != nil {
		return nil, fmt.Errorf("failed to create policy group, err: %v", err)
	}

	go func() {
		<-ctx.Done()
		select {
		case <-registerCh: // Wait the policy is unregistered
			_ = DestroyGroup(groupId)
		case <-time.After(500 * time.Millisecond):
			_ = DestroyGroup(groupId)
		}
	}()

	for _, gpuId := range gpuIds {
		err = AddToGroup(groupId, gpuId)
		if err != nil {
			return nil, fmt.Errorf("failed to add gpu %d to policy group, err: %v", gpuId, err)
		}
	}

	return registerPolicy(ctx, groupId, params)
}

// registerPolicy sets GPU usage and error policies and notifies in case of any violations on GPUs within a specific group
func registerPolicy(ctx context.Context, groupId GroupHandle, params *PolicyConditionParams) (<-chan PolicyViolation, error) {
	if params == nil {
		return nil, fmt.Errorf("PolicyConditionParams is required")
	}
	if err := validatePolicy(params); err != nil {
		return nil, err
	}

	// init policy globals for internal API
	makePolicyChannels()
	makePolicyParamsMap(params)

	// make a list of policy conditions for setting their parameters
	var paramKeys []policyIndex
	// get all conditions to be set in setPolicy()
	var condition C.dcgmPolicyCondition_t = 0
	// get length of enabled condition types
	var conTypes int = 0
	if params.DbePolicyEnabled {
		conTypes++
		paramKeys = append(paramKeys, dbePolicyIndex)
		condition |= C.DCGM_POLICY_COND_DBE
	}
	if params.PCIePolicyEnabled {
		conTypes++
		paramKeys = append(paramKeys, pciePolicyIndex)
		condition |= C.DCGM_POLICY_COND_PCI
	}
	if params.MaxRtPgPolicyEnabled {
		conTypes++
		paramKeys = append(paramKeys, maxRtPgPolicyIndex)
		condition |= C.DCGM_POLICY_COND_MAX_PAGES_RETIRED
	}
	if params.ThermalPolicyEnabled {
		conTypes++
		paramKeys = append(paramKeys, thermalPolicyIndex)
		condition |= C.DCGM_POLICY_COND_THERMAL
	}
	if params.PowerPolicyEnabled {
		conTypes++
		paramKeys = append(paramKeys, powerPolicyIndex)
		condition |= C.DCGM_POLICY_COND_POWER
	}

	var err error
	if err = setPolicy(groupId, condition, paramKeys); err != nil {
		return nil, err
	}

	result := C.dcgmPolicyRegister(handle.handle, groupId.handle,
		C.dcgmPolicyCondition_t(condition),
		C.fpRecvUpdates(C.violationNotify),
		C.fpRecvUpdates(C.voidCallback),
	)

	if err = errorString(result); err != nil {
		return nil, &DcgmError{msg: C.GoString(C.errorString(result)), Code: result}
	}
	log.Println("Listening for violations...")

	violation := make(chan PolicyViolation, conTypes)

	go func() {
		defer func() {
			log.Println("unregister policy violation...")
			unregisterPolicy(groupId, condition)
			close(violation)
			close(registerCh)
		}()
		for {
			select {
			case dbe := <-callbacks["dbe"]:
				violation <- dbe
			case pcie := <-callbacks["pcie"]:
				violation <- pcie
			case maxrtpg := <-callbacks["maxrtpg"]:
				violation <- maxrtpg
			case thermal := <-callbacks["thermal"]:
				violation <- thermal
			case power := <-callbacks["power"]:
				violation <- power
			case <-ctx.Done():
				return
			}
		}
	}()

	return violation, err
}

func unregisterPolicy(groupId GroupHandle, condition C.dcgmPolicyCondition_t) {
	result := C.dcgmPolicyUnregister(handle.handle, groupId.handle, condition)

	if err := errorString(result); err != nil {
		log.Println(fmt.Errorf("error unregistering policy: %s", err))
	}
}

func createTimeStamp(t C.longlong) time.Time {
	tm := int64(t) / 1000000
	ts := time.Unix(tm, 0)
	return ts
}

func dbeLocation(location int) string {
	switch location {
	case 0:
		return "L1"
	case 1:
		return "L2"
	case 2:
		return "Device"
	case 3:
		return "Register"
	case 4:
		return "Texture"
	}
	return "N/A"
}

// VoidCallback is a go callback function for dcgmPolicyRegister() wrapped in C.voidCallback()
//
//export VoidCallback
func VoidCallback(data unsafe.Pointer) int {
	return 0
}

// ViolationRegistration is a go callback function for dcgmPolicyRegister() wrapped in C.violationNotify()
//
//export ViolationRegistration
func ViolationRegistration(data unsafe.Pointer) int {
	var con policyCondition
	var timestamp time.Time
	var val interface{}

	response := *(*C.dcgmPolicyCallbackResponse_t)(unsafe.Pointer(data))

	switch response.condition {
	case C.DCGM_POLICY_COND_DBE:
		dbe := (*C.dcgmPolicyConditionDbe_t)(unsafe.Pointer(&response.val))
		con = DbePolicy
		timestamp = createTimeStamp(dbe.timestamp)
		val = DbePolicyCondition{
			Location:  dbeLocation(int(dbe.location)),
			NumErrors: *uintPtr(dbe.numerrors),
		}
	case C.DCGM_POLICY_COND_PCI:
		pci := (*C.dcgmPolicyConditionPci_t)(unsafe.Pointer(&response.val))
		con = PCIePolicy
		timestamp = createTimeStamp(pci.timestamp)
		val = PciPolicyCondition{
			ReplayCounter: *uintPtr(pci.counter),
		}
	case C.DCGM_POLICY_COND_MAX_PAGES_RETIRED:
		mpr := (*C.dcgmPolicyConditionMpr_t)(unsafe.Pointer(&response.val))
		con = MaxRtPgPolicy
		timestamp = createTimeStamp(mpr.timestamp)
		val = RetiredPagesPolicyCondition{
			SbePages: *uintPtr(mpr.sbepages),
			DbePages: *uintPtr(mpr.dbepages),
		}
	case C.DCGM_POLICY_COND_THERMAL:
		thermal := (*C.dcgmPolicyConditionThermal_t)(unsafe.Pointer(&response.val))
		con = ThermalPolicy
		timestamp = createTimeStamp(thermal.timestamp)
		val = ThermalPolicyCondition{
			ThermalViolation: *uintPtr(thermal.thermalViolation),
		}
	case C.DCGM_POLICY_COND_POWER:
		pwr := (*C.dcgmPolicyConditionPower_t)(unsafe.Pointer(&response.val))
		con = PowerPolicy
		timestamp = createTimeStamp(pwr.timestamp)
		val = PowerPolicyCondition{
			PowerViolation: *uintPtr(pwr.powerViolation),
		}
	}

	err := PolicyViolation{
		Condition: con,
		Timestamp: timestamp,
		Data:      val,
	}

	switch con {
	case DbePolicy:
		callbacks["dbe"] <- err
	case PCIePolicy:
		callbacks["pcie"] <- err
	case MaxRtPgPolicy:
		callbacks["maxrtpg"] <- err
	case ThermalPolicy:
		callbacks["thermal"] <- err
	case PowerPolicy:
		callbacks["power"] <- err
	}
	return 0
}
