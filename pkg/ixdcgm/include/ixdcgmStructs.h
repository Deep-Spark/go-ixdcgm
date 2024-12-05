/* Copyright © 2024 Iluvatar CoreX. All rights reserved.                                                            */
/* Copyright Declaration: This software, including all of its code and documentation, except for the third-party    */
/* software it contains, is a copyrighted work of Shanghai Iluvatar CoreX Semiconductor Co., Ltd. and its           */
/* affiliates (“Iluvatar CoreX”) in accordance with the PRC Copyright Law and relevant international treaties,      */
/* and all rights contained therein are enjoyed by Iluvatar CoreX. No user of this software shall have any right,   */
/* ownership or interest in this software and any use of this software shall be in compliance with the terms        */
/* and conditions of the End User License Agreement.                                                                */

#ifndef __IXDCGM_STRUCTS_H__
#define __IXDCGM_STRUCTS_H__

#include <stdint.h>

#include "ixdcgmFields.h"

#define MAKE_IXDCGM_VERSION(typeName, ver) (unsigned int)(sizeof(typeName) | ((unsigned long)(ver) << 24U))

#define IXDCGM_MAX_STR_LENGTH 256
#define IXDCGM_MAX_BLOB_LENGTH 4096

#define IXDCGM_VGPU_NAME_BUFFER_SIZE 64
#define IXDCGM_DEVICE_UUID_BUFFER_SIZE 80

#define IXDCGM_GROUP_MAX_ENTITIES 64

#define IXDCGM_INT32_BLANK 0x7ffffff0
#define IXDCGM_INT64_BLANK 0x7ffffffffffffff0ll
#define IXDCGM_FP64_BLANK 140737488355328.0
#define IXDCGM_STR_BLANK "<<<NULL>>>"

#define IXDCGM_INT32_NOT_FOUND (IXDCGM_INT32_BLANK + 1)
#define IXDCGM_INT64_NOT_FOUND (IXDCGM_INT64_BLANK + 1)
#define IXDCGM_FP64_NOT_FOUND (IXDCGM_FP64_BLANK + 1.0)
#define IXDCGM_STR_NOT_FOUND "<<<NOT_FOUND>>>"

#define IXDCGM_INT32_NOT_SUPPORTED (IXDCGM_INT32_BLANK + 2)
#define IXDCGM_INT64_NOT_SUPPORTED (IXDCGM_INT64_BLANK + 2)
#define IXDCGM_FP64_NOT_SUPPORTED (IXDCGM_FP64_BLANK + 2.0)
#define IXDCGM_STR_NOT_SUPPORTED "<<<NOT_SUPPORTED>>>"

/* Represents and error where fetching the value is not allowed with our current credentials */
#define IXDCGM_INT32_NOT_PERMISSIONED (IXDCGM_INT32_BLANK + 3)
#define IXDCGM_INT64_NOT_PERMISSIONED (IXDCGM_INT64_BLANK + 3)
#define IXDCGM_FP64_NOT_PERMISSIONED (IXDCGM_FP64_BLANK + 3.0)
#define IXDCGM_STR_NOT_PERMISSIONED "<<<NOT_PERM>>>"
#define IXDCGM_INT32_IS_BLANK(val) (((val) >= IXDCGM_INT32_BLANK) ? 1 : 0)
#define IXDCGM_INT64_IS_BLANK(val) (((val) >= IXDCGM_INT64_BLANK) ? 1 : 0)
#define IXDCGM_FP64_IS_BLANK(val) (((val) >= IXDCGM_FP64_BLANK ? 1 : 0))
#define IXDCGM_STR_IS_BLANK(val) (val == strstr(val, "<<<") && strstr(val, ">>>"))

#define IXDCGM_AFFINITY_BITMASK_ARRAY_SIZE 8

