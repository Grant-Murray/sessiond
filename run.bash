#!/bin/bash

# to generate a random key use:
# hexdump -n64 -v -e '1/1 "%.2x"' < /dev/urandom

DBSOURCE="user=glm1 password='with spaces' host=localhost port=5432 sslmode=disable"
SERVERKEY=c3efc2d599734b0b7488514f5bf393a7869324c0d65f27cb9d569b466d86d03f3974091e52a044a0b40c30eecf74709f3a881da3e17a9deef46ad7bd72fe308a

echo -e "$DBSOURCE\n$SERVERKEY" | go run session.go


