//go:build conster
// +build conster

package main

// PlanType
const (
	PlanType_Daily = iota
	PlanType_Weekly
	PlanType_Monthly
)

// TodoType
const (
	TodoType_Times = iota
	TodoType_Duration
	TodoType_TimesAndDuration
)
