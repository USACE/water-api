#!/bin/bash
# Usage using comment after ')' for each argument
usage(){ echo "$0 usage:" && grep " .)\ #" $0; exit 0;}
dc="docker compose"
# getopts
while getopts "bdmshu" arg; do
    case $arg in
        s) # Stop Docker
            s="stop"
            ;;
        d) # Down Docker
            d="down"
            ;;
        m) # Add minio '-f docker-compose.yml -f docker-compose.minio.yml'
            m="-f docker-compose.yml -f docker-compose.minio.yml"
            ;;
        u) # Docker compose up
            u="up"
            ;;
        b) # Start Docker with --build flag
            b="--build"
            ;;
        h | *) # Display help
            usage
            exit 0
            ;;
    esac
done

if [[ ! -z ${s} && ! -z ${m} ]]
then
    cmd="$dc $m $s"
elif [ ! -z ${s} ]
then
    cmd="$dc $s"
elif [[ ! -z ${d} && ! -z ${m} ]]
then
    cmd="$dc $m $d"
elif [ ! -z ${d} ]
then
    cmd="$dc $d"
elif [[  ! -z ${m} || ! -z ${b} ]]
then
    cmd="$dc $m up $b"
fi

[ $# -eq 0 ] && cmd="$dc up"

echo "$cmd"
eval "$cmd"