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

### Lambda

If run in "lambda" mode the code will assume it's being sent messages that have been produced by the [go-webhookd GitHubCommits transformation](https://github.com/whosonfirst/go-webhookd/#githubcommits). 

| Key | Value | Notes |
| --- | --- | --- |
| WOF_COPY_MODE | lambda | | 
| WOF_COPY_READER_URI | github://sfomuseum-data/%s?prefix=data | Any valid `go-reader` URI. |
| WOF_COPY_WRITER_URI | null:// | A comma-separated list of one or more valid `go-writer` URIs. |

Notes:

* See the way the (`go-reader`) reader URI contains a `%s` placeholder? When the Lambda handler is invoked it will parse the message looking for a GitHub repository to replace the placeholder with and create a `Reader` instance with.
* You can define multiple (`go-writer`) writer URIs in the `WOF_COPY_WRITER_URI` environment variable by passing them in as a comma-separated list.

`GitHubCommits` messages look like this:

```
[
"a66e0a6b934eb2cef1d92effb2aee1875534d454,sfomuseum-data-media,151/159/924/5/1511599245.geojson"
]
```

So assuming the code is passed the message above you'd see something like this in the logs:

```
2019/12/11 23:48:44 WROTE 151/159/924/5/1511599245.geojson 1511599245
```

## See also

* https://github.com/whosonfirst/go-reader
* https://github.com/whosonfirst/go-writer
* https://github.com/whosonfirst/go-copy
