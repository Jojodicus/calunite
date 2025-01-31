#! /bin/sh

# I'm too poor for Docker Pro and automatic builds
# Also, a few links in the readme have to be changed for Docker Hub

cat README.md \
| sed -e 's=https://hub.docker.com/r/jojodicus/calunite=https://github.com/Jojodicus/calunite=g' \
| sed -e 's=(config.yml)=(https://github.com/Jojodicus/calunite/blob/main/config.yml)=g'
