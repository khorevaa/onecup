package jobs

type Runner struct {
	params Params
	jobs   []Job
}

func (r *Runner) Run() error {

	globalInput := r.params

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

func NewRunner(jobs ...Job) *Runner {
	return &Runner{
		params: map[string]interface{}{},
		jobs:   jobs,
	}
}
