package link

import (
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type Fetcher struct {
	client *http.Client
}

func New() *Fetcher {
	f := &Fetcher{
		client: &http.Client{
			Timeout: 6 * time.Second,
		},
	}

	return f
}

func (f *Fetcher) FetchLinks(query string, number int) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://yandex.ru/search/?text="+url.QueryEscape(query), nil)
	if err != nil {
		return nil, errors.Wrap(err, "build req fetchLinks")
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://ya.ru/")

	slog.Debug("Final request URL:", slog.String("url", req.URL.String()))

	resp, err := f.client.Do(req)
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
