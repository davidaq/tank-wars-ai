cd `dirname $0`
export GOPATH=`pwd`

export HOST=localhost:8777
export GAME=HJgWiiZlM
export SIDE=red
export TACTICS=fox
# export TACTICS=proxy PROXY_PORT=8776

go run src/ai-client.go > /dev/null 2>&1 &
SIDE=blue TACTICS=cattycat go run src/ai-client.go

# run forever
# yes|while read x; do go run src/ai-client.go; done
