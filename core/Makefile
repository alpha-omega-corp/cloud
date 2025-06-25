bootstrap:
	docker compose up -d && bash ./scripts/etcd.bash

otel:
	docker compose up -d && bash ./scripts/uptrace.bash

uptrace:
	uptrace --config=uptrace.yml serve

uptrace-reset:
	uptrace --config=uptrace.yml ch reset

status:
	etcdctl --endpoints=http://localhost:2379 member list