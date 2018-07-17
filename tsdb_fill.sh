#!/bin/bash

curl -XPOST "http://localhost:8086/query" --data-urlencode "q=CREATE DATABASE test_db"

for i in $(seq 1 1000); do
	echo "next tag: host$i..."
	for j in $(seq 1 100); do
		let "t = 5 * 60 * 1000000000 * $j"
		curl -XPOST "http://localhost:8086/write?db=test_db" -d "cpu,host=host$i load=$(($RANDOM % 100)) $t"
	done
done
