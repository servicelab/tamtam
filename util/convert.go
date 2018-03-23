/*
Copyright 2018, Eelco Cramer and the TamTam contributors

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

package util

import (
	"github.com/eelcocramer/tamtam/service"
	"github.com/clockworksoul/smudge"
)

// SmudgeToTamTamStatus converts a smudge node status to a TamTam status
func SmudgeToTamTamStatus(s smudge.NodeStatus) service.Status {
	switch s {
	case smudge.StatusAlive:
		return service.Status_ALIVE
	case smudge.StatusForwardTo:
		return service.Status_FORWARD_TO
	case smudge.StatusSuspected:
		return service.Status_SUSPECTED
	case smudge.StatusDead:
		return service.Status_DEAD
	default:
		return service.Status_UNKNOWN
	}
}

// SmudgeToTamTamNode converts a smudge node to a TamTam node
func SmudgeToTamTamNode(node *smudge.Node) *service.Node {
	return &service.Node{
		Address: &service.NodeAddress{
			IP:   node.IP().String(),
			Port: uint32(node.Port()),
		},
		Age:         node.Age(),
		EmitCounter: int32(node.EmitCounter()),
		PingMillis:  int32(node.PingMillis()),
		Timestamp:   node.Timestamp(),
		Status:      SmudgeToTamTamStatus(node.Status()),
	}
}

// StatusToString converts a TamTam status to a string
func StatusToString(s service.Status) string {
	switch s {
	case service.Status_DEAD:
		return "DEAD"
	case service.Status_ALIVE:
		return "ALIVE"
	case service.Status_FORWARD_TO:
		return "FOWARD_TO"
	case service.Status_SUSPECTED:
		return "SUSPECTED"
	default:
		return "UNKNOWN"
	}
}
