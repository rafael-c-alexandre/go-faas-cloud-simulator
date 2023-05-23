package main

import "fmt"

const INSTANCE_MEMORY = 32 * 1024 * 1024 * 1024 // 32 gb memory

type Instance struct {
	id                     string
	memory                 int
	currentAvailableMemory int
	functionsRunning       map[string]*functionInvocation
}

func NewInstance() *Instance {
	i := new(Instance)
	i.id = RandStringBytes(5)
	i.memory = INSTANCE_MEMORY
	i.currentAvailableMemory = i.memory
	i.functionsRunning = map[string]*functionInvocation{}

	return i
}

func (i *Instance) RunNewFunction(id string, invocation *functionInvocation) {
	i.functionsRunning[id] = invocation
	i.currentAvailableMemory -= invocation.profile.AvgMemory
}

func (i *Instance) UpdateStatus(step int) {
	for _, invocation := range i.functionsRunning {
		invocation.remainingTime -= step

		if invocation.remainingTime <= 0 {
			fmt.Printf("Instance %s: Invocation %s finished\n", i.id, invocation.id)
			delete(i.functionsRunning, invocation.id)
		}
	}
}
