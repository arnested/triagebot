.PHONY: doc test check-env deploy logs post-fixture

ENTRY_POINT=Handle
RUNTIME=go113
MEMORY=128M

export GO111MODULE=on

doc: README.md

README.md: *.go README.md.template
	go generate

test: *.go
	go test ./...

check-env:
	@test -n "$(GCF_NAME)" || (echo "Missing environment variable NAME" ; false)
	@test -n "$(GCF_PROJECT)" || (echo "Missing environment variable PROJECT" ; false)
	@test -n "$(GCF_REGION)" || (echo "Missing environment variable REGION" ; false)

deploy: test check-env
	gcloud functions deploy $(GCF_NAME) --project=$(GCF_PROJECT) --region=$(GCF_REGION) --entry-point=$(ENTRY_POINT) --runtime=$(RUNTIME) --trigger-http --memory=$(MEMORY)

logs: check-env
	gcloud functions logs read $(GCF_NAME) --project=$(GCF_PROJECT) --region=$(GCF_REGION) --limit=100
