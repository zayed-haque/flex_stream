module github.com/zayed-haque/flex_stream/api_service

go 1.22.6

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/lib/pq v1.10.9
	github.com/zayed-haque/flex_stream v0.0.0-20240822154337-342c145db5a8
	google.golang.org/grpc v1.65.0
)

require (
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240823204242-4ba0660f739c // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace github.com/yourusername/FlexStream => ../
