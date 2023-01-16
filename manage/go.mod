module bdlinstall/bdl

go 1.18

replace bdl.local/bdl => ../src/

require (
	bdl.local/bdl v0.0.0-00010101000000-000000000000
	github.com/mattn/go-adodb v0.0.1
	golang.org/x/text v0.6.0
)

require (
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/jmoiron/sqlx v1.2.0 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/lib/pq v1.3.0 // indirect
	golang.org/x/net v0.0.0-20190603091049-60506f45cf65 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
