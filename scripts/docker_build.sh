#!/bin/bash
docker buildx build --push --platform linux/arm/v7,linux/arm64/v8,linux/amd64 --tag lyw1217/scrapingnewsgo:latest "../"