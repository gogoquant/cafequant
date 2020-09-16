VERSION?="0.1.0"
TEST?=./...
TESTARGS ?= -v
CHECK_PACKAGE?=$$(go list ./... | grep -v migrate)
GOFMT_FILES?=$$(find . -type f -name '*.go')


default: check

clean:
	rm -rf golint-report.xml  coverage.out report.json govet-report.out

check: lint fmt
	#golangci-lint run  --config=.golangci.yml


#allcheck: fmtcheck vet lint staticcheck cyclo


staticcheck:
	staticcheck $(GOFMT_FILES)

cyclo:
	gocyclo -over 15 $(GOFMT_FILES)

lint:
	go list $(CHECK_PACKAGE) | xargs  golint -set_exit_status

fmt:
	gofmt -s -w -l $(GOFMT_FILES)

