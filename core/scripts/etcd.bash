export VOLUME="config_etcd"
export NODE1="$(hostname -I | xargs)"

ETCD_VERSION=v3.4.37
REGISTRY=gcr.io/etcd-development/etcd

docker volume create --name ${VOLUME}
export DATA_DIR=${VOLUME}

docker run \
  -p 2379:2379 \
  -p 2380:2380 \
  --volume=${DATA_DIR}:/etcd-data \
  --name core-config ${REGISTRY}:${ETCD_VERSION} \
  /usr/local/bin/etcd \
  --data-dir=/etcd-data --name node1 \
  --initial-advertise-peer-urls http://${NODE1}:2380 \
  --listen-peer-urls http://0.0.0.0:2380 \
  --advertise-client-urls http://${NODE1}:2379 \
  --listen-client-urls http://0.0.0.0:2379 \
  --initial-cluster node1=http://${NODE1}:2380