typedef enum ixdcgmReturn_enum
{
    IXDCGM_RET_OK                   = 0,    //!< Success
    IXDCGM_RET_BADPARAM             = -1,   //!< A bad parameter was passed to a function
    IXDCGM_RET_GENERIC_ERROR        = -3,   //!< A generic, unspecified error
    IXDCGM_RET_MEMORY               = -4,   //!< An out of memory error occurred
    IXDCGM_RET_NOT_CONFIGURED       = -5,   //!< Setting not configured
    IXDCGM_RET_NOT_SUPPORTED        = -6,   //!< Feature not supported
    IXDCGM_RET_INIT_ERROR           = -7,   //!< DCGM Init error
    IXDCGM_RET_NVML_ERROR           = -8,   //!< When NVML returns error
    IXDCGM_RET_PENDING              = -9,   //!< Object is in pending state of something else
    IXDCGM_RET_UNINITIALIZED        = -10,  //!< Object is in undefined state
    IXDCGM_RET_TIMEOUT              = -11,  //!< Requested operation timed out
    IXDCGM_RET_VER_MISMATCH         = -12,  //!< Version mismatch between received and understood API
    IXDCGM_RET_UNKNOWN_FIELD        = -13,  //!< Unknown field id
    IXDCGM_RET_NO_DATA              = -14,  //!< No data is available
    IXDCGM_RET_STALE_DATA           = -15,  //!< Data is considered stale
    IXDCGM_RET_NOT_WATCHED          = -16,  //!< The given field id is not being updated by the cache manager
    IXDCGM_RET_NO_PERMISSION        = -17,  //!< Do not have permission to perform the desired action
    IXDCGM_RET_GPU_IS_LOST          = -18,  //!< GPU is no longer reachable
    IXDCGM_RET_RESET_REQUIRED       = -19,  //!< GPU requires a reset
    IXDCGM_RET_FUNCTION_NOT_FOUND   = -20,  //!< The function that was requested was not found (bindings only error)
    IXDCGM_RET_CONNECTION_NOT_VALID = -21,  //!< The connection to the host engine is not valid any longer
    IXDCGM_RET_GPU_NOT_SUPPORTED    = -22,  //!< This GPU is not supported by DCGM
    IXDCGM_RET_GROUP_INCOMPATIBLE = -23,  //!< The GPUs of the provided group are not compatible with each other for the
                                          //!< requested operation
    IXDCGM_RET_MAX_LIMIT                   = -24,  //!< Max limit reached for the object
    IXDCGM_RET_LIBRARY_NOT_FOUND           = -25,  //!< DCGM library could not be found
    IXDCGM_RET_DUPLICATE_KEY               = -26,  //!< Duplicate key passed to a function
    IXDCGM_RET_GPU_IN_SYNC_BOOST_GROUP     = -27,  //!< GPU is already a part of a sync boost group
    IXDCGM_RET_GPU_NOT_IN_SYNC_BOOST_GROUP = -28,  //!< GPU is not a part of a sync boost group
    IXDCGM_RET_REQUIRES_ROOT = -29,  //!< This operation cannot be performed when the host engine is running as non-root
    IXDCGM_RET_IXVS_ERROR    = -30,  //!< DCGM GPU Diagnostic was successfully executed, but reported an error.
    IXDCGM_RET_INSUFFICIENT_SIZE        = -31,  //!< An input argument is not large enough
    IXDCGM_RET_FIELD_UNSUPPORTED_BY_API = -32,  //!< The given field ID is not supported by the API being called
    IXDCGM_RET_MODULE_NOT_LOADED = -33,  //!< This request is serviced by a module of DCGM that is not currently loaded
    IXDCGM_RET_IN_USE            = -34,  //!< The requested operation could not be completed because the affected
                                         //!< resource is in use
    IXDCGM_RET_GROUP_IS_EMPTY =
        -35,  //!< This group is empty and the requested operation is not valid on an empty group
    IXDCGM_RET_PROFILING_NOT_SUPPORTED = -36,  //!< Profiling is not supported for this group of GPUs or GPU.
    IXDCGM_RET_PROFILING_LIBRARY_ERROR = -37,  //!< The third-party Profiling module returned an unrecoverable error.
    IXDCGM_RET_PROFILING_MULTI_PASS    = -38,  //!< The requested profiling metrics cannot be collected in a single pass
    IXDCGM_RET_DIAG_ALREADY_RUNNING    = -39,  //!< A diag instance is already running, cannot run a new diag until
                                               //!< the current one finishes.
    IXDCGM_RET_DIAG_BAD_JSON               = -40,  //!< The DCGM GPU Diagnostic returned JSON that cannot be parsed
    IXDCGM_RET_DIAG_BAD_LAUNCH             = -41,  //!< Error while launching the DCGM GPU Diagnostic
    IXDCGM_RET_DIAG_UNUSED                 = -42,  //!< Unused
    IXDCGM_RET_DIAG_THRESHOLD_EXCEEDED     = -43,  //!< A field value met or exceeded the error threshold.
    IXDCGM_RET_INSUFFICIENT_DRIVER_VERSION = -44,  //!< The installed driver version is insufficient for this API
    IXDCGM_RET_INSTANCE_NOT_FOUND          = -45,  //!< The specified GPU instance does not exist
    IXDCGM_RET_COMPUTE_INSTANCE_NOT_FOUND  = -46,  //!< The specified GPU compute instance does not exist
    IXDCGM_RET_CHILD_NOT_KILLED            = -47,  //!< Couldn't kill a child process within the retries
    IXDCGM_RET_3RD_PARTY_LIBRARY_ERROR     = -48,  //!< Detected an error in a 3rd-party library
    IXDCGM_RET_INSUFFICIENT_RESOURCES      = -49,  //!< Not enough resources available
    IXDCGM_RET_PLUGIN_EXCEPTION            = -50,  //!< Exception thrown from a diagnostic plugin
    IXDCGM_RET_IXVS_ISOLATE_ERROR    = -51,  //!< The diagnostic returned an error that indicates the need for isolation
    IXDCGM_RET_IXVS_BINARY_NOT_FOUND = -52,  //!< The NVVS binary was not found in the specified location
    IXDCGM_RET_IXVS_KILLED           = -53,  //!< The NVVS process was killed by a signal
    IXDCGM_RET_PAUSED                = -54,  //!< The hostengine and all modules are paused
    IXDCGM_RET_ALREADY_INITIALIZED   = -55,  //!< The object is already initialized
} ixdcgmReturn_t;

typedef enum
{
    ixdcgmLogLevelNone    = 0, /*!< No logging */
    ixdcgmLogLevelFatal   = 1, /*!< Fatal Errors */
    ixdcgmLogLevelError   = 2, /*!< Errors */
    ixdcgmLogLevelWarning = 3, /*!< Warnings */
    ixdcgmLogLevelInfo    = 4, /*!< Informative */
    ixdcgmLogLevelDebug   = 5, /*!< Debug information */
    ixdcgmLogLevelVerbose = 6  /*!< Verbose debugging information */
} ixdcgmLogLevel_t;

typedef uintptr_t ixdcgmHandle_t;    //!< Identifier for ixDCGM Handle
typedef uintptr_t ixdcgmGpuGrp_t;    //!< Identifier for a group of GPUs. A group can have one or more GPUs
typedef uintptr_t ixdcgmFieldGrp_t;  //!< Identifier for a group of fields.
typedef uintptr_t ixdcgmStatus_t;    //!< Identifier for list of status codes

typedef struct
{
    /* data */
    ixdcgmHandle_t   ixdcgmHandler;
    const char*      logFileName;
    ixdcgmLogLevel_t loglevelDefault;
} ixdcgmStartParams;

typedef struct
{
    unsigned int version;                /*!< Version number*/
    unsigned int persistAfterDisconnect; /*!< 1 = do not clean up after this connection.
                                              0 = clean up after this connection */
    unsigned int timeoutMs;              /*!< wait in milliseconds before giving up */
    unsigned int addressIsUnixSocket;    /*!< unix socket filename (1) or a TCP/IP address (0) */
} ixdcgmConnectParams;

#define ixdcgmConnectParams_version2 MAKE_IXDCGM_VERSION(ixdcgmConnectParams, 2)

typedef enum ixdcgmOperationMode_enum
{
    IXDCGM_OPERATION_MODE_AUTO   = 1,
    IXDCGM_OPERATION_MODE_MANUAL = 2
} ixdcgmOperationMode_t;

typedef enum
{
    ixdcgmModuleIdCore       = 0,  //!< Core DCGM - always loaded
    ixdcgmModuleIdNvSwitch   = 1,  //!< NvSwitch Module
    ixdcgmModuleIdVGPU       = 2,  //!< VGPU Module
    ixdcgmModuleIdIntrospect = 3,  //!< Introspection Module
    ixdcgmModuleIdHealth     = 4,  //!< Health Module
    ixdcgmModuleIdPolicy     = 5,  //!< Policy Module
    ixdcgmModuleIdConfig     = 6,  //!< Config Module
    ixdcgmModuleIdDiag       = 7,  //!< GPU Diagnostic Module
    ixdcgmModuleIdProfiling  = 8,  //!< Profiling Module
    ixdcgmModuleIdSysmon     = 9,  //!< System Monitoring Module

    ixdcgmModuleIdCount  //!< Always last. 1 greater than largest value above
} ixdcgmModuleId_t;

