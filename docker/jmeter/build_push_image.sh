#!/bin/bash

docker build . -t tsuyoshiushio/jmeter:latest -t tsuyoshiushio/jmeter:0.2
docker push tsuyoshiushio/jmeter
