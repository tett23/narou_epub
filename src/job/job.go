package job

type Job struct {
	JobType JobType
	NCode   string
}

type JobType int

const (
	JobTypeNone JobType = iota
	JobTypeFetchLatestEpisode
	JobTypeFetchAll
	JobTypeBuildLatestEpisode
	JobTypeBuildAll
	JobTypeSendToKindleLatest
	JobTypeSendToKindleAll
)

var queue chan *Job

func init() {
	queue = make(chan *Job, 1)
}

func Enqueue(t JobType, nCode string) {
	queue <- &Job{
		JobType: t,
		NCode:   nCode,
	}
}

func processJob(job *Job) error {
	var err error
	switch job.JobType {
	case JobTypeFetchLatestEpisode:
		err = fetchLatestEpisode(job)
	case JobTypeFetchAll:
		err = fetchAll(job)
	case JobTypeBuildLatestEpisode:
		err = buildLatestEpisode(job)
	case JobTypeBuildAll:
		err = buildAll(job)
	case JobTypeSendToKindleLatest:
		err = sendToKindleLatest(job)
	case JobTypeSendToKindleAll:
		err = sendToKindleAll(job)
	case JobTypeNone:
		fallthrough
	default:
		return nil
	}

	return err
}

func fetchLatestEpisode(job *Job) error {
	return nil
}

func fetchAll(job *Job) error {
	return nil
}

func buildLatestEpisode(job *Job) error {
	return nil
}

func buildAll(job *Job) error {
	return nil
}

func sendToKindleLatest(job *Job) error {
	return nil
}

func sendToKindleAll(job *Job) error {
	return nil
}
