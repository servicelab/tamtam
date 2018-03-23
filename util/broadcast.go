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

	"github.com/clockworksoul/smudge"
)

var channels = struct {
	sync.RWMutex
	m map[interface{}]chan<- []byte
}{m: make(map[interface{}]chan<- []byte)}

type broadcastListener struct {
	smudge.BroadcastListener
}

func (m broadcastListener) OnBroadcast(b *smudge.Broadcast) {
	channels.Lock()
	for _, ch := range channels.m {
		ch <- b.Bytes()
	}
	channels.Unlock()
}

// AddBroadcastChannel adds a broadcast channel
func AddBroadcastChannel(id interface{}, ch chan<- []byte) {
	channels.Lock()
	channels.m[id] = ch
	channels.Unlock()
}

// RemoveBroadcastChannel removes a broadcast channel
func RemoveBroadcastChannel(id interface{}) {
	channels.Lock()
	delete(channels.m, id)
	channels.Unlock()
}

func init() {
	smudge.AddBroadcastListener(broadcastListener{})
}
