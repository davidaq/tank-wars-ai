cd `dirname $0`
export GOPATH=`pwd`

export HOST=ml.niven.cn:8777
export GAME=ryImt-sAW
export SIDE=red

go run src/ai-client.go
