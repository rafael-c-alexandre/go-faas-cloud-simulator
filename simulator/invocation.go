package main

type functionInvocation struct {
	id            string
	profile       functionProfile
	remainingTime int
}

func NewInvocation(profile functionProfile) functionInvocation {
	return functionInvocation{
		id:            RandStringBytes(7),
		profile:       profile,
		remainingTime: profile.AvgDuration,
	}
}
