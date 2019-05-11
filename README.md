vmrenter is a tool to checkout virtual machines from a shared pool. Once a reservation is made, a configuration json file is generated that can then be used to drive the automated installation of a MapR cluster 

To set this project up via IntelliJ, certain settings are required

IntelliJ Preferences -> Languages & Frameworks -> Go -> Go Module (vgo) -> Enable Go Modules (vgo) integration

# Environment Variables
URL_DB_CONN
DEFAULT_USERNAME
DEFAULT_PASSWORD

# To import data from CSV and load it into maprdb
run pkg/utils/DataImportUtils_test.go/TestDataSeed with the csv file