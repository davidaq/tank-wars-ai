cd `dirname $0`
export GOPATH=`pwd`

export HOST=vikki.wang:8777
export GAME=SkgxX_p1f  # 四辆坦克
export SIDE=blue
export TACTICS=cattycat
# export TACTICS=proxy PROXY_PORT=8776

go run src/ai-client.go

# run forever
# yes|while read x; do go run src/ai-client.go; done
