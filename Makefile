BINARIES := code_injector_dependencies remark-inject-code

build:
	@for target in $(BINARIES); do \
		echo Installing $$target ; \
		go install ./cmd/$$target ; \
	done

all: build
