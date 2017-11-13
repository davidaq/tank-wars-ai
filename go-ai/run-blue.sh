cd `dirname $0`
export GOPATH=`pwd`

export HOST=ml.niven.cn:8777
export GAME=rklRuKn0W
export SIDE=blue

export TACTICS=simple
# export TACTICS=proxy PROXY_PORT=8775

go run src/ai-client.go

# run forever
# yes|while read x; do go run src/ai-client.go; done
