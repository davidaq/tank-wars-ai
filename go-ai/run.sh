cd `dirname $0`
export GOPATH=`pwd`

export HOST=ml.niven.cn:8777
export GAME=B1ngSs9yM
export SIDE=red
export TACTICS=nearest
# export TACTICS=proxy PROXY_PORT=8776

#go run src/ai-client.go &
SIDE=blue TACTICS=killall go run src/ai-client.go

# run forever
# yes|while read x; do go run src/ai-client.go; done
