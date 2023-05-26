package main

type functionProfile struct {
	Id          string
	Owner       string
	App         string
	Function    string
	Trigger     string
	PerMinute   [1440]int
	AvgDuration int
	AvgMemory   int64
}

type functionExecutionDuration struct {
	owner                string
	app                  string
	function             string
	average              int
	count                int
	minimum              int
	maximum              int
	percentileAverage0   int
	percentileAverage1   int
	percentileAverage25  int
	percentileAverage50  int
	percentileAverage75  int
	percentileAverage99  int
	percentileAverage100 int
}

type appMemory struct {
	owner                string
	app                  string
	count                int
	average              int64
	percentileAverage1   int64
	percentileAverage5   int64
	percentileAverage25  int64
	percentileAverage50  int64
	percentileAverage75  int64
	percentileAverage95  int64
	percentileAverage99  int64
	percentileAverage100 int64
}
