package main

import (
	"math/rand"
	"time"
)

type Simulation struct {
	timeStart   time.Time
	timeElapsed time.Time
	scheduler   Scheduler
	scaler      Scaler
	profiles    map[string]functionProfile
}

func (s *Simulation) schedule() [1440][60][]string {
	var schedule [1440][60][]string

	for _, profile := range s.profiles {
		for i, nInvocations := range profile.PerMinute {
			j := 0
			for j < nInvocations {
				invocationTs := rand.Intn(59)
				schedule[i][invocationTs] = append(schedule[i][invocationTs], profile.Id)
				j++
			}
		}
	}
	return schedule
}

func (s *Simulation) Start() {
	s.timeStart = time.Now()

	c := Cluster{instances: map[string]Instance{}}

	s.scheduler = Scheduler{cluster: c}
	s.scaler = Scaler{
		cluster:      c,
		scaleMinLoad: 0.3,
	}

	// Schedule simulation
	invocationSchedule := s.schedule()
	// Launch one instance for the start of the simulation
	c.AddInstance(NewInstance())

	for _, minute := range invocationSchedule {
		for _, second := range minute {

			c.UpdateStatus()
			orphans := s.scaler.ScanCluster()

			if orphans != nil {
				for _, invocation := range orphans {
					s.scheduler.RouteInvocation(invocation, s.scaler)
				}
			}

			for _, invocation := range second {
				newInvocation := NewInvocation(s.profiles[invocation])
				s.scheduler.RouteInvocation(newInvocation, s.scaler)
			}
		}
	}

}
