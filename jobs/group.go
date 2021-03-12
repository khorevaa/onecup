package jobs

func Group(name string, job ...TaskBuilder) Job {
	return &groupJob{
		name: name,
		jobs: job,
	}
}

type groupJob struct {
	name string
	jobs []TaskBuilder
}

func (g *groupJob) Subscribe(subscribe *Subscribe) {
	panic("implement me")
}

func (g *groupJob) simulate(input Values) (Values, error) {
	panic("implement me")
}

func (g *groupJob) Stats() Stats {
	panic("implement me")
}

func (g *groupJob) Fault() bool {
	panic("implement me")
}

func (g *groupJob) Skiped() bool {
	panic("implement me")
}

func (g *groupJob) Success() bool {
	panic("implement me")
}

func (g *groupJob) Status() CompletionStatus {
	panic("implement me")
}

func (g *groupJob) Error() error {
	panic("implement me")
}

func (g *groupJob) Name() string {
	return g.name
}

func (g *groupJob) run(input Values) (Values, error) {

	//tasks := g.getSteps()
	//
	//return tasks

	return map[string]interface{}{}, nil
}
