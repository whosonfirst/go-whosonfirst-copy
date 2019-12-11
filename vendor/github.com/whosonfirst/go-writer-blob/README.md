# go-writer-blob

Work in progress. This will replace [go-whosonfirst-readwrite](https://github.com/whosonfirst/go-whosonfirst-readwrite).

## Example

```
package main

import (
	"bytes"
	"errors"
	"context"
	"github.com/whosonfirst/go-writer"
	_ "github.com/aaronland/go-cloud-s3blob"	
	blob_writer "github.com/whosonfirst/go-writer-blob"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"	
	gc_blob "gocloud.dev/blob"	
	"io/ioutil"
)

func main() {

	bucket := "s3-bucket"
	region := "s3-region"
	
	blob_uri := fmt.Sprintf("s3blob://%s?region=%s&credentials=session", bucket, region)
	
	ctx := context.Background()
	wr, _ := writer.NewWriter(ctx, blob_uri)
	
	br := bytes.NewReader([]byte("hello"))
	fh := ioutil.NopCloser(br)

	before := func(asFunc func(interface{}) bool) error {
		
		req := &s3manager.UploadInput{}
		ok := asFunc(&req)
		
		if !ok {
			return errors.New("invalid s3 type")
		}
		
		req.ACL = aws.String("public-read")
		return nil
	}
	
	wr_opts := &gc_blob.WriterOptions{
		BeforeWrite: before,
	}

	ctx = context.WithValue(ctx, blob_writer.BlobWriterOptionsKey("options"), wr_opts)
	
	wr.Write(ctx, "test.txt", fh)
}
```

## See also

* https://github.com/whosonfirst/go-writer
* https://gocloud.dev/howto/blob/