#!/usr/bin/env bash
rm -rf ./data/redis/
docker-compose down
docker-compose up -d

#for i in 1 2 3 4 5 6; do
#    ip="$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' node_$i)"
#    #eval "var_$i=$ip"
#    declare var_$i=$ip
#done
#
#docker exec -it node_1  redis-cli --cluster create ${var_1}:8001 ${var_2}:8002 ${var_3}:8003 ${var_4}:8004 ${var_5}:8005 ${var_6}:8006 --cluster-replicas 1