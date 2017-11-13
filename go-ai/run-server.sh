cd `dirname $0`
export GOPATH=`pwd`

export PORT=8787

export TACTICS=random

go run src/ai-server.go

# run forever
# yes|while read x; do go run src/ai-client.go; done
