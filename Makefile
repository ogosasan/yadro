
build:
	go build -o xkcd yadro/cmd/xkcd/

run:
	xkcd -c configs/config.yaml