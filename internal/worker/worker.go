package worker

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Kyrbanali/API-GateWay/internal/repository"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type Job struct {
	ID string
}

func LinkWorker(workers int, jobs <-chan Job, userRepo repository.UserProvider) {
	client := &http.Client{Timeout: 6 * time.Second}

	for i := 0; i < workers; i++ {
		go func(idx int) {
			for job := range jobs {
				user, err := userRepo.GetUserByID(nil, job.ID)
				if err != nil {
					errors.Wrap(err, "LinkWorker")
					continue
				}

				links, err := getTop3(client, fmt.Sprintf("%s %d", user.Name, user.Age))
				if err != nil {
					errors.Wrap(err, "getTop3 LinkWorker")
					continue
				}

				log.Printf("worker %d user (%s, %d) links %v", idx, user.Name, user.Age, links)
			}
		}(i)
	}
}

func getTop3(client *http.Client, query string) ([]string, error) {
	u := "https://yandex.ru/search/?text=" + url.QueryEscape(query)

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "build req getTop3")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http do getTop3")
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "parse html getTop3")
	}

	var links []string

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if len(links) >= 3 {
			return
		}
		href, ok := s.Attr("href")
		if ok && href != "" {
			links = append(links, href)
		}

	})

	return links, nil
}
