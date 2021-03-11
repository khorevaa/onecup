package jobs

type Runner struct {
	params Params
	jobs   []Job
}

func (r *Runner) Run(params Params) error {

	globalInput := params

	for _, j := range r.jobs {

		output, err := j.run(Input(globalInput))

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
	params Params
	jobs   []Job
}

func (r *SimulateRunner) Run(params Params) error {

	globalInput := params

	for _, j := range r.jobs {

		output, err := j.simulate(Input(globalInput))

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
