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

typedef enum ixdcgmReturn_enum
{
    IXDCGM_RET_OK = 0,                            //!< Success
    IXDCGM_RET_BADPARAM = -1,                     //!< A bad parameter was passed to a function
    IXDCGM_RET_GENERIC_ERROR = -3,                //!< A generic, unspecified error
    IXDCGM_RET_MEMORY = -4,                       //!< An out of memory error occurred
    IXDCGM_RET_NOT_CONFIGURED = -5,               //!< Setting not configured
    IXDCGM_RET_NOT_SUPPORTED = -6,                //!< Feature not supported
    IXDCGM_RET_INIT_ERROR = -7,                   //!< DCGM Init error
    IXDCGM_RET_NVML_ERROR = -8,                   //!< When NVML returns error
    IXDCGM_RET_PENDING = -9,                      //!< Object is in pending state of something else
    IXDCGM_RET_UNINITIALIZED = -10,               //!< Object is in undefined state
    IXDCGM_RET_TIMEOUT = -11,                     //!< Requested operation timed out
    IXDCGM_RET_VER_MISMATCH = -12,                //!< Version mismatch between received and understood API
    IXDCGM_RET_UNKNOWN_FIELD = -13,               //!< Unknown field id
    IXDCGM_RET_NO_DATA = -14,                     //!< No data is available
    IXDCGM_RET_STALE_DATA = -15,                  //!< Data is considered stale
    IXDCGM_RET_NOT_WATCHED = -16,                 //!< The given field id is not being updated by the cache manager
    IXDCGM_RET_NO_PERMISSION = -17,               //!< Do not have permission to perform the desired action
    IXDCGM_RET_GPU_IS_LOST = -18,                 //!< GPU is no longer reachable
    IXDCGM_RET_RESET_REQUIRED = -19,              //!< GPU requires a reset
    IXDCGM_RET_FUNCTION_NOT_FOUND = -20,          //!< The function that was requested was not found (bindings only error)
    IXDCGM_RET_CONNECTION_NOT_VALID = -21,        //!< The connection to the host engine is not valid any longer
    IXDCGM_RET_GPU_NOT_SUPPORTED = -22,           //!< This GPU is not supported by DCGM
    IXDCGM_RET_GROUP_INCOMPATIBLE = -23,          //!< The GPUs of the provided group are not compatible with each other for the
                                                  //!< requested operation
    IXDCGM_RET_MAX_LIMIT = -24,                   //!< Max limit reached for the object
    IXDCGM_RET_LIBRARY_NOT_FOUND = -25,           //!< DCGM library could not be found
    IXDCGM_RET_DUPLICATE_KEY = -26,               //!< Duplicate key passed to a function
    IXDCGM_RET_GPU_IN_SYNC_BOOST_GROUP = -27,     //!< GPU is already a part of a sync boost group
    IXDCGM_RET_GPU_NOT_IN_SYNC_BOOST_GROUP = -28, //!< GPU is not a part of a sync boost group
    IXDCGM_RET_REQUIRES_ROOT = -29,               //!< This operation cannot be performed when the host engine is running as non-root
    IXDCGM_RET_IXVS_ERROR = -30,                  //!< DCGM GPU Diagnostic was successfully executed, but reported an error.
    IXDCGM_RET_INSUFFICIENT_SIZE = -31,           //!< An input argument is not large enough
    IXDCGM_RET_FIELD_UNSUPPORTED_BY_API = -32,    //!< The given field ID is not supported by the API being called
    IXDCGM_RET_MODULE_NOT_LOADED = -33,           //!< This request is serviced by a module of DCGM that is not currently loaded
    IXDCGM_RET_IN_USE = -34,                      //!< The requested operation could not be completed because the affected
                                                  //!< resource is in use
    IXDCGM_RET_GROUP_IS_EMPTY = -35,              //!< This group is empty and the requested operation is not valid on an empty group
    IXDCGM_RET_PROFILING_NOT_SUPPORTED = -36,     //!< Profiling is not supported for this group of GPUs or GPU.
    IXDCGM_RET_PROFILING_LIBRARY_ERROR = -37,     //!< The third-party Profiling module returned an unrecoverable error.
    IXDCGM_RET_PROFILING_MULTI_PASS = -38,        //!< The requested profiling metrics cannot be collected in a single pass
    IXDCGM_RET_DIAG_ALREADY_RUNNING = -39,        //!< A diag instance is already running, cannot run a new diag until
                                                  //!< the current one finishes.
    IXDCGM_RET_DIAG_BAD_JSON = -40,               //!< The DCGM GPU Diagnostic returned JSON that cannot be parsed
    IXDCGM_RET_DIAG_BAD_LAUNCH = -41,             //!< Error while launching the DCGM GPU Diagnostic
    IXDCGM_RET_DIAG_UNUSED = -42,                 //!< Unused
    IXDCGM_RET_DIAG_THRESHOLD_EXCEEDED = -43,     //!< A field value met or exceeded the error threshold.
    IXDCGM_RET_INSUFFICIENT_DRIVER_VERSION = -44, //!< The installed driver version is insufficient for this API
    IXDCGM_RET_INSTANCE_NOT_FOUND = -45,          //!< The specified GPU instance does not exist
    IXDCGM_RET_COMPUTE_INSTANCE_NOT_FOUND = -46,  //!< The specified GPU compute instance does not exist
    IXDCGM_RET_CHILD_NOT_KILLED = -47,            //!< Couldn't kill a child process within the retries
    IXDCGM_RET_3RD_PARTY_LIBRARY_ERROR = -48,     //!< Detected an error in a 3rd-party library
    IXDCGM_RET_INSUFFICIENT_RESOURCES = -49,      //!< Not enough resources available
    IXDCGM_RET_PLUGIN_EXCEPTION = -50,            //!< Exception thrown from a diagnostic plugin
    IXDCGM_RET_IXVS_ISOLATE_ERROR = -51,          //!< The diagnostic returned an error that indicates the need for isolation
    IXDCGM_RET_IXVS_BINARY_NOT_FOUND = -52,       //!< The NVVS binary was not found in the specified location
    IXDCGM_RET_IXVS_KILLED = -53,                 //!< The NVVS process was killed by a signal
    IXDCGM_RET_PAUSED = -54,                      //!< The hostengine and all modules are paused
    IXDCGM_RET_ALREADY_INITIALIZED = -55,         //!< The object is already initialized
} ixdcgmReturn_t;

typedef uintptr_t ixdcgmHandle_t; //!< Identifier for ixDCGM Handle

#endif // end of __IXDCGM_STRUCTS_H__