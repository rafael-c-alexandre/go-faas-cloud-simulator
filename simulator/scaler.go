package main

type Scaler struct {
	cluster      Cluster
	scaleMinLoad float32
}

func (s *Scaler) getRemainingAvailableMemory(id string) int {
	remainingAvailableMemory := 0

	for _, instance := range s.cluster.instances {
		if instance.id != id {
			remainingAvailableMemory += instance.currentAvailableMemory
		}
	}

	return remainingAvailableMemory

}

func (s *Scaler) ScanCluster() map[string]*functionInvocation {

	var orphanInvocations map[string]*functionInvocation

	remainingMemory := s.getRemainingAvailableMemory("")
	for _, instance := range s.cluster.instances {
		if float32(instance.currentAvailableMemory)/float32(instance.memory) >= (1 - s.scaleMinLoad) {
			instanceUsedMemory := instance.memory - instance.currentAvailableMemory

			if s.getRemainingAvailableMemory(instance.id) >= instanceUsedMemory &&
				instanceUsedMemory <= remainingMemory {

				orphanInvocations = s.ScaleDown(instance)
				remainingMemory -= instanceUsedMemory

			}
		}
	}

	return orphanInvocations
}

func (s *Scaler) ScaleUp() {
	newInstance := NewInstance()
	s.cluster.instances[newInstance.id] = newInstance
}

func (s *Scaler) ScaleDown(instance *Instance) map[string]*functionInvocation {
	orphanInvocations := s.cluster.DeleteInstance(instance.id)
	return orphanInvocations
}
