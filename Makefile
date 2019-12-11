lambda: lambda-copy

lambda-copy:
	if test -f main; then rm -f main; fi
	if test -f copy.zip; then rm -f copy.zip; fi
	GOOS=linux go build -mod vendor -o main cmd/wof-copy/main.go
	zip copy.zip main
	rm -f main
