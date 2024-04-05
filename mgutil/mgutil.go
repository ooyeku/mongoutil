package mgutil

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient represents a MongoDB client.
type MongoClient struct {
	client *mongo.Client
}

// NewMongoClient creates a new instance of MongoClient by connecting to a MongoDB server
// and returns a pointer to the MongoClient and an error.
// The function takes a URI string as input, which specifies the MongoDB connection details.
// If the connection is successful, a new instance of MongoClient is created with the connected
// client, and the pointer to the MongoClient is returned along with a nil error.
// If the connection fails, a nil pointer is returned along with the error.
func NewMongoClient(uri string) (*MongoClient, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &MongoClient{client: client}, nil
}

// Create inserts a document into the specified database and collection.
// It returns an error if the insertion fails.
// Parameters:
// - database: the name of the database
// - collection: the name of the collection
// - document: the document to be inserted
// Returns:
// - error: an error indicating if the insertion was successful or not
func (mc *MongoClient) Create(database, collection string, document interface{}) error {
	_, err := mc.client.Database(database).Collection(collection).InsertOne(context.Background(), document)
	return err
}

// Read retrieves a single document from the specified database and collection
// that matches the given filter and decodes the result into the provided
// result object. If an error occurs during the operation, it is returned.
// The Read operation uses the context.Background() as the context.
//
// Parameters:
// - database: The name of the database.
// - collection: The name of the collection within the database.
// - filter: The filter to apply when searching for the document.
// - result: A pointer to the object where the result will be decoded into.
//
// Returns:
// - error: An error if any occurred during the operation, otherwise it is nil.
func (mc *MongoClient) Read(database, collection string, filter interface{}, result interface{}) error {
	err := mc.client.Database(database).Collection(collection).FindOne(context.Background(), filter).Decode(result)
	return err
}

// FindAll retrieves all documents from a specified collection in a database.
// It takes the database and collection names as parameters, as well as a pointer
// to a slice of bson.M, which will store the results.
// It returns an error if the operation fails.
func (mc *MongoClient) FindAll(database string, collection string, results *[]bson.M) error {
	coll := mc.client.Database(database).Collection(collection)
	cursor, err := coll.Find(context.TODO(), bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			return err
		}
		*results = append(*results, result)
	}

	if err := cursor.Err(); err != nil {
		return err
	}
	return nil
}

// Update updates a document in the specified database and collection.
//
// Parameters:
// - database: The name of the database.
// - collection: The name of the collection.
// - filter: The filter to identify the document to update.
// - update: The update document.
//
// Returns:
// - error: An error if the operation fails, or nil if successful.
func (mc *MongoClient) Update(database, collection string, filter, update interface{}) error {
	_, err := mc.client.Database(database).Collection(collection).UpdateOne(context.Background(), filter, update)
	return err
}

// Delete deletes a document from the specified database and collection based on the given filter.
// It returns an error if the deletion operation fails.
// The filter is used to identify the document to be deleted.
// Use with caution as this operation permanently removes the matching document(s).
func (mc *MongoClient) Delete(database, collection string, filter interface{}) error {
	_, err := mc.client.Database(database).Collection(collection).DeleteOne(context.Background(), filter)
	return err
}

// Close closes the connection to the MongoDB server.
// It calls the Disconnect method on the underlying mongo.Client object.
func (mc *MongoClient) Close() error {
	return mc.client.Disconnect(context.Background())
}

func (mc *MongoClient) CountDocuments(database, collection string, filter interface{}) (int64, error) {
	return mc.client.Database(database).Collection(collection).CountDocuments(context.Background(), filter)
}

func (mc *MongoClient) Aggregate(database, collection string, pipeline []bson.M) (*mongo.Cursor, error) {
	return mc.client.Database(database).Collection(collection).Aggregate(context.Background(), pipeline)
}

func (mc *MongoClient) Distinct(database, collection, field string, filter interface{}) ([]interface{}, error) {
	return mc.client.Database(database).Collection(collection).Distinct(context.Background(), field, filter)
}

func (mc *MongoClient) DropCollection(database, collection string) error {
	return mc.client.Database(database).Collection(collection).Drop(context.Background())
}

func (mc *MongoClient) DropDatabase(database string) error {
	return mc.client.Database(database).Drop(context.Background())
}

func (mc *MongoClient) ListDatabases() (mongo.ListDatabasesResult, error) {
	return mc.client.ListDatabases(context.Background(), bson.M{})
}

func (mc *MongoClient) ListCollections(database string) ([]string, error) {
	cursor, err := mc.client.Database(database).ListCollections(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var collections []string
	for cursor.Next(context.Background()) {
		var collection string
		if err := cursor.Decode(&collection); err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}

	return collections, nil
}

func (mc *MongoClient) CreateIndex(database, collection string, keys interface{}, opts *options.IndexOptions) (string, error) {
	return mc.client.Database(database).Collection(collection).Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: keys, Options: opts})
}

// InsertMany inserts multiple documents into the specified database and collection
func (mc *MongoClient) InsertMany(database, collection string, documents []interface{}) error {
	_, err := mc.client.Database(database).Collection(collection).InsertMany(context.Background(), documents)
	return err
}

// UpdateMany updates multiple documents in the specified database and collection
func (mc *MongoClient) UpdateMany(database, collection string, filter interface{}, update interface{}) error {
	_, err := mc.client.Database(database).Collection(collection).UpdateMany(context.Background(), filter, update)
	return err
}

// DeleteMany deletes multiple documents from the specified database and collection
func (mc *MongoClient) DeleteMany(database, collection string, filter interface{}) error {
	_, err := mc.client.Database(database).Collection(collection).DeleteMany(context.Background(), filter)
	return err
}

// RunTransaction executes the provided callback function within a MongoDB transaction
func (mc *MongoClient) RunTransaction(database string, callback func(sessCtx mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	session, err := mc.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(context.Background())

	var result interface{}
	err = session.StartTransaction()
	if err != nil {
		return nil, err
	}
	defer func(session mongo.Session, ctx context.Context) {
		err := session.CommitTransaction(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}(session, context.Background())

	result, err = callback(mongo.NewSessionContext(context.Background(), session))
	if err != nil {
		err := session.AbortTransaction(context.Background())
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	err = session.CommitTransaction(context.Background())
	if err != nil {
		return nil, err
	}
	return result, nil
}
