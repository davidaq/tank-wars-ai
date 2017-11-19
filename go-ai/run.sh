cd `dirname $0`
export GOPATH=`pwd`

export HOST=ml.niven.cn:8777
export GAME=S156C6RyM
export SIDE=blue
export TACTICS=sweep
# export TACTICS=proxy PROXY_PORT=8776

go run src/ai-client.go > /dev/null 2>&1 &
SIDE=red TACTICS=forest-patrol go run src/ai-client.go

# run forever
# yes|while read x; do go run src/ai-client.go; done