typedef enum ixdcgmOrder_enum
{
    IXDCGM_ORDER_ASCENDING  = 1,  //!< Data with earliest (lowest) timestamps returned first
    IXDCGM_ORDER_DESCENDING = 2   //!< Data with latest (highest) timestamps returned first
} ixdcgmOrder_t;

typedef enum
{
    ixdcgmModuleStatusNotLoaded  = 0,  //!< Module has not been loaded yet
    ixdcgmModuleStatusDenylisted = 1,  //!< Module is on the denylist; can't be loaded
    ixdcgmModuleStatusFailed     = 2,  //!< Loading the module failed
    ixdcgmModuleStatusLoaded     = 3,  //!< Module has been loaded
    ixdcgmModuleStatusUnloaded   = 4,  //!< Module has been unloaded, happens during shutdown
    ixdcgmModuleStatusPaused     = 5,  /*!< Module has been paused. This is a temporary state that will
                                          move to ixdcgmModuleStatusLoaded once the module is resumed.
                                          This status implies that the module is loaded. */
} ixdcgmModuleStatus_t;

typedef struct
{
    unsigned int          version;        /*!< Version number. Use ixdcgmStartEmbeddedV2Params_version2 */
    ixdcgmOperationMode_t opMode;         /*!< IN: Collect data automatically or manually when asked by the user. */
    ixdcgmHandle_t        ixdcgmHandle;   /*!< OUT: DCGM Handle to use for API calls */
    const char*           logFile;        /*!< IN: File that DCGM should log to. NULL = do not log. '-' = stdout */
    ixdcgmLogLevel_t      logLevel;       /*!< IN: Severity at which DCGM should log to logFile */
    unsigned int          denyListCount;  /*!< IN: Number of modules to be added to the denylist in denyList[] */
    const char*           serviceAccount; /*!< IN: Service account for unprivileged processes */
    ixdcgmModuleId_t      denyList[ixdcgmModuleIdCount]; /*!< IN: IDs of modules to be added to the denylist */
    char                  _padding[4];                   /*!< IN: Unused. Aligns the struct to 8 bytes. */
} ixdcgmStartEmbeddedParam;

typedef unsigned int ixdcgm_connection_id_t;
#define IXDCGM_CONNECTION_ID_NONE ((ixdcgm_connection_id_t)0)

#define IXDCGM_HOSTENGINE_DEFAULT_PORT 5777
#define IXDCGM_HOSTENGINE_LOCAL_ADDR "0.0.0.0"  // Default set to listen to ALL IP addrs

#define IXDCGM_EMBEDDED_HANDLE 0x7fffffff
#define IXDCGM_MAX_NUM_DEVICES 16
#define IXDCGM_MAX_NUM_GROUPS 64

#define IXDCGM_CMI_F_WATCHED 0x00000001 /* Is this field being watched? */

typedef struct
{
    unsigned int version;
    char         rawBuildInfoString[IXDCGM_MAX_STR_LENGTH * 2];
} ixdcgmVersionInfo_v2;

#define ixdcgmVersionInfo_version2 MAKE_IXDCGM_VERSION(ixdcgmVersionInfo_v2, 2)

#define ixdcgmVersionInfo_version ixdcgmVersionInfo_version2
typedef ixdcgmVersionInfo_v2 ixdcgmVersionInfo_t;

/**
 * Type of GPU groups
 */
typedef enum ixdcgmGroupType_enum
{
    IXDCGM_GROUP_DEFAULT                   = 0,  //!< All the GPUs on the node are added to the group
    IXDCGM_GROUP_EMPTY                     = 1,  //!< Creates an empty group
    IXDCGM_GROUP_DEFAULT_NVSWITCHES        = 2,  //!< All NvSwitches of the node are added to the group
    IXDCGM_GROUP_DEFAULT_INSTANCES         = 3,  //!< All GPU instances of the node are added to the group
    IXDCGM_GROUP_DEFAULT_COMPUTE_INSTANCES = 4,  //!< All compute instances of the node are added to the group
    IXDCGM_GROUP_DEFAULT_EVERYTHING        = 5,  //!< All entities are added to this default group
} ixdcgmGroupType_t;

/**
 * Identifies for special IXDCGM groups
 */
#define IXDCGM_GROUP_ALL_GPUS 0x7fffffff
#define IXDCGM_GROUP_ALL_NVSWITCHES 0x7ffffffe
#define IXDCGM_GROUP_ALL_INSTANCES 0x7ffffffd
#define IXDCGM_GROUP_ALL_COMPUTE_INSTANCES 0x7ffffffc
#define IXDCGM_GROUP_ALL_ENTITIES 0x7ffffffb

#define IXDCGM_MAX_CLOCKS 256
#define IXDCGM_GEGE_FLAG_ONLY_SUPPORTED 0x00000001

#define IXDCGM_GROUP_MAX_ENTITIES 64
#define IXDCGM_MAX_FIELD_IDS_PER_FIELD_GROUP 128
#define IXDCGM_MAX_NUM_FIELD_GROUPS 64

#define IXDCGM_FV_FLAG_LIVE_DATA 0x00000001
/**
 * Default maximum age of samples kept (usec)
 */
#define IXDCGM_MAX_AGE_USEC_DEFAULT 30000000

typedef struct
{
    int          version;   //!< Version Number (ixdcgmClockSet_version)
    unsigned int memClock;  //!< Memory Clock (Memory Clock value OR DCGM_INT32_BLANK to Ignore/Use compatible
                            //!< value with smClk)
    unsigned int smClock;  //!< SM Clock (SM Clock value OR DCGM_INT32_BLANK to Ignore/Use compatible value with memClk)
} ixdcgmClockSet_v1;

/**
 * Typedef for \ref ixdcgmClockSet_v1
 */
typedef ixdcgmClockSet_v1 ixdcgmClockSet_t;

/**
 * Version 1 for \ref ixdcgmClockSet_v1
 */
#define ixdcgmClockSet_version1 MAKE_IXDCGM_VERSION(ixdcgmClockSet_v1, 1)

/**
 * Latest version for \ref ixdcgmClockSet_t
 */
#define ixdcgmClockSet_version ixdcgmClockSet_version1

typedef struct
{
    unsigned int version;  //!< Version Number (ixdcgmDeviceSupportedClockSets_version)
    unsigned int count;    //!< Number of supported clocks
    ixdcgmClockSet_t
        clockSet[IXDCGM_MAX_CLOCKS];  //!< Valid clock sets for the device. Upto \ref count entries are filled
} ixdcgmDeviceSupportedClockSets_v1;
/**
 * Typedef for \ref ixdcgmDeviceSupportedClockSets_v1
 */
