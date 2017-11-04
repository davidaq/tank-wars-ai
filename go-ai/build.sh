cd `dirname $0`
export GOPATH=`pwd`

# 构建ai客户端
go build src/ai-client.go

# 构建ai服务器
go build src/ai-server.go
