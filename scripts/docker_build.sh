#!/bin/bash
TAG="0.13"
#docker buildx build --push --platform linux/arm/v7,linux/arm64/v8,linux/amd64 --tag lyw1217/scrapingnewsgo:latest "../"
docker build --tag lyw1217/scrapingnewsgo:${TAG} "../"
docker push lyw1217/scrapingnewsgo:${TAG}
