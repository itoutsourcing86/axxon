
protoc -I %GOPATH%\src\axxon\api\proto\v1\ --go_out=plugins=grpc:%GOPATH%\src\axxon\pkg\api\v1\ %GOPATH%\src\axxon\api\proto\v1\fetch-service.proto
protoc -I %GOPATH%\src\axxon\api\proto\v1\ --grpc-gateway_out=logtostderr=true:%GOPATH%\src\axxon\pkg\api\v1\ %GOPATH%\src\axxon\api\proto\v1\fetch-service.proto