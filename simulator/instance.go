package main

import (
	"log"
)

const INSTANCE_MEMORY = 32 * 1024 * 1024 * 1024 // 32 gb memory

type Instance struct {
	id                     string
	memory                 int
	currentAvailableMemory int
	functionsRunning       map[string]functionInvocation
}

func NewInstance() Instance {
	i := Instance{
		id:                     RandStringBytes(5),
		memory:                 INSTANCE_MEMORY,
		currentAvailableMemory: INSTANCE_MEMORY,
		functionsRunning:       map[string]functionInvocation{},
	}

	return i
}

func (i *Instance) RunNewFunction(id string, invocation functionInvocation) {
	i.functionsRunning[id] = invocation
	i.currentAvailableMemory -= invocation.profile.AvgMemory
}

func (i *Instance) UpdateStatus(step int) {
	for _, invocation := range i.functionsRunning {
		invocation.remainingTime -= step

		if invocation.remainingTime <= 0 {
			log.Printf("Instance %s: Invocation %s finished\n", i.id, invocation.id)
			delete(i.functionsRunning, invocation.id)
		}
	}
}
