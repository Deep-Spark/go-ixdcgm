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

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmInit(void);
    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmStartEmbedded(ixdcgmStartEmbeddedParam *params);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmEngineStart(unsigned short portNum, char const *socketPath, bool overTCP);
    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmEngineRun(unsigned short portNumber,
                                                     char const *socketPath,
                                                     unsigned int isConnectionTCP);
    IXDCGM_PUBLIC_API const char *ixdcgmErrorString(ixdcgmReturn_t result);
    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmDisconnect(ixdcgmHandle_t pixdcgmHandle);
    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmConnect(const char *ipAddress,
                                                   ixdcgmConnectParams *connectParams,
                                                   ixdcgmHandle_t *pixdcgmHandle);
    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGetEntityGroupEntities(ixdcgmHandle_t pixdcgmHandle,
                                                                  ixdcgm_field_entity_group_t entityGroup,
                                                                  ixdcgm_field_eid_t *entities,
                                                                  int *numEntities,
                                                                  unsigned int flags);
    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGetDeviceAttributes(ixdcgmHandle_t pixdcgmHandle,
                                                               unsigned int gpuId,
                                                               ixdcgmDeviceAttributes_t *pixdcgmAttr);
    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGetAllDevices(ixdcgmHandle_t pixdcgmHandle,
                                                         unsigned int gpuIdList[IXDCGM_MAX_NUM_DEVICES],
                                                         int *count);
    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGetAllSupportedDevices(ixdcgmHandle_t pixdcgmHandle,
                                                                  unsigned int gpuIdList[IXDCGM_MAX_NUM_DEVICES],
                                                                  int *count);
    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmEntitiesGetLatestValues(ixdcgmHandle_t pDcgmHandle,
                                                                   ixdcgmGroupEntityPair_t entities[],
                                                                   unsigned int entityCount,
                                                                   unsigned short fields[],
                                                                   unsigned int fieldCount,
                                                                   unsigned int flags,
                                                                   ixdcgmFieldValue_v2 values[]);
    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmHostengineVersionInfo(ixdcgmHandle_t pixdcgmHandle,
                                                                 ixdcgmVersionInfo_t *pVersionInfo);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmVersionInfo(ixdcgmVersionInfo_t *pVersionInfo);

    /*Grouping APIs*/
    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGroupCreate(ixdcgmHandle_t pixdcgmHandle,
                                                       ixdcgmGroupType_t type,
                                                       const char *groupName,
                                                       ixdcgmGpuGrp_t *groupId);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGroupDestroy(ixdcgmHandle_t pixdcgmHandle, ixdcgmGpuGrp_t groupId);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGroupAddEntity(ixdcgmHandle_t pixdcgmHandle,
                                                          ixdcgmGpuGrp_t groupId,
                                                          ixdcgm_field_entity_group_t entityGroupId,
                                                          ixdcgm_field_eid_t entityId);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGroupAddDevice(ixdcgmHandle_t pixdcgmHandle,
                                                          ixdcgmGpuGrp_t groupId,
                                                          unsigned int gpuId);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGroupRemoveDevice(ixdcgmHandle_t pixdcgmHandle,
                                                             ixdcgmGpuGrp_t groupId,
                                                             unsigned int gpuId);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGroupRemoveEntity(ixdcgmHandle_t pixdcgmHandle,
                                                             ixdcgmGpuGrp_t groupId,
                                                             ixdcgm_field_entity_group_t entityGroupId,
                                                             ixdcgm_field_eid_t entityId);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGroupGetInfo(ixdcgmHandle_t pixdcgmHandle,
                                                        ixdcgmGpuGrp_t groupId,
                                                        ixdcgmGroupInfo_t *pDcgmGroupInfo);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGroupGetAllIds(ixdcgmHandle_t pixdcgmHandle,
                                                          ixdcgmGpuGrp_t groupIdList[],
                                                          unsigned int *count);

    /* Field Grouping APIs*/
    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmFieldGroupCreate(ixdcgmHandle_t pixdcgmHandle,
                                                            int numFieldIds,
                                                            unsigned short *fieldIds,
                                                            const char *fieldGroupName,
                                                            ixdcgmFieldGrp_t *fieldGroupId);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmFieldGroupDestroy(ixdcgmHandle_t pixdcgmHandle,
                                                             ixdcgmFieldGrp_t fieldGroupId);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmFieldGroupGetInfo(ixdcgmHandle_t pixdcgmHandle,
                                                             ixdcgmFieldGroupInfo_t *fieldGroupInfo);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmFieldGroupGetAll(ixdcgmHandle_t pixdcgmHandle,
                                                            ixdcgmAllFieldGroup_t *allGroupInfo);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmWatchFields(ixdcgmHandle_t pixdcgmHandle,
                                                       ixdcgmGpuGrp_t groupId,
                                                       ixdcgmFieldGrp_t fieldGroupId,
                                                       long long updateFreq,
                                                       double maxKeepAge,
                                                       int maxKeepSamples);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmUnwatchFields(ixdcgmHandle_t pixdcgmHandle,
                                                         ixdcgmGpuGrp_t groupId,
                                                         ixdcgmFieldGrp_t fieldGroupId);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmStatusCreate(ixdcgmStatus_t *statusHandle);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmStatusDestroy(ixdcgmStatus_t statusHandle);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmStatusGetCount(ixdcgmStatus_t statusHandle, unsigned int *count);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmStatusPopError(ixdcgmStatus_t statusHandle,
                                                          ixdcgmErrorInfo_t *pixdcgmErrorInfo);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmStatusClear(ixdcgmStatus_t statusHandle);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmConfigGet(ixdcgmHandle_t pixdcgmHandle,
                                                     ixdcgmGpuGrp_t groupId,
                                                     ixdcgmConfigType_t type,
                                                     int count,
                                                     ixdcgmConfig_t deviceConfigList[],
                                                     ixdcgmStatus_t statusHandle);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmConfigSet(ixdcgmHandle_t pixdcgmHandle,
                                                     ixdcgmGpuGrp_t groupId,
                                                     ixdcgmConfig_t *pDeviceConfig,
                                                     ixdcgmStatus_t statusHandle);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmConfigEnforce(ixdcgmHandle_t pixdcgmHandle,
                                                         ixdcgmGpuGrp_t groupId,
                                                         ixdcgmStatus_t statusHandle);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGetValuesSince_v2(ixdcgmHandle_t pixdcgmHandle,
                                                             ixdcgmGpuGrp_t groupId,
                                                             ixdcgmFieldGrp_t fieldGroupId,
                                                             long long sinceTimestamp,
                                                             long long *nextSinceTimestamp,
                                                             ixdcgmFieldValueEntityEnumeration_f enumCB,
                                                             void *userData);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGetLatestValues_v2(ixdcgmHandle_t pDcgmHandle,
                                                              ixdcgmGpuGrp_t groupId,
                                                              ixdcgmFieldGrp_t fieldGroupId,
                                                              ixdcgmFieldValueEntityEnumeration_f enumCB,
                                                              void *userData);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmStopEmbedded(ixdcgmHandle_t pixdcgmHandle);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGetFieldSummary(ixdcgmHandle_t pixdcgmHandle,
                                                           ixdcgmFieldSummaryRequest_t *request);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmShutdown(void);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmModuleIdToName(ixdcgmModuleId_t id, char const **name);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGetLatestValuesForFields(ixdcgmHandle_t pixdcgmHandle,
                                                                    int gpuId,
                                                                    unsigned short fields[],
                                                                    unsigned int count,
                                                                    ixdcgmFieldValue_v1 values[]);
    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmUpdateAllFields(ixdcgmHandle_t pixdcgmHandle, int waitForUpdate);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmHostengineSetLoggingSeverity(ixdcgmHandle_t pixdcgmHandle,
                                                                        ixdcgmSettingsSetLoggingSeverity_t *logging);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmDeviceOnSameBoard(ixdcgmHandle_t pixdcgmHandle,
                                                             unsigned int gpuId1,
                                                             unsigned int gpuId2,
                                                             int *onSameBoard);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmModuleGetStatuses(ixdcgmHandle_t pixdcgmHandle,
                                                             ixdcgmModuleGetStatuses_t *moduleStatuses);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmModuleDenylist(ixdcgmHandle_t pixdcgmHandle, ixdcgmModuleId_t moduleId);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGetDeviceTopology(ixdcgmHandle_t pixdcgmHandle,
                                                             unsigned int gpuId,
                                                             ixdcgmDeviceTopology_t *deviceTopology);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGetGroupTopology(ixdcgmHandle_t pixdcgmHandle,
                                                            ixdcgmGpuGrp_t groupId,
                                                            ixdcgmGroupTopology_t *pixdcgmGroupTopology);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmSelectGpusByTopology(ixdcgmHandle_t pixdcgmHandle,
                                                                uint64_t inputGpuIds,
                                                                uint32_t numGpus,
                                                                uint64_t *outputGpuIds,
                                                                uint64_t hintFlags);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmHealthGet(ixdcgmHandle_t pixdcgmHandle,
                                                     ixdcgmGpuGrp_t groupId,
                                                     ixdcgmHealthSystems_t *systems);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmHealthSet_v2(ixdcgmHandle_t pixdcgmHandle, ixdcgmHealthSetParams_v2 *params);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmHealthCheck_v4(ixdcgmHandle_t pixdcgmHandle,
                                                          ixdcgmGpuGrp_t groupId,
                                                          ixdcgmHealthResponse_v4 *response);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGetLinkStatus(ixdcgmHandle_t pixdcgmHandle, ixdcgmLinkStatus_v3 *linkStatus);

    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmProfGetSupportedMetricGroups(ixdcgmHandle_t pixdcgmHandle,
                                                                        ixdcgmProfGetMetricGroups_t *metricGroups);

    /**
     * Gets all the running process info corresponding to the gpuId .
     * @param pixdcgmHandle         IN: ixDCGM Handle
     * @param gpuId                 IN: GPU Id corresponding to which the processes info should be fetched
     * @param infoCount             IN/OUT:
     *  IN - max number of the info could be stored in to the pids and usedMemoryBytes buffer
     *  OUT - When API return DCGM_ST_OK, stored number of valid pids/usedMemoryBytes info collected.
     *    When API return IXDCGM_RET_INSUFFICIENT_SIZE, stored the number of buffer needed.
     * @param pids                  OUT: Buffer to store returned processes pid
     * @param usedMemoryBytes       OUT: Buffer to store returned processes used memory in byte
     *
     * @return
     * - \ref DCGM_ST_OK                         if the call was successful.
     * - \ref IXDCGM_RET_INSUFFICIENT_SIZE       if the infoCount input is smaller than the buffer needed.
     * - \ref DCGM_ST_BADPARAM                   if gpuId, infoCunt, pids or usedMemoryBytes not valid.
     **/
    ixdcgmReturn_t IXDCGM_PUBLIC_API ixdcgmGetDeviceRunningProcesses(ixdcgmHandle_t pixdcgmHandle,
                                                                     unsigned int gpuId,
                                                                     unsigned int *infoCount,
                                                                     uint64_t *pids,
                                                                     uint64_t *usedMemoryBytes);

#ifdef __cplusplus
}
#endif

#endif // end of __IXDCGM_API_EXPORT_H__