typedef ixdcgmDeviceSupportedClockSets_v1 ixdcgmDeviceSupportedClockSets_t;

/**
 * Version 1 for \ref ixdcgmDeviceSupportedClockSets_v1
 */
#define ixdcgmDeviceSupportedClockSets_version1 MAKE_IXDCGM_VERSION(ixdcgmDeviceSupportedClockSets_v1, 1)

/**
 * Latest version for \ref ixdcgmDeviceSupportedClockSets_t
 */
#define ixdcgmDeviceSupportedClockSets_version ixdcgmDeviceSupportedClockSets_version1

typedef struct
{
    ixdcgm_field_entity_group_t entityGroupId;  //!< Entity Group ID entity belongs to
    ixdcgm_field_eid_t          entityId;       //!< Entity ID of the entity
} ixdcgmGroupEntityPair_t;

typedef struct
{
    unsigned int            version;                                //!< Version Number (use ixdcgmGroupInfo_version2)
    unsigned int            count;                                  //!< count of entityIds returned in \a entityList
    char                    groupName[IXDCGM_MAX_STR_LENGTH];       //!< Group Name
    ixdcgmGroupEntityPair_t entityList[IXDCGM_GROUP_MAX_ENTITIES];  //!< List of the entities that are in this group
} ixdcgmGroupInfo_v2;

/**
 * Typedef for \ref ixdcgmGroupInfo_v2
 */
typedef ixdcgmGroupInfo_v2 ixdcgmGroupInfo_t;

/**
 * Version 2 for \ref ixdcgmGroupInfo_v2
 */
#define ixdcgmGroupInfo_version2 MAKE_IXDCGM_VERSION(ixdcgmGroupInfo_v2, 2)

/**
 * Latest version for \ref ixdcgmGroupInfo_t
 */
#define ixdcgmGroupInfo_version ixdcgmGroupInfo_version2

typedef struct
{
    unsigned int version;       //!< Version Number
    unsigned int slowdownTemp;  //!< Slowdown temperature
    unsigned int shutdownTemp;  //!< Shutdown temperature
} ixdcgmDeviceThermals;
typedef ixdcgmDeviceThermals ixdcgmDeviceThermals_t;

typedef struct
{
    unsigned int version;             //!< Version Number
    unsigned int curPowerLimit;       //!< Power management limit associated with this device (in W)
    unsigned int defaultPowerLimit;   //!< Power management limit effective at device boot (in W)
    unsigned int enforcedPowerLimit;  //!< Effective power limit that the driver enforces after taking into account
                                      //!< all limiters (in W)
    unsigned int minPowerLimit;       //!< Minimum power management limit (in W)
    unsigned int maxPowerLimit;       //!< Maximum power management limit (in W)
} ixdcgmDevicePowerLimits;
typedef ixdcgmDevicePowerLimits ixdcgmDevicePowerLimits_t;

typedef struct
{
    unsigned int version;                                     //!< Version Number (ixdcgmDeviceIdentifiers_version)
    char         brandName[IXDCGM_MAX_STR_LENGTH];            //!< Brand Name
    char         deviceName[IXDCGM_MAX_STR_LENGTH];           //!< Name of the device
    char         pciBusId[IXDCGM_MAX_STR_LENGTH];             //!< PCI Bus ID
    char         serial[IXDCGM_MAX_STR_LENGTH];               //!< Serial for the device
    char         uuid[IXDCGM_MAX_STR_LENGTH];                 //!< UUID for the device
    char         vbios[IXDCGM_MAX_STR_LENGTH];                //!< VBIOS version
    char         inforomImageVersion[IXDCGM_MAX_STR_LENGTH];  //!< Inforom Image version
    unsigned int pciDeviceId;                                 //!< The combined 16-bit device id and 16-bit vendor id
    unsigned int pciSubSystemId;                              //!< The 32-bit Sub System Device ID
    char         driverVersion[IXDCGM_MAX_STR_LENGTH];        //!< Driver Version
    unsigned int virtualizationMode;                          //!< Virtualization Mode
} ixdcgmDeviceIdentifiers_v1;
typedef ixdcgmDeviceIdentifiers_v1 ixdcgmDeviceIdentifiers_t;

typedef struct
{
    unsigned int version;    //!< Version Number (ixdcgmDeviceMemoryUsage_version)
    unsigned int bar1Total;  //!< Total BAR1 size in megabytes
    unsigned int fbTotal;    //!< Total framebuffer memory in megabytes
    unsigned int fbUsed;     //!< Used framebuffer memory in megabytes
    unsigned int fbFree;     //!< Free framebuffer memory in megabytes
} ixdcgmDeviceMemoryUsage_v1;
typedef ixdcgmDeviceMemoryUsage_v1 ixdcgmDeviceMemoryUsage_t;

typedef struct
{
    unsigned int version;
    unsigned int persistenceModeEnabled;
    unsigned int migModeEnabled;
    unsigned int confidentialComputeMode;
} ixdcgmDeviceSettings;

typedef ixdcgmDeviceSettings ixdcgmDeviceSettings_t;

typedef struct
{
    unsigned int                     version;          //!< Version number (ixdcgmDeviceAttributes_version)
    ixdcgmDeviceSupportedClockSets_t clockSets;        //!< Supported clocks for the device
    ixdcgmDeviceThermals_t           thermalSettings;  //!< Thermal settings for the device
    ixdcgmDevicePowerLimits_t        powerLimits;      //!< Various power limits for the device
    ixdcgmDeviceIdentifiers_t        identifiers;      //!< Identifiers for the device
    ixdcgmDeviceMemoryUsage_t        memoryUsage;      //!< Memory usage info for the device
    ixdcgmDeviceSettings_t           settings;         //!< Basic device settings
} ixdcgmDeviceAttributes;

typedef ixdcgmDeviceAttributes ixdcgmDeviceAttributes_t;
#define ixdcgmDeviceAttributes_version3 MAKE_IXDCGM_VERSION(ixdcgmDeviceAttributes, 3)
#define ixdcgmDeviceAttributes_version ixdcgmDeviceAttributes_version3

