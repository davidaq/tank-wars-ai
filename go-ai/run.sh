cd `dirname $0`
export GOPATH=`pwd`

export HOST=ml.niven.cn:8777
export GAME=B1Wjp58kf
export SIDE=red

#export TACTICS=random
export TACTICS=simple

go run src/ai-client.go

# run forever
# yes|while read x; do go run src/ai-client.go; done
