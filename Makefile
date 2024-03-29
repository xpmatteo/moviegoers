

test:
	go test ./...

watch:
	@air -v > /dev/null 2> /dev/null || go install github.com/cosmtrek/air@latest
	@air

open:
	open http://localhost:8080

M=github.com/xpmatteo/gomovies
depgraph:
	godepgraph -s -onlyprefixes $(M) $(M) \
		| sed -e s@$(M)/@@g -e s@$(M)@main@g -e s/splines=ortho/splines=curved/ \
		| dot -Tpng -o /tmp/godepgraph.png \
		&& open /tmp/godepgraph.png

staticcheck:
	@staticcheck -version > /dev/null 2> /dev/null || go install honnef.co/go/tools/cmd/staticcheck@latest
	@staticcheck ./...

# heroku

heroku-remote:
	open https://moviegoers-ef258a681da3.herokuapp.com/

heroku-deploy:
	git push heroku main

heroku-turnoff:
	heroku ps:scale web=0

heroku-turnon:
	heroku ps:scale web=1

heroku-logs:
	heroku logs --tail

