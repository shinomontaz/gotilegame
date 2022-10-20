// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"time"
)

func inform(e int) {
	switch e {
	case EVENT_DONE:
		fmt.Println("event done")
		setStage(currStage.GetNext(EVENT_DONE))
	case EVENT_ENTER:
		fmt.Println("event enter")
		setStage(currStage.GetNext(EVENT_ENTER))
	case EVENT_NOTREADY:
		fmt.Println("event not ready")
		loadingStage.SetUp(WithJob(currStage.Init), WithNext(EVENT_DONE, currStage.GetID()))
		setStage(loadingStage.GetID())
	}
}

func setStage(id int) {
	currStage = stages[id]
	currStage.Start()
}

var currStage IStage
var isquit bool
var stages map[int]IStage
var loadingStage IStage

func main() {
	stages = make(map[int]IStage, 0)
	loadingStage = NewStage2(1, inform)
	stages[1] = loadingStage
	stages[2] = NewStage2(2, inform)
	//	stages[2].SetUp(WithNext(EVENT_ENTER, 3))
	stages[3] = NewStage3(3, inform)

	currStage = stages[1]

	currStage.SetUp(WithJob(stages[2].Init), WithNext(EVENT_DONE, 2))
	currStage.Init()
	currStage.Start()

	last := time.Now()

	for !isquit {
		dt := time.Since(last).Seconds()
		currStage.Run(dt)
	}
}