typedef struct
{
    // version must always be first
    unsigned int version;  //!< version number (ixdcgmFieldValue_version1)

    unsigned short fieldId;    //!< One of IXDCGM_FI_?
    unsigned short fieldType;  //!< One of IXDCGM_FT_?
    int            status;     //!< Status for the querying the field. IXDCGM_ST_OK or one of IXDCGM_ST_?
    int64_t        ts;         //!< Timestamp in usec since 1970
    union {
        int64_t i64;                           //!< Int64 value
        double  dbl;                           //!< Double value
        char    str[IXDCGM_MAX_STR_LENGTH];    //!< NULL terminated string
        char    blob[IXDCGM_MAX_BLOB_LENGTH];  //!< Binary blob
    } value;                                   //!< Value
} ixdcgmFieldValue_v1;
#define ixdcgmFieldValue_version1 MAKE_IXDCGM_VERSION(ixdcgmFieldValue_v1, 1)

typedef struct
{
    // version must always be first
    unsigned int                version;        //!< version number (ixdcgmFieldValue_version2)
    ixdcgm_field_entity_group_t entityGroupId;  //!< Entity group this field value's entity belongs to
    ixdcgm_field_eid_t          entityId;       //!< Entity this field value belongs to
    unsigned short              fieldId;        //!< One of IXDCGM_FI_?
    unsigned short              fieldType;      //!< One of IXDCGM_FT_?
    int                         status;  //!< Status for the querying the field. IXDCGM_ST_OK or one of IXDCGM_ST_?
    unsigned int                unused;  //!< Unused for now to align ts to an 8-byte boundary.
    int64_t                     ts;      //!< Timestamp in usec since 1970
    union {
        int64_t i64;                           //!< Int64 value
        double  dbl;                           //!< Double value
        char    str[IXDCGM_MAX_STR_LENGTH];    //!< NULL terminated string
        char    blob[IXDCGM_MAX_BLOB_LENGTH];  //!< Binary blob
    } value;                                   //!< Value
} ixdcgmFieldValue_v2;
#define ixdcgmFieldValue_version2 MAKE_IXDCGM_VERSION(ixdcgmFieldValue_v2, 2)

/* Bitmask values for ixdcgmGetFieldIdSummary - Sync with DcgmcmSummaryType_t */
#define IXDCGM_SUMMARY_MIN 0x00000001
#define IXDCGM_SUMMARY_MAX 0x00000002
#define IXDCGM_SUMMARY_AVG 0x00000004
#define IXDCGM_SUMMARY_SUM 0x00000008
#define IXDCGM_SUMMARY_COUNT 0x00000010
#define IXDCGM_SUMMARY_INTEGRAL 0x00000020
#define IXDCGM_SUMMARY_DIFF 0x00000040
#define IXDCGM_SUMMARY_SIZE 7

/* ixdcgmSummaryResponse_t is part of ixdcgmFieldSummaryRequest, so it uses ixdcgmFieldSummaryRequest's version. */

typedef struct
{
    unsigned int fieldType;     //!< type of field that is summarized (int64 or fp64)
    unsigned int summaryCount;  //!< the number of populated summaries in \ref values
    union {
        int64_t i64;
        double  fp64;
    } values[IXDCGM_SUMMARY_SIZE];  //!< array for storing the values of each summary. The summaries are stored
                                    //!< in order. For example, if MIN AND MAX are requested, then 0 will be MIN
                                    //!< and 1 will be MAX. If AVG and DIFF were requested, then AVG would be 0
                                    //!< and 1 would be DIFF
} ixdcgmSummaryResponse_t;

typedef struct
{
    unsigned int                version;          //!< version of this message - ixdcgmFieldSummaryRequest_v1
    unsigned short              fieldId;          //!< field id to be summarized
    ixdcgm_field_entity_group_t entityGroupId;    //!< the type of entity whose field we're getting
    ixdcgm_field_eid_t          entityId;         //!< ordinal id for this entity
    uint32_t                    summaryTypeMask;  //!< bit-mask of IXDCGM_SUMMARY_*, the requested summaries
    uint64_t                    startTime;        //!< start time for the interval being summarized. 0 means to use
                                                  //!< any data before.
    uint64_t endTime;                             //!< end time for the interval being summarized. 0 means to use
                                                  //!< any data after.
    ixdcgmSummaryResponse_t response;             //!< response data for this request
} ixdcgmFieldSummaryRequest;

typedef ixdcgmFieldSummaryRequest ixdcgmFieldSummaryRequest_t;

#define ixdcgmFieldSummaryRequest_version1 MAKE_IXDCGM_VERSION(ixdcgmFieldSummaryRequest, 1)

typedef struct
{
    ixdcgmFieldSummaryRequest_t fsr;     //!< IN/OUT: field summary populated on success
    unsigned int                cmdRet;  //!< OUT: Error code generated
} ixdcgmGetFieldSummary_v1;

typedef int (*ixdcgmFieldValueEntityEnumeration_f)(ixdcgm_field_entity_group_t entityGroupId,
                                                   ixdcgm_field_eid_t          entityId,
                                                   ixdcgmFieldValue_v1*        values,
                                                   int                         numValues,
                                                   void*                       userData);
typedef enum ixdcgmPerGpuTestIndices_enum
{
    IXDCGM_MEMORY_INDEX           = 0,  //!< Memory test index
    IXDCGM_DIAGNOSTIC_INDEX       = 1,  //!< Diagnostic test index
    IXDCGM_PCIE_INDEX             = 2,  //!< PCIe test index
    IXDCGM_SM_STRESS_INDEX        = 3,  //!< SM Stress test index
    IXDCGM_TARGETED_STRESS_INDEX  = 4,  //!< Targeted Stress test index
    IXDCGM_TARGETED_POWER_INDEX   = 5,  //!< Targeted Power test index
    IXDCGM_MEMORY_BANDWIDTH_INDEX = 6,  //!< Memory bandwidth test index
    IXDCGM_MEMTEST_INDEX          = 7,  //!< Memtest test index
    IXDCGM_PULSE_TEST_INDEX       = 8,  //!< Pulse test index
    IXDCGM_EUD_TEST_INDEX         = 9,  //!< EUD test index
    // Remaining tests are included for convenience but have different execution rules
    // See IXDCGM_PER_GPU_TEST_COUNT
    IXDCGM_UNUSED2_TEST_INDEX   = 10,
    IXDCGM_UNUSED3_TEST_INDEX   = 11,
    IXDCGM_UNUSED4_TEST_INDEX   = 12,
    IXDCGM_UNUSED5_TEST_INDEX   = 13,
    IXDCGM_SOFTWARE_INDEX       = 14,  //!< Software test index
    IXDCGM_CONTEXT_CREATE_INDEX = 15,  //!< Context create test index
    IXDCGM_UNKNOWN_INDEX        = 16   //!< Unknown test
} ixdcgmPerGpuTestIndices_t;

