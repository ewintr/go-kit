# this Makefile purpose is to help testing all packages, consolidate coverage
# report, exame Go source and ensure format.
SRC = $(shell find . -type f -name '*.go' | \
		awk -F'__' '{ sub ("/[^/]*$$", "/", $$1); print $1 }' | sort | uniq)

PACKAGES = log slugify herror test

all: dep fmt vet test_all

dep:
	@for pkg in $(PACKAGES); do \
		echo "- Checking dependencies for $$pkg"; \
		cd $$pkg && go get && cd ..; \
	done

fmt:
	@echo "- Checking code format"
	@GO_FMT=$$(gofmt -e -l ${SRC}) && \
		if [ -n "$$GO_FMT" ]; then \
		  	echo '$@: Incorrect format has been detected in your code run `make fmt-fix`'; \
			exit 1; \
		fi

fmt-fix:
	@echo "- Checking code format"
	@for file in $$(go fmt ${SRC}) ; do \
 		echo "$@: $$file fixed and staged"; \
		git add "./${file}"; \
	done

vet:
	@for pkg in $(PACKAGES); do \
		echo "- Examine source code for $$pkg"; \
		cd $$pkg && go vet . && cd ..; \
	done

test_all:
	@rm -f ./coverage/*.out ./coverage/*.html
	@for pkg in $(PACKAGES); do \
		echo "- Testing package $$pkg"; \
		go test ./$$pkg -coverprofile=./coverage/$$pkg.cover.out; \
	done
	@echo "- Merging coverage output files"
	@echo "mode: set" > ./coverage/coverage.out && \
		cat ./coverage/*.cover.out | grep -v mode: | sort -r | \
		awk -f ./coverage/merge.awk >> ./coverage/coverage.out
	@go tool cover -html=./coverage/coverage.out \
		-o ./coverage/coverage.html
	@go tool cover --func=./coverage/coverage.out | \
		awk -f ./coverage/total_coverage.awk
