# Mongoutil

General purpose utility for using MongoDB's go driver, allowing you to make cleaner, more readable code for your DB operations.  
The core of the package is the MongoClient struct, which represents a MongoDB client and has several methods for performing standard database CRUD
operations for documents.

**This project is under active development**

## Installation

To install Mongoutil, use `go get`:

```bash
go get github.com/ooyeku/mongoutil
```

## Usage
```Go
import "github.com/ooyeku/mongoutil/mgutil"
```
### Create a new MongoDB Client
```go
mc, err := mgutil.NewMongoClient(uri)
```
### Insert a Document
```go
err := mc.Create(database, collection, document)
```
### Read a Document
```go
err := mc.Read(database, collection, filter, &result)
```
### Find All Documents
```go
err := mc.FindAll(database, collection, &results)
```
### Update a Document
```go
err := mc.Update(database, collection, filter, update)
```
### Delete a Document
```go
err := mc.Delete(database, collection, filter)
```
### Close the MongoDB Connection
```go
err := mc.Close()
```



Additionally, the package has various other methods that provide more functionality like counting documents, running aggregations, finding distinct values, dropping collections and databases, listing databases and collections, creating indexes, inserting many documents, updating many documents, deleting many documents, and running a MongoDB transaction.

## Notes
All MongoDB operations in this package are using the context.Background() which means these operations will not be cancelled until they've finished. Depending on your use case, you might want to provide a context with set deadline or cancellation signals.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details