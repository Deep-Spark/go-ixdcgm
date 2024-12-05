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

import "C"
import "unsafe"

type Short C.ushort

type FieldValue_v1 struct {
	Version   uint
	FieldId   uint
	FieldType uint
	Status    int
	Ts        int64
	Value     [4096]byte
}

func (fv FieldValue_v1) Int64() int64 {
	return *(*int64)(unsafe.Pointer(&fv.Value[0]))
}

func (fv FieldValue_v1) Float64() float64 {
	return *(*float64)(unsafe.Pointer(&fv.Value[0]))
}
