package mgutil

import (
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
	_ "time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMongoClient(t *testing.T) {
	/*
		Note: Test assumes Mongo DB connection string established as MONGO_URI
	*/
	mongoURI := os.Getenv("MONGO_URI")
	// Checking if MongoClient instance is created successfully
	mc, err := NewMongoClient(mongoURI)
	require.NoError(t, err)
	defer mc.Close()

	database, collection := "testDB", "testColl"

	// TestCase for Create.
	t.Run("Create", func(t *testing.T) {
		doc := bson.M{"foo": "bar"}
		assert.NoError(t, mc.Create(database, collection, doc))
	})

	// TestCase for Read.
	t.Run("Read", func(t *testing.T) {
		filter := bson.M{"foo": "bar"}
		var result bson.M
		err := mc.Read(database, collection, filter, &result)
		assert.NoError(t, err)
		assert.Equal(t, filter["foo"], result["foo"])
	})

	// TestCase for FindAll.
	t.Run("FindAll", func(t *testing.T) {
		var results []bson.M
		assert.NoError(t, mc.FindAll(database, collection, &results))
	})

	// TestCase for Update.
	t.Run("Update", func(t *testing.T) {
		filter := bson.M{"foo": "bar"}
		update := bson.M{"$set": bson.M{"foo": "baz"}}
		assert.NoError(t, mc.Update(database, collection, filter, update))
	})

	// TestCase for Delete.
	t.Run("Delete", func(t *testing.T) {
		filter := bson.M{"foo": "baz"}
		assert.NoError(t, mc.Delete(database, collection, filter))
	})

	// TestCase for CountDocuments.
	t.Run("CountDocuments", func(t *testing.T) {
		filter := bson.M{}
		count, err := mc.CountDocuments(database, collection, filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	// TestCase for InsertMany.
	t.Run("InsertMany", func(t *testing.T) {
		docs := []interface{}{bson.M{"foo1": "bar1"}, bson.M{"foo2": "bar2"}}
		assert.NoError(t, mc.InsertMany(database, collection, docs))
	})

	// TestCase for UpdateMany.
	t.Run("UpdateMany", func(t *testing.T) {
		filter := bson.M{"foo1": "bar1"}
		update := bson.M{"$set": bson.M{"bar1": "foo1"}}
		assert.NoError(t, mc.UpdateMany(database, collection, filter, update))
	})

	// TestCase for DeleteMany.
	t.Run("DeleteMany", func(t *testing.T) {
		filter := bson.M{"bar1": "foo1"}
		assert.NoError(t, mc.DeleteMany(database, collection, filter))
	})

}
