package main

func WithJob(j Job) StageOpt {
	return func(s IStage) {
		s.SetJob(j)
	}
}

func WithNext(event, id int) StageOpt {
	return func(s IStage) {
		s.SetNext(event, id)
	}
}
