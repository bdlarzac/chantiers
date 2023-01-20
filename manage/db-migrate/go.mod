module bdl.dbmigrate/bdl

go 1.18

// replace bdl.local/bdl => ../../src/
// replace bdl.dbinstall/bdl => ../dbinstall

require (
//	bdl.local/bdl v0.0.0-00010101000000-000000000000
	github.com/jmoiron/sqlx v1.3.5
)

require (
//	bdl.dbinstall/bdl v0.0.0-00010101000000-000000000000 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/lib/pq v1.3.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
