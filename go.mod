module github.com/Liquid-Labs/lc-entities-model

require (
	github.com/Liquid-Labs/env v1.0.0-beta.0
	github.com/Liquid-Labs/lc-rdb-service v1.0.0-alpha.1
	github.com/Liquid-Labs/terror v1.0.0-alpha.0
	github.com/go-pg/pg v8.0.5+incompatible
	github.com/stretchr/testify v1.4.0
)

replace github.com/Liquid-Labs/lc-rdb-service => ../lc-rdb-service

replace github.com/Liquid-Labs/terror => ../terror
