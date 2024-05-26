package main

import (
	"fmt"
	"runtime/metrics"
)

func main() {
	desc := metrics.All()
	fmt.Println("total:", len(desc))
	for _, m := range desc {
		fmt.Println(
			"\nName:", m.Name,
			// "\nKind:", m.Kind,
			// "\nDesc:", m.Description,
			// "\nCumulative:", m.Cumulative,
		)
	}
}
