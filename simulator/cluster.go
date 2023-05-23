package main

type Cluster struct {
	instances map[string]*Instance
}

func (c *Cluster) AddInstance(instance *Instance) {
	c.instances[instance.id] = instance
}

func (c *Cluster) DeleteInstance(id string) map[string]*functionInvocation {
	orphanInvocations := c.instances[id].functionsRunning
	delete(c.instances, id)

	return orphanInvocations
}