typedef enum ixdcgmChipArchitecture_enum
{
    IXDCGM_CHIP_ARCH_OLDER     = 1,  //!< All GPUs older than Kepler
    IXDCGM_CHIP_ARCH_NVKEPLER  = 2,  //!< All Kepler-architecture parts
    IXDCGM_CHIP_ARCH_NVMAXWELL = 3,  //!< All Maxwell-architecture parts
    IXDCGM_CHIP_ARCH_NVPASCAL  = 4,  //!< All Pascal-architecture parts
    IXDCGM_CHIP_ARCH_NVVOLTA   = 5,  //!< All Volta-architecture parts
    IXDCGM_CHIP_ARCH_NVTURING  = 6,  //!< All Turing-architecture parts
    IXDCGM_CHIP_ARCH_NVAMPERE  = 7,  //!< All Ampere-architecture parts
    IXDCGM_CHIP_ARCH_NVADA     = 8,  //!< All Ada-architecture parts
    IXDCGM_CHIP_ARCH_NVHOPPER  = 9,  //!< All Hopper-architecture parts

    IXDCGM_CHIP_ARCH_IX_BI = 100,
    IXDCGM_CHIP_ARCH_IX_MR = 101,
    IXDCGM_CHIP_ARCH_COUNT,                //!< Keep this second to last, exclude unknown
    IXDCGM_CHIP_ARCH_UNKNOWN = 0xffffffff  //!< Anything else, presumably something newer
} ixdcgmChipArchitecture_t;

typedef enum
{
    IXDCGM_GPU_VIRTUALIZATION_MODE_NONE        = 0,  //!< Represents Bare Metal GPU
    IXDCGM_GPU_VIRTUALIZATION_MODE_PASSTHROUGH = 1,  //!< Device is associated with GPU-Passthrough
    IXDCGM_GPU_VIRTUALIZATION_MODE_VGPU        = 2,  //!< Device is associated with vGPU inside virtual machine.
    IXDCGM_GPU_VIRTUALIZATION_MODE_HOST_VGPU   = 3,  //!< Device is associated with VGX hypervisor in vGPU mode
    IXDCGM_GPU_VIRTUALIZATION_MODE_HOST_VSGA   = 4,  //!< Device is associated with VGX hypervisor in vSGA mode
} ixdcgmGpuVirtualizationMode_t;

typedef struct
{
    unsigned int syncBoost;  //!< Sync Boost Mode (0: Disabled, 1 : Enabled, DCGM_INT32_BLANK : Ignored). Note that
                             //!< using this setting may result in lower clocks than targetClocks
    ixdcgmClockSet_t targetClocks;  //!< Target clocks. Set smClock and memClock to DCGM_INT32_BLANK to ignore/use
                                    //!< compatible values. For GPUs > Maxwell, setting this implies autoBoost=0
} ixdcgmConfigPerfStateSettings_t;

typedef enum ixdcgmConfigPowerLimitType_enum
{
    IXDCGM_CONFIG_POWER_CAP_INDIVIDUAL = 0,  //!< Represents the power cap to be applied for each member of the group
    IXDCGM_CONFIG_POWER_BUDGET_GROUP   = 1,  //!< Represents the power budget for the entire group
} ixdcgmConfigPowerLimitType_t;

typedef struct
{
    ixdcgmConfigPowerLimitType_t
                 type;  //!< Flag to represent power cap for each GPU or power budget for the group of GPUs
    unsigned int val;   //!< Power Limit in Watts (Set a value OR DCGM_INT32_BLANK to Ignore)
} ixdcgmConfigPowerLimit_t;

typedef struct
{
    unsigned int version;      //!< Version number (ixdcgmConfig_version)
    unsigned int gpuId;        //!< GPU ID
    unsigned int eccMode;      //!< ECC Mode  (0: Disabled, 1 : Enabled, DCGM_INT32_BLANK : Ignored)
    unsigned int computeMode;  //!< Compute Mode (One of DCGM_CONFIG_COMPUTEMODE_? OR DCGM_INT32_BLANK to Ignore)
    ixdcgmConfigPerfStateSettings_t perfState;   //!< Performance State Settings (clocks / boost mode)
    ixdcgmConfigPowerLimit_t        powerLimit;  //!< Power Limits
} ixdcgmConfig_v1;

typedef ixdcgmConfig_v1 ixdcgmConfig_t;

#define ixdcgmConfig_version1 MAKE_IXDCGM_VERSION(ixdcgmConfig_v1, 1)

#define ixdcgmConfig_version ixdcgmConfig_version1

typedef enum ixdcgmConfigType_enum
{
    IXDCGM_CONFIG_TARGET_STATE  = 0,  //!< The target configuration values to be applied
    IXDCGM_CONFIG_CURRENT_STATE = 1,  //!< The current configuration state
} ixdcgmConfigType_t;

typedef enum ixdcgmLinkState_enum
{
    ixdcgmLinkStateNotSupported = 0,  //!< Link is unsupported by this GPU (Default for GPUs)
    ixdcgmLinkStateDisabled     = 1,  //!< Link is supported for this link but this link is disabled
    ixdcgmLinkStateDown         = 2,  //!< This Link link is down (inactive)
    ixdcgmLinkStateUp           = 3   //!< This Link link is up (active)
} ixdcgmLinkState_t;

#define IXDCGM_MAX_LINKS_PER_GPU 16

typedef struct
{
    unsigned int version;            //!< Version Number. Should match ixdcgmDevicePidAccountingStats_version
    unsigned int pid;                //!< Process id of the process these stats are for
    unsigned int gpuUtilization;     //!< Percent of time over the process's lifetime during which one or more kernels
                                     //!< was executing on the GPU.
                                     //!< Set to DCGM_INT32_NOT_SUPPORTED if is not supported
    unsigned int memoryUtilization;  //!< Percent of time over the process's lifetime during which global (device)
                                     //!< memory was being read or written.
                                     //!< Set to DCGM_INT32_NOT_SUPPORTED if is not supported
    unsigned long long maxMemoryUsage;  //!< Maximum total memory in bytes that was ever allocated by the process.
                                        //!< Set to DCGM_INT64_NOT_SUPPORTED if is not supported
    unsigned long long startTimestamp;  //!< CPU Timestamp in usec representing start time for the process
    unsigned long long activeTimeUsec;  //!< Amount of time in usec during which the compute context was active.
                                        //!< Note that this does not mean the context was being used. endTimestamp
                                        //!< can be computed as startTimestamp + activeTime
} ixdcgmDevicePidAccountingStats_v1;

