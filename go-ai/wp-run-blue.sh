cd `dirname $0`
export GOPATH=`pwd`

export HOST=vikki.wang:8777
# export GAME=B1cuWR8JG  # 四辆坦克
export GAME=HJdFmjtyG    # 五辆坦克
export SIDE=blue
# export TACTICS=simple
export TACTICS=nearest
# export TACTICS=proxy PROXY_PORT=8775

go run src/ai-client.go

# run forever
# yes|while read x; do go run src/ai-client.go; done
