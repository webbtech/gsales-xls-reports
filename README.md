# Gales Sales Spreadsheet Report File Download Service

Golang service to manage various xls report file downloads

## Version 2 Updates

- decided to remove the `mongodb.Close()` method. We're not at risk of having too many open connections
