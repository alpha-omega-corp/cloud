build:
	docker build --tag redis -f Dockerfile .

run:
	docker run --network=cloud --name cache --publish 6379:6379 redis --save 60 1 --loglevel warning

uptrace:
	uptrace --config=uptrace.yml serve

uptrace-reset:
	uptrace --config=uptrace.yml ch reset