#!/bin/bash

sudo chown ubuntu:ubuntu /var/lib/docker/volumes/nitro_poster-data/_data/nitro.log
sudo chmod 755 /var/lib/docker/volumes/nitro_poster-data/_data/nitro.log

go run .