/**
 * Typedef for \ref ixdcgmDevicePidAccountingStats_v1
 */
typedef ixdcgmDevicePidAccountingStats_v1 ixdcgmDevicePidAccountingStats_t;
#define ixdcgmDevicePidAccountingStats_version1 MAKE_IXDCGM_VERSION(ixdcgmDevicePidAccountingStats_v1, 1)
#define ixdcgmDevicePidAccountingStats_version ixdcgmDevicePidAccountingStats_version1

typedef struct
{
    unsigned int pid;
    double       smUtil;
    double       memUtil;
} ixdcgmProcessUtilInfo_t;

typedef struct
{
    double       util;
    unsigned int pid;
} ixdcgmProcessUtilSample_t;

typedef struct
{
    unsigned int version;  //!< Version Number (ixdcgmDeviceVgpuProcessUtilInfo_version)
    union {
        unsigned int vgpuId;                   //!< vGPU instance ID
        unsigned int vgpuProcessSamplesCount;  //!< Count of processes running in the vGPU VM,for which utilization
                                               //!< rates are being reported in this cycle.
    } vgpuProcessUtilInfo;
    unsigned int pid;                                        //!< Process ID of the process running in the vGPU VM.
    char         processName[IXDCGM_VGPU_NAME_BUFFER_SIZE];  //!< Process Name of process running in the vGPU VM.
    unsigned int smUtil;                                     //!< GPU utilization of process running in the vGPU VM.
    unsigned int memUtil;                                    //!< Memory utilization of process running in the vGPU VM.
    unsigned int encUtil;                                    //!< Encoder utilization of process running in the vGPU VM.
    unsigned int decUtil;                                    //!< Decoder utilization of process running in the vGPU VM.
} ixdcgmDeviceVgpuProcessUtilInfo_v1;

/**
 * Typedef for \ref ixdcgmDeviceVgpuProcessUtilInfo_v1
 */
typedef ixdcgmDeviceVgpuProcessUtilInfo_v1 ixdcgmDeviceVgpuProcessUtilInfo_t;

/**
 * Version 1 for \ref ixdcgmDeviceVgpuProcessUtilInfo_v1
 */
#define ixdcgmDeviceVgpuProcessUtilInfo_version1 MAKE_IXDCGM_VERSION(ixdcgmDeviceVgpuProcessUtilInfo_v1, 1)

typedef enum ixdcgmGpuLevel_enum
{
    IXDCGM_TOPOLOGY_UNINITIALIZED = 0x0,

    /** \name PCI connectivity states */
    /**@{*/
    IXDCGM_TOPOLOGY_BOARD      = 0x1,  //!< multi-GPU board
    IXDCGM_TOPOLOGY_SINGLE     = 0x2,  //!< all devices that only need traverse a single PCIe switch
    IXDCGM_TOPOLOGY_MULTIPLE   = 0x4,  //!< all devices that need not traverse a host bridge
    IXDCGM_TOPOLOGY_HOSTBRIDGE = 0x8,  //!< all devices that are connected to the same host bridge
    IXDCGM_TOPOLOGY_CPU = 0x10,  //!< all devices that are connected to the same CPU but possibly multiple host bridges
    IXDCGM_TOPOLOGY_SYSTEM = 0x20,  //!< all devices in the system
    /**@}*/

    /** \name LINK connectivity states */
    /**@{*/
    IXDCGM_TOPOLOGY_LINK1  = 0x0100,     //!< GPUs connected via a single LINK link
    IXDCGM_TOPOLOGY_LINK2  = 0x0200,     //!< GPUs connected via two LINK links
    IXDCGM_TOPOLOGY_LINK3  = 0x0400,     //!< GPUs connected via three LINK links
    IXDCGM_TOPOLOGY_LINK4  = 0x0800,     //!< GPUs connected via four LINK links
    IXDCGM_TOPOLOGY_LINK5  = 0x1000,     //!< GPUs connected via five LINK links
    IXDCGM_TOPOLOGY_LINK6  = 0x2000,     //!< GPUs connected via six LINK links
    IXDCGM_TOPOLOGY_LINK7  = 0x4000,     //!< GPUs connected via seven LINK links
    IXDCGM_TOPOLOGY_LINK8  = 0x8000,     //!< GPUs connected via eight LINK links
    IXDCGM_TOPOLOGY_LINK9  = 0x10000,    //!< GPUs connected via nine LINK links
    IXDCGM_TOPOLOGY_LINK10 = 0x20000,    //!< GPUs connected via ten LINK links
    IXDCGM_TOPOLOGY_LINK11 = 0x40000,    //!< GPUs connected via 11 LINK links
    IXDCGM_TOPOLOGY_LINK12 = 0x80000,    //!< GPUs connected via 12 LINK links
    IXDCGM_TOPOLOGY_LINK13 = 0x100000,   //!< GPUs connected via 13 LINK links
    IXDCGM_TOPOLOGY_LINK14 = 0x200000,   //!< GPUs connected via 14 LINK links
    IXDCGM_TOPOLOGY_LINK15 = 0x400000,   //!< GPUs connected via 15 LINK links
    IXDCGM_TOPOLOGY_LINK16 = 0x800000,   //!< GPUs connected via 16 LINK links
    IXDCGM_TOPOLOGY_LINK17 = 0x1000000,  //!< GPUs connected via 17 LINK links
    IXDCGM_TOPOLOGY_LINK18 = 0x2000000,  //!< GPUs connected via 18 LINK links
    /**@}*/
} ixdcgmGpuTopologyLevel_t;

// the PCI paths are the lower 8 bits of the path information
#define IXDCGM_TOPOLOGY_PATH_PCI(x) (ixdcgmGpuTopologyLevel_t)((unsigned int)(x) & 0xFF)

// the LINK paths are the upper 24 bits of the path information
#define IXDCGM_TOPOLOGY_PATH_LINK(x) (ixdcgmGpuTopologyLevel_t)((unsigned int)(x) & 0xFFFFFF00)

#define IXDCGM_AFFINITY_BITMASK_ARRAY_SIZE 8
/**
 * Device topology information
 */
