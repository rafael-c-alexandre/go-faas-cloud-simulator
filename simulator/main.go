package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

const N_NODES = 100
const NODE_MEMORY = 64000
const N_THREADS = 8
const UNLOAD_POLICY = "random"
const MAX_DATASET_SIZE = 1000
const RESOURCES_PATH = "./resources"

type Statistics struct {
	invocations    [N_NODES]int
	coldStarts     [N_NODES]int
	failed         [N_NODES]int
	minuteProgress [MINUTES_IN_DAY + 1]int
	minutesLock    sync.Mutex
}

// This function adds the average duration of a function to the invocation count structure
func addDurations(functionInvocations []functionProfile, durations []functionExecutionDuration) []functionProfile {
	for i := range functionInvocations {
		functionInvocations[i].AvgDuration = -1
		for j := range durations {
			if functionInvocations[i].Function == durations[j].function {
				functionInvocations[i].AvgDuration = durations[j].average
				break
			}
		}
	}
	for i := range functionInvocations {
		if functionInvocations[i].AvgDuration == -1 || functionInvocations[i].AvgDuration < 1000 {
			RemoveFromList(functionInvocations, i)
		}
	}
	return functionInvocations
}

// This function adds the average memory of a function to the invocation count structure
func addMemories(functionInvocations []functionProfile, memoryUsages []appMemory) []functionProfile {
	for i := range functionInvocations {
		functionInvocations[i].AvgMemory = -1
		for j := range memoryUsages {
			if functionInvocations[i].App == memoryUsages[j].app {
				functionInvocations[i].AvgMemory = memoryUsages[j].average
				break
			}
		}
	}
	for i := range functionInvocations {
		if functionInvocations[i].AvgMemory == -1 {
			RemoveFromList(functionInvocations, i)
		}
	}
	return functionInvocations
}

func estimateRelevantInvocations(listInvocations []functionProfile) {

	relevantInvocations1, relevantInvocations10, relevantInvocations20, relevantInvocations30,
		relevantInvocations40, relevantInvocations50, relevantInvocations60,
		allInvocations := 0, 0, 0, 0, 0, 0, 0, 0

	for i := 0; i < len(listInvocations); i++ {
		for _, minuteCardinality := range listInvocations[i].PerMinute {
			if listInvocations[i].AvgDuration > 1000 {
				relevantInvocations1 += minuteCardinality
			}
			if listInvocations[i].AvgDuration > 10000 {
				relevantInvocations10 += minuteCardinality
			}
			if listInvocations[i].AvgDuration > 20000 {
				relevantInvocations20 += minuteCardinality
			}
			if listInvocations[i].AvgDuration > 30000 {
				relevantInvocations30 += minuteCardinality
			}
			if listInvocations[i].AvgDuration > 40000 {
				relevantInvocations40 += minuteCardinality
			}
			if listInvocations[i].AvgDuration > 50000 {
				relevantInvocations50 += minuteCardinality
			}
			if listInvocations[i].AvgDuration > 60000 {
				relevantInvocations60 += minuteCardinality
			}
			allInvocations += minuteCardinality
		}
	}
	log.Println()
	log.Println("-------------- Dataset Statistics -------------")
	log.Printf("Number of total invocations: %d\n", allInvocations)
	log.Printf("Fraction of invocations with > 1 sec: %.3f\n", float32(relevantInvocations1)/float32(allInvocations))
	log.Printf("Fraction of invocations with > 10 sec: %.3f\n", float32(relevantInvocations10)/float32(allInvocations))
	log.Printf("Fraction of invocations with > 20 sec: %.3f\n", float32(relevantInvocations20)/float32(allInvocations))
	log.Printf("Fraction of invocations with > 30 sec: %.3f\n", float32(relevantInvocations30)/float32(allInvocations))
	log.Printf("Fraction of invocations with > 40 sec: %.3f\n", float32(relevantInvocations40)/float32(allInvocations))
	log.Printf("Fraction of invocations with > 50 sec: %.3f\n", float32(relevantInvocations50)/float32(allInvocations))
	log.Printf("Fraction of invocations with > 60 sec: %.3f\n", float32(relevantInvocations60)/float32(allInvocations))
}

func prepareSimulation() {
	//Read the csv files into structure arrays
	log.Println("Reading the invocations per function files")
	var listInvocations []functionProfile
	for i := 5; i < 6; i++ {
		listInvocations = append(listInvocations,
			readInvocationCsvFile(fmt.Sprintf("dataset/invocations_per_function_md.anon.d0%d.csv", i))...)
	}

	log.Println("Reading the app memory files")
	var listMemory []appMemory
	for i := 5; i < 6; i++ {
		listMemory = append(listMemory,
			readAppMemoryCsvFile(fmt.Sprintf("dataset/app_memory_percentiles.anon.d0%d.csv", i))...)
	}

	log.Println("Reading the function duration files")
	var functionDuration []functionExecutionDuration
	for i := 5; i < 6; i++ {
		functionDuration = append(functionDuration,
			readFunctionDurationCsvFile(fmt.Sprintf("dataset/function_durations_percentiles.anon.d0%d.csv", i))...)
	}

	//Add the durations and memory to the invocation structure, so we have everything in the same array
	log.Println("Joining the average durations to each function")
	listInvocations = addDurations(listInvocations, functionDuration)
	log.Println("Joining the average memory usage to each function")
	listInvocations = addMemories(listInvocations, listMemory)

	estimateRelevantInvocations(listInvocations)

	// Create folder for serialized dataset object
	err := os.MkdirAll(RESOURCES_PATH, os.ModePerm)
	if err != nil {
		log.Fatalf("Caught error: %s", err)
	}

	var fileBuffer bytes.Buffer

	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register([]functionProfile{})

	// Create an encoder interface and send values
	enc := gob.NewEncoder(&fileBuffer)
	interfaceEncode(enc, listInvocations)

	file, err := os.OpenFile(fmt.Sprintf("%s/serialized_dataset", RESOURCES_PATH), os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("Caught error: %s", err)
	}

	_, err = file.Write(fileBuffer.Bytes())
	if err != nil {
		log.Fatalf("Caught error: %s", err)
	}

	if err = file.Close(); err != nil {
		log.Fatalf("Caught error: %s", err)
	}

}
func run() {

	fileData, err := os.ReadFile(fmt.Sprintf("%s/serialized_dataset", RESOURCES_PATH))
	if err != nil {
		log.Fatalf("Caught error: %s", err)
	}

	var listInvocations []functionProfile

	// Create a decoder interface and receive values
	dec := gob.NewDecoder(bytes.NewBuffer(fileData))
	listInvocations = interfaceDecode(dec, listInvocations)

}

func main() {

	//Measure the execution time
	//timeStart := time.Now()

	//Initialize statistics struct
	//stats := new(Statistics)

	var operation string

	flag.StringVar(&operation, "operation", "operation", "Operation name (Options: prepare, run)")
	flag.Parse()

	switch operation {
	case "prepare":
		log.Println("Operation: Prepare")
		prepareSimulation()
	case "run":
		log.Println("Operation: Run")
		run()
	default:
		log.Fatal("Missing operation argument.")
	}

	log.Println("Done")

	/*
		fmt.Printf("The simulation took %s\n", timeElapsed)
		fmt.Printf("Keep Alive: %d\n", KEEP_ALIVE_WINDOW)
		fmt.Printf("Invocations: %d\n", invocationsSum)
		fmt.Printf("Failed Invocations: %d\n", failedInvocationsSum)
		fmt.Printf("Cold Starts: %d\n", coldSum)*/
}
