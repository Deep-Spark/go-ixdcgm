/* Copyright © 2024 Iluvatar CoreX. All rights reserved.                                                            */
/* Copyright Declaration: This software, including all of its code and documentation, except for the third-party    */
/* software it contains, is a copyrighted work of Shanghai Iluvatar CoreX Semiconductor Co., Ltd. and its           */
/* affiliates (“Iluvatar CoreX”) in accordance with the PRC Copyright Law and relevant international treaties,      */
/* and all rights contained therein are enjoyed by Iluvatar CoreX. No user of this software shall have any right,   */
/* ownership or interest in this software and any use of this software shall be in compliance with the terms        */
/* and conditions of the End User License Agreement.                                                                */

#ifndef __IXDCGM_API_EXPORT_H__
#define __IXDCGM_API_EXPORT_H__

#include "ixdcgmStructs.h"

#include "ixdcgmFields.h"

#include <stdbool.h>

#ifdef __cplusplus
extern "C"
{
#endif

#if defined(IXDCGM_API_EXPORT)
#define IXDCGM_PUBLIC_API __attribute((visibility("default")))
#else
#define IXDCGM_PUBLIC_API
#endif

#define IXDCGM_PRIVATE_API __attribute((visibility("hidden")))
    /**
     * Get the IxLink link status for every IxLink in this system.
     *
     * @param pixdcgmHandle  IN: IxDCGM Handle
     * @param linkStatus     OUT: Structure in which to store IxLink link statuses.
     *
     * @return
     *        - \ref IXDCGM_RET_OK                if the call was successful.
     *        - \ref IXDCGM_RET_NOT_SUPPORTED     if the given entityGroup does not support enumeration.
     *        - \ref IXDCGM_RET_BADPARAM          if any parameter is invalid
     *        - \ref IXDCGM_RET_VER_MISMATCH      if the version of linkStatus is not ixdcgmLinkStatus_v3
     */
    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGetLinkStatus(ixdcgmHandle_t pixdcgmHandle, ixdcgmLinkStatus_v3 *linkStatus);

    /**
     * Gets the 2 GPUs are on the same board or not.
     * @param pixdcgmHandle    IN: IxDCGM Handle
     * @param gpuId1           IN: GPU1 Id
     * @param gpuId2           IN: GPU2 Id
     * @param onSameBoard IN/OUT: On same board info of the GPU pair.   0= not on the same board; 1= on the same board
     *
     * @return
     *        - \ref IXDCGM_RET_OK                   if the call was successful.
     *        - \ref IXDCGM_RET_BADPARAM             if gpuId1, gpuId2 or onSameBoard were not valid.
     */
    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmDeviceOnSameBoard(ixdcgmHandle_t pixdcgmHandle,
                                                             unsigned int gpuId1,
                                                             unsigned int gpuId2,
                                                             int *onSameBoard);

    /**
     * Gets all the running process info corresponding to the gpuId .
     * @param pixdcgmHandle         IN: IxDCGM Handle
     * @param gpuId                 IN: GPU Id corresponding to which the processes info should be fetched
     * @param infoCount             IN/OUT:
     *  IN - max number of the info could be stored in to the pids and usedMemoryBytes buffer
     *  OUT - When API return DCGM_ST_OK, stored number of valid pids/usedMemoryBytes info collected.
     *    When API return IXDCGM_RET_INSUFFICIENT_SIZE, stored the number of buffer needed.
     * @param pids                  OUT: Buffer to store returned processes pid
     * @param usedMemoryBytes       OUT: Buffer to store returned processes used memory in byte
     *
     * @return
     * - \ref IXDCGM_RET_OK                      if the call was successful.
     * - \ref IXDCGM_RET_INSUFFICIENT_SIZE       if the infoCount input is smaller than the buffer needed.
     * - \ref IXDCGM_RET_BADPARAM                if gpuId, infoCunt, pids or usedMemoryBytes not valid.
     **/
    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGetDeviceRunningProcesses(ixdcgmHandle_t pixdcgmHandle,
                                                                     unsigned int gpuId,
                                                                     unsigned int *infoCount,
                                                                     uint64_t *pids,
                                                                     uint64_t *usedMemoryBytes);

    IXDCGM_PUBLIC_API const char *ixdcgmErrorString(ixdcgmReturn_t result);

#ifdef __cplusplus
}
#endif

#endif // end of __IXDCGM_API_EXPORT_H__