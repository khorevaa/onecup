package jobs

type List []Job

func (l *List) Add(j Job) {

	*l = append(*l, j)

}
func (l *List) Append(job ...Job) {
	for _, j := range job {
		*l = append(*l, j)
	}
}
