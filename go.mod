module github.com/thisXYH/cache

go 1.18

require (
	github.com/cmstar/go-conv v0.3.1
	github.com/go-redis/redis/v8 v8.10.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
)

require (
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	go.opentelemetry.io/otel v0.20.0 // indirect
	go.opentelemetry.io/otel/metric v0.20.0 // indirect
	go.opentelemetry.io/otel/trace v0.20.0 // indirect
)

retract (
	v1.0.1
	v0.1.3
	v0.1.2
	v0.1.1
)
