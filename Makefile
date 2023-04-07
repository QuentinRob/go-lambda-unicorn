.PHONY: clean build local sync
.SILENT: build sync local clean

HANDLERS=$$(dir ./handler)
NETWORK=mysql

BG_GREEN=\033[0;42m
BG_RED=\033[0;41m
BLUE=\033[0;34m
END=\033[0m

.DEFAULT_GOAL := build

sync:
	echo "${BLUE} Synchronizing workspace...${END}"
	go work sync \
		&& echo "\t${BG_GREEN}Success${END} workspace synchronized" \
		|| echo "\t${BG_RED}Failure${END} workspace synchronized"

build: sync
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
