.PHONY = build test

build-docker:
	docker build -t forum-tp .

run:
	docker run  --memory 2G --log-opt max-size=5M --log-opt max-file=3 -p 5000:5000 -p 5432:5432 --name forum-tp -t forum-tp

run-build: build-docker run


rm-docker:
	docker rm -vf $$(docker ps -a -q) || true
