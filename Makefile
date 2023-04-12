.PHONY: clean build local sync
.SILENT: build sync local clean check go-deps pre-check

HANDLERS=$$(dir ./handler)
NETWORK=mysql
GOLANG_VERSION = go1.19

BG_GREEN=\033[0;42m
BG_RED=\033[0;41m
BLUE=\033[0;34m
END=\033[0m

.DEFAULT_GOAL := build

pre-check:
	echo "${BLUE} Verifying requirements...${END}"

check: pre-check go-deps
	echo "\t${BG_GREEN}Success${END} All dependencies verified"


go-deps:
	which go > /dev/null || echo "${BG_RED}Missing${END} Required dependency \"golang\""
	if [ "$$(go version | cut -d" " -f3)" != "${GOLANG_VERSION}" ]; then \
		echo "${BG_RED}Failure${END} Version of golang doesn't match requirements\n\tInstalled $$(go version | cut -d" " -f3) require ${GOLANG_VERSION}"; \
		exit 1; \
  	fi


sync:
	echo "${BLUE} Synchronizing workspace...${END}"
	go work sync \
		&& echo "\t${BG_GREEN}Success${END} workspace synchronized" \
		|| echo "\t${BG_RED}Failure${END} workspace synchronized"

build: check sync
	echo "${BLUE} Building lambdas...${END}"
	for handler in $(HANDLERS); do \
  		go build -o dist/ ./handler/$$handler \
  			&& echo "\t${BG_GREEN}Success${END} build module \"./dist/$${handler}\"" \
  			|| echo "\t${BG_RED}Failure${END} build module \"./dist/$${handler}\""; \
  	done

clean:
	echo "${BLUE} Removing go cache...${END}"
	go clean -cache && echo "\t${BG_GREEN}Success${END} Cache cleaned" || echo "\t${BG_RED}Failure${END}"
	echo "${BLUE} Removing dist artifacts...${END}"
	rm -f ./dist/* && echo "\t${BG_GREEN}Success${END} Artifacts deleted" || echo "\t${BG_RED}Failure${END}"

local: build
	echo "${BLUE} Launching local environment...${END}"
	sam local start-api --docker-network $(NETWORK)
