/*
   Copyright 2023 The Kubernetes Authors.

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

// Package enqueue is defined internally so that it can export Pod as
// cyclestate.Pod, without circular dependencies or exporting it publicly.
package enqueue

import (
	"runtime"
	"unsafe"

	"sigs.k8s.io/kube-scheduler-wasm-extension/guest/api"
	"sigs.k8s.io/kube-scheduler-wasm-extension/guest/internal/plugin"
)

// enqueue is the current plugin assigned with SetPlugin.
var enqueue api.EnqueueExtensions

// SetPlugin is exposed to prevent package cycles.
func SetPlugin(enqueueExtensions api.EnqueueExtensions) {
	if enqueueExtensions == nil {
		panic("nil enqueueExtensions")
	}
	enqueue = enqueueExtensions
	plugin.MustSet(enqueueExtensions)
}

// prevent unused lint errors (lint is run with normal go).
var _ func() = _enqueue

// enqueue is only exported to the host.
//
//go:wasmexport enqueue
func _enqueue() {
	println("0")
	if enqueue == nil { // Then, the user didn't define one.
		// This is likely caused by use of plugin.Set(p), where 'p' didn't
		// implement EnqueueExtensions: return to use default events.
		return
	}

	println("1")

	clusterEvents := enqueue.EventsToRegister()
	println("2")

	// If plugin returned clusterEvents, encode them and call the host with the
	// count and memory region.
	encoded := encodeClusterEvents(clusterEvents)
	println("3")
	if encoded != nil {
		println("4")
		ptr := uint32(uintptr(unsafe.Pointer(&encoded[0])))
		println("5")
		size := uint32(len(encoded))
		println("6")
		setClusterEventsResult(ptr, size)
		println("7")
		runtime.KeepAlive(encoded) // until ptr is no longer needed.
		println("8")
	}
}
