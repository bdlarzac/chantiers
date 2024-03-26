module bdl.dbinstall/bdl

go 1.22

// replace bdl.local/bdl => ../../src/
// replace bdl.dbmigrate/bdl => ../db-migrate

//	bdl.local/bdl v0.0.0-00010101000000-000000000000
require golang.org/x/text v0.6.0

require (
	github.com/jmoiron/sqlx v1.2.0 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/lib/pq v1.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
