package gomongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateMongoClient creates a new MongoDB client and establishes
// a connection
func CreateMongoClient(service string) (*mongo.Client, error) {
	uri := "mongodb://" + service + ":27017"

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// defaultContext returns a default context and cancel function.
func defaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

// jsonToBson converts a JSON document to a BSON document.
func jsonToBson(data map[string]interface{}) ([]byte, error) {
	doc, err := bson.Marshal(data)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// Add inserts a JSON document into a collection and returns the document ID.
func Add(collection *mongo.Collection, data map[string]interface{}) (string, error) {
	ctx, cancel := defaultContext()
	defer cancel()

	doc, err := jsonToBson(data)
	if err != nil {
		return "", err
	}

	result, err := collection.InsertOne(ctx, doc)
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, err
}

// Update edits a document in a collection based on the filter JSON. It currently
// sets every field given in the update JSON, meaning it will overwrite arrays and
// nested structures with whatever is provided.
func Update(collection *mongo.Collection, filter map[string]interface{}, update map[string]interface{}) (interface{}, error) {
	ctx, cancel := defaultContext()
	defer cancel()

	filterDoc, err := jsonToBson(filter)
	if err != nil {
		return nil, err
	}

	updateDoc, err := jsonToBson(update)
	if err != nil {
		return nil, err
	}

	after := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	result := collection.FindOneAndUpdate(ctx, filterDoc, updateDoc, opts)
	if result.Err() != nil {
		return nil, result.Err()
	}

	doc := bson.D{}
	err = result.Decode(&doc)
	if err != nil {
		return nil, err
	}

	return doc.Map(), nil
}

// Update edits a document in a collection based on the document ID. It currently
// sets every field given in the update JSON, meaning it will overwrite arrays and
// nested structures with whatever is provided.
func UpdateByID(collection *mongo.Collection, ID string, update map[string]interface{}) (interface{}, error) {
	ctx, cancel := defaultContext()
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	filterDoc := bson.M{
		"_id": objectID,
	}

	updateDoc, err := jsonToBson(update)
	if err != nil {
		return nil, err
	}

	after := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	result := collection.FindOneAndUpdate(ctx, filterDoc, updateDoc, opts)
	if result.Err() != nil {
		return nil, result.Err()
	}

	doc := bson.D{}
	err = result.Decode(&doc)
	if err != nil {
		return nil, err
	}

	return doc.Map(), nil
}

// Remove deletes a document from a collection based on the filter JSON.
func Remove(collection *mongo.Collection, filter map[string]interface{}) (interface{}, error) {
	ctx, cancel := defaultContext()
	defer cancel()

	filterDoc, err := jsonToBson(filter)
	if err != nil {
		return nil, err
	}

	result := collection.FindOneAndDelete(ctx, filterDoc)
	if result.Err() != nil {
		return nil, result.Err()
	}

	doc := bson.D{}
	err = result.Decode(&doc)
	if err != nil {
		return nil, err
	}

	return doc.Map(), nil
}

// Remove deletes a document from a collection based on the document ID.
func RemoveByID(collection *mongo.Collection, ID string) (interface{}, error) {
	ctx, cancel := defaultContext()
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	filterDoc := bson.M{
		"_id": objectID,
	}

	result := collection.FindOneAndDelete(ctx, filterDoc)
	if result.Err() != nil {
		return nil, result.Err()
	}

	doc := bson.D{}
	err = result.Decode(&doc)
	if err != nil {
		return nil, err
	}

	return doc.Map(), nil
}

// GetAll returns all documents in the collection given.
func GetAll(collection *mongo.Collection) ([]map[string]interface{}, error) {
	ctx, cancel := defaultContext()
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	documents := []map[string]interface{}{}
	err = cur.All(ctx, &documents)
	if err != nil {
		return nil, err
	}

	return documents, nil
}
