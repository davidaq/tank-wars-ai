cd `dirname $0`
export GOPATH=`pwd`

export HOST=vikki.wang:8777
export GAME=SkgxX_p1f
export SIDE=red
export TACTICS=nearest
# export TACTICS=proxy PROXY_PORT=8775

go run src/ai-client.go

# run forever
# yes|while read x; do go run src/ai-client.go; done
