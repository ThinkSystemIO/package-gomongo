package gomongo

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// handleError is a generic error handling function
// for testing with errors.
func handleError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

// TestIntegration runs integration testing for a local MongoDB database
// running on localhost with exposed port 27017.
// TODO - Test each individual database call
func TestIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	handleError(t, err)

	collection := client.Database("test").Collection("test")

	// ADD
	fmt.Println()
	fmt.Println("ADD")
	data := map[string]interface{}{
		"item": "a",
	}
	dataAdd, err := Add(collection, data)
	fmt.Println(dataAdd)
	handleError(t, err)
	all, err := GetAll(collection)
	handleError(t, err)
	fmt.Println(all)
	fmt.Println()

	// UPDATE
	fmt.Println("UPDATE")
	filter := map[string]interface{}{
		"item": "a",
	}
	update := map[string]interface{}{
		"$set": map[string]interface{}{"item": "b"},
	}
	dataUpdate, err := Update(collection, filter, update)
	fmt.Println(dataUpdate)
	handleError(t, err)
	all, err = GetAll(collection)
	handleError(t, err)
	fmt.Println(all)
	fmt.Println()

	// DELETE
	fmt.Println("DELETE")
	filter = map[string]interface{}{
		"item": "b",
	}
	dataDelete, err := Remove(collection, filter)
	fmt.Println(dataDelete)
	handleError(t, err)
	all, err = GetAll(collection)
	handleError(t, err)
	fmt.Println(all)
	fmt.Println()

	// ADD With ID
	fmt.Println("ADD ID")
	data = map[string]interface{}{
		"item": "a",
	}
	id, err := Add(collection, data)
	handleError(t, err)
	all, err = GetAll(collection)
	handleError(t, err)
	fmt.Println(all)
	fmt.Println()

	// UPDATE With ID
	fmt.Println("UPDATE ID")
	update = map[string]interface{}{
		"$set": map[string]interface{}{"item": "b"},
	}
	dataUpdateByID, err := UpdateByID(collection, id, update)
	fmt.Println(dataUpdateByID)
	fmt.Println(reflect.TypeOf(dataUpdateByID))
	handleError(t, err)
	all, err = GetAll(collection)
	handleError(t, err)
	fmt.Println(all)
	fmt.Println()

	// DELETE With ID
	fmt.Println("DELETE ID")
	dataRemoveByID, err := RemoveByID(collection, id)
	fmt.Println(dataRemoveByID)
	fmt.Println(reflect.TypeOf(dataRemoveByID))
	handleError(t, err)
	all, err = GetAll(collection)
	handleError(t, err)
	fmt.Println(all)
	fmt.Println()
}
