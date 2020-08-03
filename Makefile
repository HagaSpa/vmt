all: bin/main

PLATFORM=local

.PHONY: bin/main

bin/main:
	@docker build . --target bin \
	--output bin/ \
	--platform ${PLATFORM}
