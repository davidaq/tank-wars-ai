cd `dirname $0`
export GOPATH=`pwd`

export HOST=vikki.wang:8777
export GAME=B1cuWR8JG
export SIDE=red
export TACTICS=simple
# export TACTICS=proxy PROXY_PORT=8776

go run src/ai-client.go

# run forever
# yes|while read x; do go run src/ai-client.go; done
