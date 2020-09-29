module bdl.local/install

go 1.14

require bdl.local/bdl v0.0.0

replace bdl.local/bdl => ../src

require (
	github.com/recursionpharma/go-csv-map v0.0.0-20160524001940-792523c65ae9
	github.com/tealeg/xlsx/v3 v3.2.0
	golang.org/x/text v0.3.2
)
