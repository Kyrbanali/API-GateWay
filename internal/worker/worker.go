package worker

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type Task func(ctx context.Context) error

type Worker struct {
	jobs   chan Task
	client *http.Client
}

func New(workers int) *Worker {
	w := &Worker{
		jobs: make(chan Task, 100),
		client: &http.Client{
			Timeout: 6 * time.Second,
		},
	}

	for i := 0; i < workers; i++ {
		go w.worker()
	}
	return w
}

func (w *Worker) worker() {
	for job := range w.jobs {
		if err := job(context.Background()); err != nil {
			slog.Error("worker task", slog.String("error", err.Error()))
		}
	}
}

func (w *Worker) Push(ctx context.Context, task Task) {
	go func() {
		select {
		case w.jobs <- task:
		case <-ctx.Done():
		}
	}()
}

func (w *Worker) FetchLinks(query string, number int) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://yandex.ru/search/?text="+url.QueryEscape(query), nil)
	if err != nil {
		return nil, errors.Wrap(err, "build req fetchLinks")
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://ya.ru/")

	log.Println("Final request URL:", req.URL.String())
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http do fetchLinks")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Debug("close response body", slog.String("error", err.Error()))
		}
	}()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "parse html fetchLinks")
	}

	var links []string

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if len(links) >= number {
			return
		}

		if href, ok := s.Attr("href"); ok && href != "" {
			links = append(links, href)
		}

	})

	return links, nil
}
