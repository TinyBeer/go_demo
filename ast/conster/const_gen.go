// Code generated by conster. DO NOT EDIT.

package main

type PlanType int 
var ( 
    PlanType_Daily PlanType  = 0
    PlanType_Weekly PlanType  = 1
    PlanType_Monthly PlanType  = 2
)
func (v PlanType)String() string {
    switch v { 
    case 0: 
      return "PlanType_Daily"
    case 1: 
      return "PlanType_Weekly"
    case 2: 
      return "PlanType_Monthly"
    default:
      return ""
    }
}

type TodoType int 
var ( 
    TodoType_Times TodoType  = 0
    TodoType_Duration TodoType  = 1
    TodoType_TimesAndDuration TodoType  = 2
)
func (v TodoType)String() string {
    switch v { 
    case 0: 
      return "TodoType_Times"
    case 1: 
      return "TodoType_Duration"
    case 2: 
      return "TodoType_TimesAndDuration"
    default:
      return ""
    }
}
