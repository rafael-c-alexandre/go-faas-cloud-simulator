package main

const INSTANCE_MEMORY = 32 * 1024 * 1024 * 1024 // 32 gb memory

type Instance struct {
	id                     string
	memory                 int64
	currentAvailableMemory int64
	functionsRunning       map[string]*functionInvocation
	launchTs               int
}

func NewInstance(start int) *Instance {
	i := Instance{
		id:                     RandStringBytes(5),
		memory:                 INSTANCE_MEMORY,
		currentAvailableMemory: INSTANCE_MEMORY,
		functionsRunning:       map[string]*functionInvocation{},
		launchTs:               start,
	}

	return &i
}

func (i *Instance) RunNewFunction(id string, invocation *functionInvocation) {
	i.functionsRunning[id] = invocation
	i.currentAvailableMemory -= invocation.profile.AvgMemory
}

func (i *Instance) UpdateStatus(step int) bool {

	for _, invocation := range i.functionsRunning {
		invocation.remainingTime -= step

		if invocation.remainingTime <= 0 {
			delete(i.functionsRunning, invocation.id)
			i.currentAvailableMemory += invocation.profile.AvgMemory
		}
	}
	return len(i.functionsRunning) > 0
}
