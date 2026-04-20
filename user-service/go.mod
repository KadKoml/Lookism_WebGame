module github.com/kozie/lookism-rpg/user-service

go 1.25.0

require (
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/kozie/lookism-rpg/api v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.48.0
	google.golang.org/grpc v1.81.1
)

require (
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260226221140-a57be14db171 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/kozie/lookism-rpg/api => ../api
