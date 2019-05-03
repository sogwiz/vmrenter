
To set this project up via IntelliJ, certain settings are required

IntelliJ Preferences -> Languages & Frameworks -> Go -> Go Module (vgo) -> Enable Go Modules (vgo) integration

# Environment Variables
URL_DB_CONN
DEFAULT_USERNAME
DEFAULT_PASSWORD

# To import data from CSV and load it into maprdb
run pkg/utils/DataImportUtils_test.go/TestDataSeed with the csv file