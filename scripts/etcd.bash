function main() {
  export NODE1="$(ip route get 1 | sed -nE 's/.* src ([0-9.]+).*/\1/p')"
  export DATA_DIR="etcd-data"

  ETCD_VERSION=v3.5.19
  REGISTRY=gcr.io/etcd-development/etcd

  docker volume create --name etcd-data

  docker run \
    -p 2379:2379 \
    -p 2380:2380 \
    --volume=${DATA_DIR}:/etcd-data \
    --name etcd ${REGISTRY}:${ETCD_VERSION} \
    /usr/local/bin/etcd \
    --data-dir=/etcd-data --name node1 \
    --initial-advertise-peer-urls http://${NODE1}:2380 --listen-peer-urls http://0.0.0.0:2380 \
    --advertise-client-urls http://${NODE1}:2379 --listen-client-urls http://0.0.0.0:2379 \
    --initial-cluster node1=http://${NODE1}:2380

    docker network connect
}

main "$@"