package elasticsearch

import (
	"context"
	"fmt"
	es "github.com/elastic/go-elasticsearch"
	esapi "github.com/elastic/go-elasticsearch/esapi"
	wof_writer "github.com/whosonfirst/go-writer"
	"io"
	"net/url"
	"strconv"
)

type ElasticsearchWriter struct {
	wof_writer.Writer
	client *es.Client
	index  string
}

func init() {
	wr := NewElasticsearchWriter()
	wof_writer.Register("elasticsearch", wr)
}

func NewElasticsearchWriter() wof_writer.Writer {

	wr := ElasticsearchWriter{}
	return &wr
}

func (wr *ElasticsearchWriter) Open(ctx context.Context, uri string) error {

	u, err := url.Parse(uri)

	if err != nil {
		return err
	}

	host := u.Host
	index := u.Path

	port := 9200

	q := u.Query()

	str_port := q.Get("port")

	if str_port != "" {

		p, err := strconv.Atoi(str_port)

		if err != nil {
			return err
		}

		port = p
	}

	var es_endpoint string

	switch port {
	case 443:
		es_endpoint = fmt.Sprintf("https://%s", host)
	default:
		es_endpoint = fmt.Sprintf("http://%s:%d", host, port)
	}

	es_cfg := es.Config{
		Addresses: []string{es_endpoint},
	}

	client, err := es.NewClient(es_cfg)

	if err != nil {
		return err
	}

	wr.client = client
	wr.index = index

	return nil
}

func (wr *ElasticsearchWriter) Write(ctx context.Context, uri string, fh io.ReadCloser) error {

	req := esapi.IndexRequest{
		Index:      wr.index,
		DocumentID: uri,
		Body:       fh,
		Refresh:    "true",
	}

	_, err := req.Do(ctx, wr.client)

	if err != nil {
		return err
	}

	return nil
}
