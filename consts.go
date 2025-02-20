/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package holmes

import (
	"os"
	"sync"
	"time"
)

const (
	defaultThreadTriggerMin  = 10 // 10 threads
	defaultThreadTriggerAbs  = 70 // 70 threads
	defaultThreadTriggerDiff = 25 // 25%

	defaultCPUTriggerMin   = 10              // 10%
	defaultCPUTriggerAbs   = 70              // 70%
	defaultCPUTriggerDiff  = 25              // 25%
	defaultCPUSamplingTime = 5 * time.Second // collect 5s cpu profile

	defaultGoroutineTriggerMin  = 3000   // 3000 goroutines
	defaultGoroutineTriggerAbs  = 200000 // 200k goroutines
	defaultGoroutineTriggerDiff = 20     // 20%  diff

	defaultMemTriggerMin  = 10 // 10%
	defaultMemTriggerAbs  = 80 // 80%
	defaultMemTriggerDiff = 25 // 25%

	defaultGCHeapTriggerMin  = 10 // 10%
	defaultGCHeapTriggerAbs  = 40 // 40%
	defaultGCHeapTriggerDiff = 20 // 20%

	defaultCooldown          = time.Minute
	defaultThreadCoolDown    = time.Hour
	defaultGoroutineCoolDown = time.Minute * 10

	defaultInterval        = 5 * time.Second
	defaultDumpProfileType = binaryDump
	defaultDumpPath        = "/tmp"
	defaultLoggerName      = "holmes.log"
	defaultLoggerFlags     = os.O_RDWR | os.O_CREATE | os.O_APPEND
	defaultLoggerPerm      = 0644
	defaultShardLoggerSize = 5242880 // 5m
)

type dumpProfileType int

const (
	binaryDump dumpProfileType = 0
	textDump   dumpProfileType = 1
)

type configureType int

const (
	mem configureType = iota
	cpu
	thread
	goroutine
	gcHeap
)

// check type to profile name, just align to pprof
var type2name = map[configureType]string{
	mem:       "heap",
	cpu:       "cpu",
	thread:    "threadcreate",
	goroutine: "goroutine",
	gcHeap:    "heap",
}

// check type to check name
var check2name = map[configureType]string{
	mem:       "mem",
	cpu:       "cpu",
	thread:    "thread",
	goroutine: "goroutine",
	gcHeap:    "GCHeap",
}

const (
	cgroupMemLimitPath  = "/sys/fs/cgroup/memory/memory.limit_in_bytes"
	cgroupCpuQuotaPath  = "/sys/fs/cgroup/cpu/cpu.cfs_quota_us"
	cgroupCpuPeriodPath = "/sys/fs/cgroup/cpu/cpu.cfs_period_us"
)

const (
	cgroupMemLimitPathV2 = "/sys/fs/cgroup/memory.max"
	cgroupCpuMaxPathV2   = "/sys/fs/cgroup/cpu.max"
)

const minCollectCyclesBeforeDumpStart = 10

const (
	// TrimResultTopN trimResult return only reserve the top n.
	TrimResultTopN = 10

	// TrimResultMaxBytes trimResultFront return only reserve the front n bytes.
	TrimResultMaxBytes = 512000

	// NotSupportTypeMaxConfig means this profile type is
	// not support control dump profile by max parameter.
	NotSupportTypeMaxConfig = 0

	// UniformLogFormat is the format of uniform logging.
	UniformLogFormat = "[Holmes] %v %v, config_min : %v, config_diff : %v, config_abs : %v, config_max : %v, previous : %v, current: %v"
)

var (
	checkMode sync.Once
	cgMode    CGMode
)

const unifiedMountpoint = "/sys/fs/cgroup"

// CGMode is the cgroups mode of the host system
type CGMode int

const (
	// Unavailable cgroup mountpoint
	Unavailable CGMode = iota
	// Legacy cgroups v1
	Legacy
	// Hybrid with cgroups v1 and v2 controllers mounted
	Hybrid
	// Unified with only cgroups v2 mounted
	Unified
)
