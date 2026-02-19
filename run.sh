#! /bin/sh

# helper script to build and test without compose
docker build -t calunite .
docker run --rm -p 8080:8080 -v ./testdata:/config -e LOG_LEVEL=debug calunite
