# go-whosonfirst-copy

This is work in progress. Things will probably change, still.

## Example

### Command-line

```
go run -mod vendor cmd/wof-copy/main.go \
	-reader-uri 'github://sfomuseum-data/sfomuseum-data-media?prefix=data' \
	-writer-uri 'null://' \
	-writer-uri 'file:///usr/local/data' \
	-writer-uri 'elasticsearch://localhost/data?port=9999' \	
	151/159/924/5/1511599245.geojson
```

## See also

* https://github.com/whosonfirst/go-reader
* https://github.com/whosonfirst/go-writer
* https://github.com/whosonfirst/go-copy
