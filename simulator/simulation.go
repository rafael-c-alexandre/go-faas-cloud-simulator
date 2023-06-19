package main

import (
	"github.com/schollz/progressbar/v3"
	"log"
	"math/rand"
	"sync"
	"time"
)

const SIMULATION_MINUTES = 1444
const SIMULATION_DURATION = SIMULATION_MINUTES * 60

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
	totalResources     int
}

func (s *Simulation) schedule(minute int) [60][]string {
	var schedule [60][]string

	for _, profile := range s.profiles {
		nInvocations := profile.PerMinute[minute]
		j := 0
		for j < nInvocations {
			invocationTs := rand.Intn(59)
			schedule[invocationTs] = append(schedule[invocationTs], profile.Id)
			j++
		}
	}

	return schedule
}

func (s *Simulation) Run() {

	s.Statistics.timeStart = time.Now()

	s.lock = sync.RWMutex{}

	epoch := 0
	minute := 0

	// Schedule simulation
	invocationSchedule := s.schedule(minute)

	// Launch one instance for the start of the simulation
	s.cluster.AddInstance(NewInstance(epoch), &s.lock)

	bar := progressbar.Default(SIMULATION_DURATION)

	for i := 0; i < SIMULATION_MINUTES; i += 1 {
		for _, second := range invocationSchedule {
			epoch++
			var wg sync.WaitGroup
			wg.Add(3)
			go s.updateStatus(&wg)
			go s.scanCluster(epoch, &wg)
			go s.newRound(second, epoch, &wg)
			wg.Wait()
			bar.Add(1)
		}
		minute++
		invocationSchedule = s.schedule(minute)
	}

	for _, instance := range s.cluster.instances {
		s.cluster.totalResources += SIMULATION_DURATION - instance.launchTs
	}

	s.Statistics.totalResources = s.cluster.totalResources

}

func (s *Simulation) updateStatus(wg *sync.WaitGroup) {
	s.cluster.UpdateStatus(&s.lock)
	wg.Done()
}

func (s *Simulation) scanCluster(now int, wg *sync.WaitGroup) {

	// Scaler check for scaling down the cluster
	orphans := s.scaler.ScanCluster(now, &s.lock)
	//log.Printf("Cluster scanned. Captured %d orphan invocations. Re-routing.\n", len(orphans))
	if orphans != nil {
		for _, invocation := range orphans {
			s.scheduler.RouteInvocation(invocation, s.scaler, now, &s.lock)
			s.Statistics.evictedInvocations += 1
		}
	}
	wg.Done()
}

func (s *Simulation) newRound(second []string, now int, wg *sync.WaitGroup) {
	// New period invocations
	for _, invocation := range second {
		newInvocation := NewInvocation(s.profiles[invocation])
		s.scheduler.RouteInvocation(newInvocation, s.scaler, now, &s.lock)
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
	log.Printf("Total Instance-seconds: %d\n", stats.totalResources)

}
