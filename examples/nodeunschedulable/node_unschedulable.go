package main

import (
	guestapi "sigs.k8s.io/kube-scheduler-wasm-extension/guest/api"
	protoapi "sigs.k8s.io/kube-scheduler-wasm-extension/kubernetes/proto/api"
	"sigs.k8s.io/kube-scheduler-wasm-extension/guest/api/proto"
)

// NodeUnschedulable plugin filters nodes that set node.Spec.Unschedulable=true unless
// the pod tolerates {key=node.kubernetes.io/unschedulable, effect:NoSchedule} taint.
type NodeUnschedulable struct{}

// ErrReasonUnschedulable is used for NodeUnschedulable predicate error.
const ErrReasonUnschedulable = "node(s) were unschedulable"

func (n *NodeUnschedulable) Filter(_ guestapi.CycleState, pod proto.Pod, nodeInfo guestapi.NodeInfo) *guestapi.Status {
	node := nodeInfo.Node()

	if node.Spec().Unschedulable == nil || !*node.Spec().Unschedulable {
		return nil
	}

	// TaintNodeUnschedulable will be added when node becomes unschedulable
	// and removed when node becomes schedulable.
	taintNodeUnschedulable := "node.kubernetes.io/unschedulable"
	// TaintNodeMemoryPressure will be added when node has memory pressure
	// and removed when node has enough memory.
	taintEffectNoSchedule := "NoSchedule"

	// If pod tolerate unschedulable taint, it's also tolerate `node.Spec.Unschedulable`.
	podToleratesUnschedulable := tolerationsTolerateTaint(pod.Spec().Tolerations, &protoapi.Taint{
	// podToleratesUnschedulable := tolerationsTolerateTaint(pod.Spec().Tolerations(), &guestapi.TaintInfo().Taint(),{
		Key:    &taintNodeUnschedulable,
		Effect: &taintEffectNoSchedule,
	})

	if !podToleratesUnschedulable {
		return &guestapi.Status{Code: guestapi.StatusCodeUnschedulable, Reason: ErrReasonUnschedulable}
	}

	return nil
}
