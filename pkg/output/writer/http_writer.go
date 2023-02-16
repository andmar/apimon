package writer

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/ncarlier/apimon/pkg/model"
	"github.com/ncarlier/apimon/pkg/output/format"
)

// HTTPWriter HTTP writer
type HTTPWriter struct {
	Formatter format.Formatter
	URL       string
	Headers   []string
}

func newHTTPWriter(url string, formatter format.Formatter, headers []string) *HTTPWriter {
	return &HTTPWriter{
		Formatter: formatter,
		URL:       url,
		Headers:   headers,
	}
}

// Write post metric to HTTP endpoint
func (w *HTTPWriter) Write(metric model.Metric) error {

	contentType := w.Formatter.ContentType()
	body := w.Formatter.Format(metric)

	req, err := http.NewRequest("POST", w.URL, bytes.NewBufferString(body))
	if err != nil {
		return fmt.Errorf("could not create POST request to target endpoint")
	}

	for _, header := range w.Headers {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			req.Header.Add(parts[0], parts[1])
		}
	}

	req.Header.Add("Content-Type", contentType)

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; APImon/1.0; +https://github.com/ncarlier/apimon)")
	}

	// fmt.Println(req)
	// fmt.Println(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)

	// contentType := w.Formatter.ContentType()
	// body := w.Formatter.Format(metric)
	// resp, err := http.Post(w.URL, contentType, bytes.NewBufferString(body))

	//fmt.Println(body, resp)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	return nil
}

// Close close the metric writer
func (w *HTTPWriter) Close() error {
	return nil
}
