cd `dirname $0`
export GOPATH=`pwd`

export HOST=ml.niven.cn:8777
export GAME=S1VaH8eJf
export SIDE=red

#export TACTICS=random
export TACTICS=proxy PROXY_PORT=8776

go run src/ai-client.go

# run forever
# yes|while read x; do go run src/ai-client.go; done
