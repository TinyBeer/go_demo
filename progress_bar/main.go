package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/fatih/color"
)

const max = 40

func main() {
	bar()
}

func simple() {
	bar := bytes.Repeat([]byte("-"), max)
	for i := 0; i <= max; i++ {
		if i == max {
			fmt.Printf("\rprogress: %3d%%  [%s]   √", 100, bar)
			return
		}
		bar[i] = '>' // #
		fmt.Printf("\rprogress: %3d%% [%s]", i*100/max, bar)
		time.Sleep(time.Millisecond * 200)
	}
}

func colorful() {
	bar := bytes.Repeat([]byte("-"), max)
	for i := 0; i <= max; i++ {
		if i == max {
			color.Green("\rprogress: %3d%%  [%s]   √", 100, bar)
			return
		}
		bar[i] = '>' // #
		fmt.Printf("\r%s%s%s",
			color.YellowString("progress: %3d%% [", i*100/max),
			color.GreenString("%s", bar[:i+1]),
			color.YellowString("%s]", bar[i+1:]),
		)
		time.Sleep(time.Millisecond * 200)
	}
}

func bar() {
	bar := bytes.Repeat([]byte("-"), max)
	i := 0
	spin := []rune{0x285f, 0x283f, 0x28bb, 0x28f9, 0x28fc, 0x28f6, 0x28e7, 0x28cf}
	exit := make(chan bool)

	go func() {
		v := 0
		for {
			time.Sleep(time.Millisecond * 200)
			select {
			case <-exit:
				return
			default:
				cur := i
				if cur < max {
					fmt.Printf("\r %s %s%s%s",
						color.HiMagentaString("%c", spin[v%len(spin)]),
						color.YellowString("progress: %3d%%  [", i*100/max),
						color.GreenString("%s", bar[:i+1]),
						color.YellowString("%s]", bar[i+1:]),
					)
				}
				v++
			}
		}
	}()
	for ; i <= max; i++ {
		if i == max {
			exit <- true
			color.Green("\r √  progress: %3d%%  [%s]", 100, bar)
			return
		}
		bar[i] = '>'
		time.Sleep(time.Millisecond * 200)
	}
}
