package main

import (
	"github.com/schollz/progressbar/v3"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Simulation struct {
	cluster    *Cluster
	scheduler  *Scheduler
	scaler     *Scaler
	profiles   map[string]functionProfile
	Statistics *Statistics
	lock       sync.RWMutex
}

type Statistics struct {
	timeStart          time.Time
	timeElapsed        time.Duration
	invocations        int
	evictedInvocations int
}

func (s *Simulation) schedule() [500][60][]string {
	var schedule [500][60][]string

	for _, profile := range s.profiles {
		for i, nInvocations := range profile.PerMinute[:500] {
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

func (s *Simulation) Run() {

	s.Statistics.timeStart = time.Now()

	s.lock = sync.RWMutex{}

	// Schedule simulation
	invocationSchedule := s.schedule()
	// Launch one instance for the start of the simulation
	s.cluster.AddInstance(NewInstance(), &s.lock)

	bar := progressbar.Default(500 * 60)

	for _, minute := range invocationSchedule {
		for _, second := range minute {

			var wg sync.WaitGroup
			wg.Add(3)

			go s.updateStatus(&wg)
			go s.scanCluster(&wg)
			go s.newRound(second, &wg)
			wg.Wait()
			bar.Add(1)
		}
	}

	// Pending invocations after the 24-hour period
	//for s.cluster.UpdateStatus() {
	//
	//	// Scaler check for scaling down the cluster
	//	orphans := s.scaler.ScanCluster()
	//	log.Printf("Cluster scanned. Captured %d orphan invocations. Re-routing.\n", len(orphans))
	//	if orphans != nil {
	//		for _, invocation := range orphans {
	//			s.scheduler.RouteInvocation(invocation, s.scaler)
	//			s.Statistics.evictedInvocations += 1
	//		}
	//	}
	//}

}

func (s *Simulation) updateStatus(wg *sync.WaitGroup) {
	s.cluster.UpdateStatus(&s.lock)
	wg.Done()
}

func (s *Simulation) scanCluster(wg *sync.WaitGroup) {

	// Scaler check for scaling down the cluster
	orphans := s.scaler.ScanCluster(&s.lock)
	//log.Printf("Cluster scanned. Captured %d orphan invocations. Re-routing.\n", len(orphans))
	if orphans != nil {
		for _, invocation := range orphans {
			s.scheduler.RouteInvocation(invocation, s.scaler, &s.lock)
			s.Statistics.evictedInvocations += 1
		}
	}
	wg.Done()
}

func (s *Simulation) newRound(second []string, wg *sync.WaitGroup) {
	// New period invocations
	for _, invocation := range second {
		newInvocation := NewInvocation(s.profiles[invocation])
		s.scheduler.RouteInvocation(newInvocation, s.scaler, &s.lock)
		s.Statistics.invocations += 1
	}
	wg.Done()
}

func (s *Simulation) Finalize() {
	s.Statistics.timeElapsed = time.Since(s.Statistics.timeStart)
}

func (stats *Statistics) Display() {
	log.Println("")
	log.Println("---------- Simulation stats ------------")
	log.Printf("Simulation duration: %.3fs\n", stats.timeElapsed.Seconds())
	log.Printf("Total invocations: %d\n", stats.invocations)
	log.Printf("Total evicted invocations: %d\n", stats.evictedInvocations)
	log.Printf("Percentage evicted invocations: %.3f\n", float32(stats.evictedInvocations)/float32(stats.invocations))

}
