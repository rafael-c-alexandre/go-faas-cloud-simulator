package main

import "time"

type Simulation struct {
	timeStart   time.Time
	timeElapsed time.Time
	scheduler   Scheduler
	scaler      Scaler
	profiles    []*functionProfile
}

func (s *Simulation) Start() {
	s.timeStart = time.Now()

	c := Cluster{instances: map[string]*Instance{}}

	s.scheduler = Scheduler{cluster: c}
	s.scaler = Scaler{
		cluster:      c,
		scaleMinLoad: 0.3,
	}

	// Schedule simulation

	for i := 1; i < 1441; i += 1 {
		for j := 0; j < 60; j += 1 {

			orphans := s.scaler.ScanCluster()

			if orphans != nil {
				for _, invocation := range orphans {
					s.scheduler.RouteInvocation(invocation)
				}
			}

		}
	}

}
