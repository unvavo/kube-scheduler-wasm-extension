package main

import (
	protoapi "sigs.k8s.io/kube-scheduler-wasm-extension/kubernetes/proto/api"
)

// These are valid values for TolerationOperator
const (
	TolerationOpExists string = "Exists"
	TolerationOpEqual  string = "Equal"
)

// ToleratesTaint checks if the toleration tolerates the taint.
// The matching follows the rules below:
//
//  1. Empty toleration.effect means to match all taint effects,
//     otherwise taint effect must equal to toleration.effect.
//  2. If toleration.operator is 'Exists', it means to match all taint values.
//  3. Empty toleration.key means to match all taint keys.
//     If toleration.key is empty, toleration.operator must be 'Exists';
//     this combination means to match all taint values and all taint keys.
func toleratesTaint(toleration *protoapi.Toleration, taint *protoapi.Taint) bool {
	if len(*toleration.Effect) > 0 && toleration.Effect != taint.Effect {
		return false
	}

	if len(*toleration.Key) > 0 && toleration.Key != taint.Key {
		return false
	}

	// TODO: Use proper defaulting when Toleration becomes a field of PodSpec
	switch *toleration.Operator {
	// empty operator means Equal
	case "", TolerationOpEqual:
		return toleration.Value == taint.Value
	case TolerationOpExists:
		return true
	default:
		return false
	}
}

// TolerationsTolerateTaint checks if taint is tolerated by any of the tolerations.
func tolerationsTolerateTaint(tolerations []*protoapi.Toleration, taint *protoapi.Taint) bool {
	for i := range tolerations {
		if toleratesTaint(tolerations[i], taint) {
			return true
		}
	}
	return false
}
