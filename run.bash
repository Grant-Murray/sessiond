#!/bin/bash

# This is a convenience script for starting the sessiond server. It should only be used in isolated test environments.
# In production environments the SERVERKEY should be stored securely and entered on the terminal at startup.

# to generate a random server key use: 
#     hexdump -n64 -v -e '1/1 "%.2x"' < /dev/urandom
#   or
#     openssl rand -hex 64

cd $GOPATH/src/code.grantmurray.com/sessiond
LOGDIR="/tmp/sessiond.glog"

killall --quiet sessiond
mkdir -p $LOGDIR
rm "$LOGDIR"/sessiond.*

go clean
go install || exit 1

# execute the session-pg-schema.sql
cd $GOPATH/src/code.grantmurray.com/session
PSQL="psql --username=postgres --dbname=sessdb"
$PSQL -c 'DROP SCHEMA IF EXISTS session CASCADE;'
$PSQL -c 'DROP ROLE IF EXISTS sessionr;'
$PSQL -f session-pg-schema.sql
$PSQL -f test.config.sql


DBSOURCE="user=sessionr password='big-secret' dbname=sessdb host=localhost port=5432 sslmode=disable"
SERVERKEY=fa1725ba8034485170912d8c29d4ef118f3fddd43e21437f0ee167835921b786d4bc6f52027fb858e6a138d6dfa1875d4ec12488464af3dbe79984bc23ffdece

echo -e "$DBSOURCE\n$SERVERKEY" | sessiond -instance="sessiond-A" -stderrthreshold=FATAL -log_dir="$LOGDIR" -v=2 &
echo -e "$DBSOURCE\n$SERVERKEY" | sessiond -instance="sessiond-B" -stderrthreshold=FATAL -log_dir="$LOGDIR" -v=2 &
echo -e "$DBSOURCE\n$SERVERKEY" | sessiond -instance="sessiond-C" -stderrthreshold=FATAL -log_dir="$LOGDIR" -v=2 &
echo -e "$DBSOURCE\n$SERVERKEY" | sessiond -instance="sessiond-D" -stderrthreshold=FATAL -log_dir="$LOGDIR" -v=2 &
echo -e "$DBSOURCE\n$SERVERKEY" | sessiond -instance="sessiond-E" -stderrthreshold=FATAL -log_dir="$LOGDIR" -v=2 &

sleep 1 # cleans up the terminal output
echo "PID of sessiond: $(pidof sessiond)"
