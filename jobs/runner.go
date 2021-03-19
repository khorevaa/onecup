package jobs

import "context"

type Runner struct {
	params Values
	jobs   []Job
}

func (r *Runner) Run(ctx context.Context, params Values) error {

	globalInput := params

	for _, j := range r.jobs {

		output, err := j.Run(ctx, globalInput)

		if err != nil {
			return err
		}

		for key, val := range output {
			globalInput[key] = val
		}

	}

	return nil

}

func NewJobsRunner(jobs ...Job) *Runner {
	return &Runner{
		params: map[string]interface{}{},
		jobs:   jobs,
	}
}
