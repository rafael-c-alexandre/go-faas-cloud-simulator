package main

import "sync"

type Scheduler struct {
	cluster *Cluster
}

func (s *Scheduler) RouteInvocation(invocation *functionInvocation, scaler *Scaler, globalLock *sync.RWMutex) {
	globalLock.Lock()
	defer globalLock.Unlock()
	chosenInstance, err := s.getSuitableInstance(invocation.profile.AvgMemory)

	if err != nil || chosenInstance.currentAvailableMemory < invocation.profile.AvgMemory {
		chosenInstance = scaler.ScaleUp()
	}

	chosenInstance.RunNewFunction(invocation.id, invocation)
}

func (s *Scheduler) getSuitableInstance(memory int64) (*Instance, error) {
	suitableInstance, err := s.cluster.GetOne()

	if err != nil {
		return &Instance{}, err
	}

	currentSuitableMemory := suitableInstance.memory

	for _, instance := range s.cluster.instances {
		if instance.currentAvailableMemory > memory && instance.currentAvailableMemory < currentSuitableMemory {
			currentSuitableMemory = instance.currentAvailableMemory
			suitableInstance = instance
		}
	}
	return suitableInstance, nil
}
