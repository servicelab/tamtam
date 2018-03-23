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
	"sync"

	"github.com/eelcocramer/tamtam/service"
	"github.com/clockworksoul/smudge"
)

var monchannels = struct {
	sync.RWMutex
	m map[interface{}]chan<- *service.Node
}{m: make(map[interface{}]chan<- *service.Node)}

type statusListener struct {
	smudge.StatusListener
}

func (s statusListener) OnChange(node *smudge.Node, status smudge.NodeStatus) {
	n := SmudgeToTamTamNode(node)
	monchannels.Lock()
	for _, ch := range monchannels.m {
		ch <- n
	}
	monchannels.Unlock()
}

// AddMonitorChannel adds a monitor channel
func AddMonitorChannel(id interface{}, ch chan<- *service.Node) {
	monchannels.Lock()
	monchannels.m[id] = ch
	monchannels.Unlock()
}

// RemoveMonitorChannel removes a monitor channel
func RemoveMonitorChannel(id interface{}) {
	monchannels.Lock()
	delete(monchannels.m, id)
	monchannels.Unlock()
}

func init() {
	smudge.AddStatusListener(statusListener{})
}
