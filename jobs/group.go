package jobs

func NewGroup(job2 *job, name string, steps Steps, inputsOutputs ...Inputs) Task {

	var inputs, outputs Inputs

	if len(inputsOutputs) == 1 {
		inputs = inputsOutputs[0]
	}
	if len(inputsOutputs) == 2 {
		outputs = inputsOutputs[1]
	}

	return &groupTask{
		task{
			job:     job2,
			name:    name,
			inputs:  inputs,
			outputs: outputs,
			steps:   steps,
		},
	}

}

var _ Task = (*groupTask)(nil)

type groupTask struct {
	task
}
