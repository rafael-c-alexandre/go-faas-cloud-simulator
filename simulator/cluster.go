package main

import (
	"errors"
	"sync"
)

type Cluster struct {
	instances map[string]*Instance
}

func (c *Cluster) AddInstance(instance *Instance, globalLock *sync.RWMutex) {
	c.instances[instance.id] = instance
	//log.Printf("Adding instance. Number of current active instances %d\n", len(c.instances))
}

func (c *Cluster) DeleteInstance(id string) map[string]*functionInvocation {
	orphanInvocations := c.instances[id].functionsRunning
	delete(c.instances, id)
	//log.Printf("Deleting instance. Number of current active instances %d\n", len(c.instances))

	return orphanInvocations
}

func (c *Cluster) UpdateStatus(globalLock *sync.RWMutex) {
	globalLock.Lock()
	defer globalLock.Unlock()
	for _, instance := range c.instances {
		instance.keepAlive += 1
		instance.UpdateStatus(1000)
	}
}

func (c *Cluster) GetOne() (*Instance, error) {
	var i string

	if len(c.instances) == 0 {
		return &Instance{}, errors.New("no instances available")
	}
	for s, _ := range c.instances {
		i = s
		break
	}
	return c.instances[i], nil
}

func (c *Cluster) getRemainingAvailableMemory(id string) int64 {
	remainingAvailableMemory := int64(0)

	for _, instance := range c.instances {
		if instance.id != id {
			remainingAvailableMemory += instance.currentAvailableMemory
		}
	}

	return remainingAvailableMemory

}

func (c *Cluster) all() int {
	return len(c.instances)
}
