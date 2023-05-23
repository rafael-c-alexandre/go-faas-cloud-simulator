package main

import "errors"

type Cluster struct {
	instances map[string]Instance
}

func (c *Cluster) AddInstance(instance Instance) {
	c.instances[instance.id] = instance
}

func (c *Cluster) DeleteInstance(id string) map[string]functionInvocation {
	orphanInvocations := c.instances[id].functionsRunning
	delete(c.instances, id)

	return orphanInvocations
}

func (c *Cluster) UpdateStatus() {
	for _, instance := range c.instances {
		instance.UpdateStatus(1000)
	}
}

func (c *Cluster) GetOne() (Instance, error) {
	var i string
	if len(c.instances) == 0 {
		return Instance{}, errors.New("no instances available")
	}
	for s, _ := range c.instances {
		i = s
		break
	}
	return c.instances[i], nil
}

func (c *Cluster) getRemainingAvailableMemory(id string) int {
	remainingAvailableMemory := 0

	for _, instance := range c.instances {
		if instance.id != id {
			remainingAvailableMemory += instance.currentAvailableMemory
		}
	}

	return remainingAvailableMemory

}
