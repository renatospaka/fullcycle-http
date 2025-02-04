.PHONY: 
	go.mod 
	run 
	test
	tidy 
	build 
	net 
	kill 

#########################################################
##### RUNNING

setup.go.environment:
	go env -w GOPRIVATE=github.com/renatospaka/*;
	go env -w GOPROXY=direct;
	go env -w GOSUMDB=off;

run: 
	$(MAKE) setup.go.environment
	cd cmd/server && \
	go run . && \
	cd ../../

run.test:
	$(MAKE) --no-print-directory run

#########################################################
##### CLEANNING AND TESTING

clean-cache:
	go clean -cache
	
test:
	go env -w CGO_ENABLED=1;
	clear && \
	go test ./... -cover -race;

test.clean:
	clear && \
	go clean -cache && \
	go test ./... -cover

#########################################################
##### GO MOD AND STUFF

tidy: 
	go mod tidy
	
go.mod: 
	$(MAKE) setup.go.environment && \
	$(MAKE) tidy;

go.get.t:
	@go get -t -u ./... 
	$(MAKE) go.mod

go.get.all:
	@go get -u all 
	$(MAKE) go.mod

net:
	lsof -t -i:8080;

kill:
	@kill -9 $(shell lsof -t -i:8080) || true
