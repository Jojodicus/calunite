#! /bin/sh

# helper script to build and test without compose
docker build -t calunite .
docker run -p 8080:8080 -v ./testdata:/config calunite
