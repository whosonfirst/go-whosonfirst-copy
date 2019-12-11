package main

/*

go run -mod vendor cmd/wof-copy/main.go -reader-uri 'github://sfomuseum-data/sfomuseum-data-media?prefix=data' -writer-uri 'elasticsearch://localhost?port=9200&index=wof' -writer-uri 'null://' 151/159/924/5/1511599245.geojson

*/

import (
	_ "github.com/aaronland/go-cloud-s3blob"
	_ "gocloud.dev/blob/fileblob"
)

import (
	"bytes"
	"context"
	"flag"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/whosonfirst/go-reader"
	_ "github.com/whosonfirst/go-reader-github"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-writer"
	_ "github.com/whosonfirst/go-writer-blob"
	"github.com/whosonfirst/go-writer-elasticsearch"
	"io/ioutil"
	"log"
	"strings"
)

// basically all of the Copier / Copy stuff is already in whosonfirst/go-copy
// but I haven't figured out how to define a custom URI function for wrangling
// a fully-qualified URI in to an ID for the ES writer (20191211/thisisaaronland)

type Copier struct {
	Reader  reader.Reader
	Writers []writer.Writer
}

func CopyMany(ctx context.Context, cp *Copier, uris ...string) error {

	// do these in parallel (see notes above wrt/ go-copy)	

	for _, uri := range uris {

		err := Copy(ctx, cp, uri)

		if err != nil {
			return err
		}
	}

	return nil
}

func Copy(ctx context.Context, cp *Copier, uri string) error {

	fh, err := cp.Reader.Read(ctx, uri)

	if err != nil {
		return err
	}

	defer fh.Close()

	f, err := feature.LoadGeoJSONFeatureFromReader(fh)

	if err != nil {
		return err
	}

	// do these in parallel (see notes above wrt/ go-copy)

	for _, wr := range cp.Writers {

		br := bytes.NewReader(f.Bytes())
		fh := ioutil.NopCloser(br)

		var wr_uri string

		switch wr.(type) {
		case *elasticsearch.ElasticsearchWriter:
			wr_uri = f.Id()
		default:
			wr_uri = uri
		}

		err := wr.Write(ctx, wr_uri, fh)

		if err != nil {
			return err
		}

		log.Println("WROTE", wr_uri, f.Id())
	}

	return nil
}

func main() {

	reader_uri := flag.String("reader-uri", "", "A valid go-reader Reader URI.")
	mode := flag.String("mode", "cli", "Valid options are: cli, lambda.")

	var writer_uris flags.MultiString
	flag.Var(&writer_uris, "writer-uri", "A valid go-writer Writer URI.")

	flag.Parse()

	err := flags.SetFlagsFromEnvVars("WOF_COPY")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reader, err := reader.NewReader(ctx, *reader_uri)

	if err != nil {
		log.Fatal(err)
	}

	var writers []writer.Writer

	if *mode == "lambda" {

		uris := make([]string, 0)

		for _, wr_uri := range writer_uris {

			for _, u := range strings.Split(wr_uri, ",") {
				uris = append(uris, u)
			}
		}

		writer_uris = uris
	}

	for _, wr_uri := range writer_uris {

		wr, err := writer.NewWriter(ctx, wr_uri)

		if err != nil {
			log.Fatal(err)
		}

		writers = append(writers, wr)
	}

	cp := &Copier{
		Reader:  reader,
		Writers: writers,
	}

	switch *mode {
	case "cli":

		uris := flag.Args()
		err := CopyMany(ctx, cp, uris...)

		if err != nil {
			log.Fatal(err)
		}

	case "invoke":

		// as in: invoke a lambda function (below) with a list of URIs

		log.Fatal("Please write me")

	case "lambda":

		// specifically, this assumes: https://github.com/whosonfirst/go-webhookd/#githubcommits
		
		handler := func(ctx context.Context, args []interface{}) error {

			uris := make([]string, len(args))

			for i, uri := range args {
				uris[i] = uri.(string)
			}
			
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			err = CopyMany(ctx, cp, uris...)

			if err != nil {
				return err
			}

			return nil
		}

		aws_lambda.Start(handler)

	default:
		log.Fatal("Unsupported mode")
	}

}
