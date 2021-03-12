package jobs

type Runner struct {
	params Params
	jobs   []Job
}

func (r *Runner) Run(params Values) error {

	globalInput := params

	for _, j := range r.jobs {

		output, err := j.run(globalInput)

		if err != nil {
			return err
		}

		for key, val := range output {
			globalInput[key] = val
		}

	}

	return nil

}

type SimulateRunner struct {
	params Values
	jobs   []Job
}

func (r *SimulateRunner) Run(params Values) error {

	globalInput := params
	simulateSub := NewSubscriber(AllEvents)

	for _, j := range r.jobs {

		j.Subscribe(simulateSub)

		output, err := j.simulate(globalInput)

		if err != nil {
			return err
		}

		for key, val := range output {
			globalInput[key] = val
		}

	}

	return nil

}

func NewRunner(jobs ...Job) *Runner {
	return &Runner{
		params: map[string]interface{}{},
		jobs:   jobs,
	}
}

func NewSimulateRunner(jobs ...Job) *SimulateRunner {
	return &SimulateRunner{
		params: map[string]interface{}{},
		jobs:   jobs,
	}
}
