package worker

import "log"

type Job struct {
	ID string
}

func Start(workers int, jobs <-chan Job) {
	for i := 0; i < workers; i++ {
		go func(idx int) {
			for job := range jobs {
				log.Printf("worker %d user id %s", idx, job.ID)
			}
		}(i)
	}
}