typedef struct
{
    unsigned int version;  //!< version number (ixdcgmDeviceTopology_version)

    unsigned long cpuAffinityMask[IXDCGM_AFFINITY_BITMASK_ARRAY_SIZE];  //!< affinity mask for the specified GPU
                                                                        //!< a 1 represents affinity to the CPU in that
                                                                        //!< bit position supports up to 256 cores
    unsigned int numGpus;                                               //!< number of valid entries in gpuPaths

    struct
    {
        unsigned int             gpuId;  //!< gpuId to which the path represents
        ixdcgmGpuTopologyLevel_t path;   //!< path to the gpuId from this GPU. Note that this is a bit-mask
                                         //!< of IXDCGM_TOPOLOGY_* values and can contain both PCIe topology
                                         //!< and NvLink topology where applicable. For instance:
                                         //!< 0x210 = IXDCGM_TOPOLOGY_CPU | IXDCGM_TOPOLOGY_LINK2
                                         //!< Use the macros IXDCGM_TOPOLOGY_PATH_LINK and
                                         //!< IXDCGM_TOPOLOGY_PATH_PCI to mask the NvLink and PCI paths, respectively.
        unsigned int localNvLinkIds;     //!< bits representing the local links connected to gpuId
                                         //!< e.g. if this field == 3, links 0 and 1 are connected,
                                         //!< field is only valid if LINKS actually exist between GPUs
    } gpuPaths[IXDCGM_MAX_NUM_DEVICES - 1];
} ixdcgmDeviceTopology_v1;

/**
 * Typedef for \ref ixdcgmDeviceTopology_v1
 */
typedef ixdcgmDeviceTopology_v1 ixdcgmDeviceTopology_t;

/**
 * Version 1 for \ref ixdcgmDeviceTopology_v1
 */
#define ixdcgmDeviceTopology_version1 MAKE_IXDCGM_VERSION(ixdcgmDeviceTopology_v1, 1)

/**
 * Latest version for \ref ixdcgmDeviceTopology_t
 */
#define ixdcgmDeviceTopology_version ixdcgmDeviceTopology_version1

/**
 * Group topology information
 */
typedef struct
{
    unsigned int version;  //!< version number (ixdcgmGroupTopology_version)

    unsigned long
        groupCpuAffinityMask[IXDCGM_AFFINITY_BITMASK_ARRAY_SIZE];  //!< the CPU affinity mask for all GPUs in the group
                                                                   //!< a 1 represents affinity to the CPU in that bit
                                                                   //!< position supports up to 256 cores
    unsigned int numaOptimalFlag;                                  //!< a zero value indicates that 1 or more GPUs
                                   //!< in the group have a different CPU affinity and thus
                                   //!< may not be optimal for certain algorithms
    ixdcgmGpuTopologyLevel_t slowestPath;  //!< the slowest path amongst GPUs in the group
} ixdcgmGroupTopology_v1;

/**
 * Typedef for \ref ixdcgmGroupTopology_v1
 */
typedef ixdcgmGroupTopology_v1 ixdcgmGroupTopology_t;

/**
 * Version 1 for \ref ixdcgmGroupTopology_v1
 */
#define ixdcgmGroupTopology_version1 MAKE_IXDCGM_VERSION(ixdcgmGroupTopology_v1, 1)

/**
 * Latest version for \ref ixdcgmGroupTopology_t
 */
#define ixdcgmGroupTopology_version ixdcgmGroupTopology_version1

/**
 * Running process information for a compute or graphics process
 */
typedef struct
{
    unsigned int       version;     //!< Version of this message (ixdcgmRunningProcess_version)
    unsigned int       pid;         //!< PID of the process
    unsigned long long memoryUsed;  //!< GPU memory used by this process in bytes.
} ixdcgmRunningProcess_v1;

/**
 * Typedef for \ref ixdcgmRunningProcess_v1
 */
typedef ixdcgmRunningProcess_v1 ixdcgmRunningProcess_t;

/**
 * Version 1 for \ref ixdcgmRunningProcess_v1
 */
#define ixdcgmRunningProcess_version1 MAKE_IXDCGM_VERSION(ixdcgmRunningProcess_v1, 1)

/**
 * Latest version for \ref ixdcgmRunningProcess_t
 */
#define ixdcgmRunningProcess_version ixdcgmRunningProcess_version1

typedef struct
{
    unsigned int     version;                                //!< Version number (ixdcgmFieldGroupInfo_version)
    unsigned int     numFieldIds;                            //!< Number of entries in fieldIds[] that are valid
    ixdcgmFieldGrp_t fieldGroupId;                           //!< ID of this field group
    char             fieldGroupName[IXDCGM_MAX_STR_LENGTH];  //!< Field Group Name
    unsigned short   fieldIds[IXDCGM_MAX_FIELD_IDS_PER_FIELD_GROUP];  //!< Field ids that belong to this group
} ixdcgmFieldGroupInfo_v1;

typedef ixdcgmFieldGroupInfo_v1 ixdcgmFieldGroupInfo_t;

/**
 * Version 1 for ixdcgmFieldGroupInfo_v1
 */
#define ixdcgmFieldGroupInfo_version1 MAKE_IXDCGM_VERSION(ixdcgmFieldGroupInfo_v1, 1)

/**
 * Latest version for ixdcgmFieldGroupInfo_t
 */
#define ixdcgmFieldGroupInfo_version ixdcgmFieldGroupInfo_version1

typedef struct
{
    unsigned int           version;         //!< Version number (ixdcgmAllFieldGroupInfo_version)
    unsigned int           numFieldGroups;  //!< Number of entries in fieldGroups[] that are populated
    ixdcgmFieldGroupInfo_t fieldGroups[IXDCGM_MAX_NUM_FIELD_GROUPS];  //!< Info about each field group
} ixdcgmAllFieldGroup_v1;

typedef ixdcgmAllFieldGroup_v1 ixdcgmAllFieldGroup_t;

/**
 * Version 1 for ixdcgmAllFieldGroup_v1
 */
#define ixdcgmAllFieldGroup_version1 MAKE_IXDCGM_VERSION(ixdcgmAllFieldGroup_v1, 1)

/**
 * Latest version for ixdcgmAllFieldGroup_t
 */
#define ixdcgmAllFieldGroup_version ixdcgmAllFieldGroup_version1

typedef struct
{
    int              targetLogger;
    ixdcgmLogLevel_t targetLogLevel;
} ixdcgmSettingsSetLoggingSeverity_v1;

#define ixdcgmSettingsSetLoggingSeverity_version1 MAKE_IXDCGM_VERSION(ixdcgmSettingsSetLoggingSeverity_v1, 1)
#define ixdcgmSettingsSetLoggingSeverity_version ixdcgmSettingsSetLoggingSeverity_version1
typedef ixdcgmSettingsSetLoggingSeverity_v1 ixdcgmSettingsSetLoggingSeverity_t;

#endif  // end of __IXDCGM_STRUCTS_H__