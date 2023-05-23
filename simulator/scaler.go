package main

type Scaler struct {
	cluster      Cluster
	scaleMinLoad float32
}

func (s *Scaler) ScanCluster() map[string]functionInvocation {

	var orphanInvocations map[string]functionInvocation

	remainingMemory := s.cluster.getRemainingAvailableMemory("")

	for _, instance := range s.cluster.instances {
		if float32(instance.currentAvailableMemory)/float32(instance.memory) >= (1 - s.scaleMinLoad) {
			instanceUsedMemory := instance.memory - instance.currentAvailableMemory

			if s.cluster.getRemainingAvailableMemory(instance.id) >= instanceUsedMemory &&
				instanceUsedMemory <= remainingMemory {

				orphanInvocations = s.ScaleDown(instance)
				remainingMemory -= instanceUsedMemory

			}
		}
	}

	return orphanInvocations
}

func (s *Scaler) ScaleUp() Instance {
	newInstance := NewInstance()
	s.cluster.instances[newInstance.id] = newInstance
	return newInstance
}

func (s *Scaler) ScaleDown(instance Instance) map[string]functionInvocation {
	orphanInvocations := s.cluster.DeleteInstance(instance.id)
	return orphanInvocations
}
