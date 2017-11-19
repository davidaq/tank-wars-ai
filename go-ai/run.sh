cd `dirname $0`
export GOPATH=`pwd`

export HOST=ml.niven.cn:8777
<<<<<<< HEAD
export GAME=Sy12Akylf
export SIDE=blue
export TACTICS=waitsweep
# export TACTICS=proxy PROXY_PORT=8776

go run src/ai-client.go > /dev/null 2>&1 &
SIDE=red TACTICS=nearest go run src/ai-client.go
=======
export GAME=By6SwT0kz
export SIDE=blue
export TACTICS=nearest
# export TACTICS=proxy PROXY_PORT=8776

go run src/ai-client.go > /dev/null 2>&1 &
SIDE=red TACTICS=less go run src/ai-client.go
>>>>>>> 32fb8224692bce83760fe4006a9f56c102d245b7

# run forever
# yes|while read x; do go run src/ai-client.go; done
