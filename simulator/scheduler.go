package main

type Scheduler struct {
	cluster Cluster
}

func (s *Scheduler) RouteInvocation(invocation *functionInvocation) {
	chosenInstance := s.getSuitableInstance(invocation.profile.AvgDuration)

	if chosenInstance.currentAvailableMemory > invocation.profile.AvgMemory {
		chosenInstance.RunNewFunction(invocation.id, invocation)
	} else {
	}
}

func (s *Scheduler) getSuitableInstance(duration int) *Instance {
	var suitableInstance *Instance
	currentSuitableMemory := INSTANCE_MEMORY

	for _, instance := range s.cluster.instances {
		if instance.currentAvailableMemory > duration && instance.currentAvailableMemory < currentSuitableMemory {
			currentSuitableMemory = instance.currentAvailableMemory
			suitableInstance = instance
		}
	}
	return suitableInstance
}
