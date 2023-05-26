package main

import "sync"

type Scaler struct {
	cluster      *Cluster
	scaleMinLoad float32
}

func (s *Scaler) ScanCluster(globalLock *sync.RWMutex) map[string]*functionInvocation {
	globalLock.Lock()
	defer globalLock.Unlock()
	orphanInvocations := map[string]*functionInvocation{}

	for _, instance := range s.cluster.instances {
		if float32(instance.currentAvailableMemory)/float32(instance.memory) >= (1 - s.scaleMinLoad) {
			instanceUsedMemory := instance.memory - instance.currentAvailableMemory
			if s.cluster.getRemainingAvailableMemory(instance.id) >= instanceUsedMemory &&
				len(s.cluster.instances) > 1 {
				orphanInvocations = MergeMaps(orphanInvocations, s.ScaleDown(instance))

			}
		}
	}

	return orphanInvocations
}

func (s *Scaler) ScaleUp() *Instance {
	newInstance := NewInstance()
	s.cluster.instances[newInstance.id] = newInstance
	return newInstance
}

func (s *Scaler) ScaleDown(instance *Instance) map[string]*functionInvocation {
	orphanInvocations := s.cluster.DeleteInstance(instance.id)
	return orphanInvocations
}
