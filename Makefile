.PHONY: test watch run open e2e


test:
	go test ./...

watch:
	@air -v > /dev/null 2> /dev/null || go install github.com/cosmtrek/air@latest
	@air -build.exclude_dir e2e

server:
	go run server.go

open:
	open http://localhost:8080

M=github.com/xpmatteo/gomovies
depgraph:
	godepgraph -s -onlyprefixes $(M) $(M) \
		| sed -e s@$(M)/@@g -e s@$(M)@main@g -e s/splines=ortho/splines=curved/ \
		| dot -Tpng -o /tmp/godepgraph.png \
		&& open /tmp/godepgraph.png

