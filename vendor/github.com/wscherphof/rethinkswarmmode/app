#!/bin/bash

usage ()
{
    echo
    echo "Usage: $(basename $0) [OPTIONS] NAME IMAGE SWARM"
    echo
    echo "Create or update a swarm application service"
    echo
    echo "NAME   service name"
    echo "IMAGE  repo/name:tag identifying the image"
    echo "SWARM  swarm to create the service on"
    echo
    echo "Options:"
    echo "  -e key=value ...  environment variables"
    echo "  -n network        swarm overlay network the service connects to (default: dbnet)"
    echo "  -p port ...       ports to publish (without creating an ssh tunnel)"
    echo "  -t port ...       ports to publish (and create an ssh tunnel to)"
    echo "  -r replicas       number of replicas to run (default: 1)"
    echo
    echo "A volume appdata is mounted on /appdata"
    echo
}

while getopts "e:n:p:t:r:h" opt; do
    case $opt in
        e  ) ENVS+=("$OPTARG");;
        n  ) NETWORK="$OPTARG";;
        p  ) PORTS+=("$OPTARG");;
        t  ) TUNNELS+=("$OPTARG");;
        r  ) REPLICAS="$OPTARG";;
        h  ) usage; exit;;
        \? ) echo "Unknown option: -$OPTARG" >&2; exit 1;;
        :  ) echo "Missing option argument for -$OPTARG" >&2; exit 1;;
        *  ) echo "Unimplemented option: -$OPTARG" >&2; exit 1;;
    esac
done
shift $((OPTIND -1))

NAME="$1"
TAG="$2"
SWARM="$3"
if [ ! "$TAG" -o ! "$NAME" -o ! "$SWARM" ]; then
    usage
    exit 1
fi
REPLICAS=${REPLICAS-1}
NETWORK=${NETWORK-dbnet}

DOCKER="docker-machine ssh ${SWARM}-manager-1 sudo docker"

echo "* creating appdata..."
${DOCKER} volume create --name appdata

${DOCKER} service ps $NAME &>/dev/null
if [ "$?" = "0" ]; then
    for env in "${ENVS[@]}"; do
        echo "* setting env ${env}..."
        ${DOCKER} service update --env-add ${env} ${NAME}
    done
    echo "* updating image to ${TAG}..."
	${DOCKER} service update --image ${TAG} ${NAME}
    echo "* scaling..."
	${DOCKER} service scale ${NAME}=${REPLICAS}
else
	PUBLISH=""
    for port in "${PORTS[@]}"; do
        PUBLISH="${PUBLISH} --publish ${port}:${port}"
    done
    for port in "${TUNNELS[@]}"; do
        PUBLISH="${PUBLISH} --publish ${port}:${port}"
    done
    ENVIRONMENT=""
    for env in "${ENVS[@]}"; do
        ENVIRONMENT="${ENVIRONMENT} -e ${env}"
    done
    echo "* starting service..."
	${DOCKER} service create --mount src=appdata,dst=/appdata --name ${NAME} --replicas ${REPLICAS} ${ENVIRONMENT} --network ${NETWORK} ${PUBLISH} ${TAG}
fi

if [ "${TUNNELS[@]}" ]; then
    echo "* connecting..."
    sleep 15
    for port in "${TUNNELS[@]}"; do
        $(dirname "$0")/util/tunnel $SWARM $port
    done
fi
