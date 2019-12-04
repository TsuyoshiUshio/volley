#!/bin/bash

docker build . -t tsuyoshiushio/jmeter:latest -t tsuyoshiushio/jmeter:0.1
docker push tsuyoshiushio/jmeter
