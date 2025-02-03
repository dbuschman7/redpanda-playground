#! /bin/sh 
set -ex

# Make the tmp dirs for databases
mkdir -p tmp/backend

#
# Watch the sqlite database dirs 
# watch 'find tmp -exec ls -l {} \;' 
#
# create the topics 
# rpk topic create -p 3 -r 3 cluster1_topic_demo cluster2_udp_topic cluster1_topic1 cluster2_topic1

BACKEND="$(   ls backend-*.yaml    | tr '\n' ' ')"
FRONTEND="$(  ls frontend-*.yaml   | tr '\n' ' ')"
GENERATOR="$( ls generator-*.yaml  | tr '\n' ' ')"

#STREAMS = $BACKEND $FRONTEND $GENERATOR
#STREAMS = $GENERATOR
# STREAMS="$BACKEND $FRONTEND"
#
# Example - simple straight through from generator to topic
# ##############################################################
# STREAMS="generator-to-topic.yaml"

#
# Example - simple straight through from generator to topic via inproc
# ##############################################################
# STREAMS="generator-to-demo-inproc.yaml backend-to-demo-topic.yaml"

#
# Example - straight through from generator to topic via brokered dual inproc backends
#STREAMS="backend-buffer-topic-1.yaml backend-buffer-topic-2.yaml generator-to-brokered-inproc.yaml "

# All of the above
STREAMS="$BACKEND $FRONTEND $GENERATOR"

echo $STREAMS
UDP_PORT=1514 SERVICE_NAME=forward ADMIN_PORT=4195 BASE_DIR=$(pwd) rpk connect streams -o server.yaml $STREAMS
