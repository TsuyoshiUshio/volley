#!/bin/bash

base_name="volley"
platforms=("windows/amd64" "darwin/amd64" "linux/amd64" "windows/386" "darwin/386" "linux/386")

sudo apt-get install zip -y

rm -rf output

for platform in "${platforms[@]}"
do 
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_dir='output/'${base_name}'-'$GOOS'-'$GOARCH
    output_name="volley"
    if [ $GOOS = "windows" ]; then
      output_name="volley.exe"
    fi
    mkdir -p $output_dir
    cd $output_dir
    pwd
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name ../../pkg/cmd/main.go
    cd ..
    if [ $GOOS != "windows" ]; then
      tar -zcvf ${base_name}'-'$GOOS'-'$GOARCH.tgz ${base_name}'-'$GOOS'-'$GOARCH
      # On the pipeline, we should have zip compression for windows. 
    else
      zip ${base_name}'-'$GOOS'-'$GOARCH -r ${base_name}'-'$GOOS'-'$GOARCH
    fi 
    cd ..
done
