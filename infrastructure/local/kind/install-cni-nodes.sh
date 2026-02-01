#!/bin/bash

download_cni() {
    mkdir -p /tmp/cni-plugins
    cd /tmp/cni-plugins
    curl -L -o cni-plugins.tgz https://github.com/containernetworking/plugins/releases/download/v1.6.2/cni-plugins-linux-amd64-v1.6.2.tgz
    tar -xzf cni-plugins.tgz
}

copy_to_nodes() {
    for node in $(docker ps --filter "name=ng-voice-cluster" --format "{{.Names}}"); do
        docker exec $node mkdir -p /opt/cni/bin
        for plugin in /tmp/cni-plugins/*; do
            docker cp $plugin $node:/opt/cni/bin/$(basename $plugin)
            docker exec $node chmod +x /opt/cni/bin/$(basename $plugin)
        done
    done
}

download_cni
copy_to_nodes
