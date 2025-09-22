package worker

import (
	"fmt"
	"log"
	"net/url"

	"github.com/Kyrbanali/API-GateWay/internal/repository"
	"github.com/pkg/errors"
)

type Job struct {
	ID string
}

func LinkWorker(workers int, jobs <-chan Job, userRepo repository.UserProvider) {
	for i := 0; i < workers; i++ {
		go func(idx int) {
			for job := range jobs {
				user, err := userRepo.GetUserByID(nil, job.ID)
				if err != nil {
					errors.Wrap(err, "LinkWorker")
					continue
				}

				q := url.QueryEscape(fmt.Sprintf("%s %d", user.Name, user.Age))

				links := []string{
					fmt.Sprintf("https://yandex.ru/search/?text=%s&p=1", q),
					fmt.Sprintf("https://yandex.ru/search/?text=%s&p=2", q),
					fmt.Sprintf("https://yandex.ru/search/?text=%s&p=3", q),
				}

				log.Printf("worker %d user (%s, %d) links %v", idx, user.Name, user.Age, links)
			}
		}(i)
	}
}
