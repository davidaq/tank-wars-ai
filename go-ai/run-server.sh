cd `dirname $0`
export GOPATH=`pwd`

export PORT=8787
export TACTICS=simple

go run src/ai-server.go
