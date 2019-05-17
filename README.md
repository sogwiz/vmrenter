vmrenter is a tool to checkout virtual machines from a shared pool. Once a reservation is made, a configuration json file is generated that can then be used to drive the automated installation of a MapR cluster 

The backing database is a MapR DB with 2 JSON tables. This client project uses the MapRDB Go client bindings 
to access the Data Access Gateway, typically via grpc port 5678

# Docker
## build the image
docker build -t sogwiz/vmrenter .
## run the container
docker run -v ${workspace}/private-installer/testing/configuration:/src/configuration sogwiz/vmrenter -u "<DAG_HOST>:5678?auth=basic;user=mapr;password=<PASSWORD>;ssl=false" -f /src/configuration/config.json -n $NODES -o $OS -e $EMAIL

# Data 
## To create the MapRDB tables
run pkg/mapr/DBUtils_test.go/TestReset
- env vars: URL_DB_CONN

## To import data from CSV and load it into maprdb
run pkg/utils/DataImportUtils_test.go/TestDataSeed with the csv file
- env vars: URL_DB_CONN, DEFAULT_USERNAME (root ssh user for each node), DEFAULT_PASSWORD (root ssh for each node)

## Reservations
# Make reservation via test
run pkg/mapr/DBUtils_test.go/TestMakeReservation
- env vars: URL_DB_CONN

# Make reservation
run cmd/ConfigLoader.go
- parameters: -f
              "/Users/sargonbenjamin/dev/src/private-installer/testing/configuration/config.json"
              -u
              "<host>:5678?auth=basic;user=mapr;password=<PASSWORD>;ssl=false"
              -n
              1
              -o
              ubuntu

To set this project up via IntelliJ, certain settings are required

IntelliJ Preferences -> Languages & Frameworks -> Go -> Go Module (vgo) -> Enable Go Modules (vgo) integration

# Environment Variables
URL_DB_CONN
DEFAULT_USERNAME
DEFAULT_PASSWORD




