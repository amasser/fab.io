package mongorecords

import (
	"fmt"
	"time"

	"github.com/kooinam/fab.io/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Base used to represent base classes for all models
type Base struct {
	collection   *models.Collection
	hooksHandler *models.HooksHandler
	item         models.Modellable
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
}

// InitializeBase used for setting up base attributes for a mongo record
func (base *Base) InitializeBase(context *models.Context) {
	base.collection = context.Collection()
	base.hooksHandler = context.HooksHandler()
	base.item = context.Item()
}

// GetCollectionName used to retrieve collection's name
func (base *Base) GetCollectionName() string {
	return base.collection.Name()
}

// GetHooksHandler used to retrieve hooks handler
func (base *Base) GetHooksHandler() *models.HooksHandler {
	return base.hooksHandler
}

// GetID used to retrieve record's ID
func (base *Base) GetID() string {
	return base.ID.Hex()
}

// Save used to save record in adapter
func (base *Base) Save() error {
	var err error

	err = base.hooksHandler.ExecuteValidationHooks()

	if err != nil {
		return err
	}

	adapter := base.collection.Adapter().(*Adapter)

	collection := adapter.getCollection(base.collection.Name())
	ctx := adapter.getTimeoutContext()

	if base.IsNewRecord() {
		base.ID = primitive.NewObjectID()
		base.CreatedAt = time.Now()
		base.UpdatedAt = time.Now()

		results, err2 := collection.InsertOne(ctx, base.item)

		err = err2

		if err != nil {
			base.ID = primitive.NilObjectID
			base.CreatedAt = time.Time{}
			base.UpdatedAt = time.Time{}
		} else {
			base.ID = results.InsertedID.(primitive.ObjectID)
		}
	} else {
		previousUpdatedAt := base.UpdatedAt
		base.UpdatedAt = time.Now()

		results, err2 := collection.ReplaceOne(ctx, bson.M{"_id": base.ID}, base.item)

		if err2 == nil && results.MatchedCount != 1 {
			err2 = fmt.Errorf("no matched record for update")
		}

		err = err2

		if err != nil {
			base.UpdatedAt = previousUpdatedAt
		}
	}

	return err
}

// IsNewRecord used to check if record is new unsaved record
func (base *Base) IsNewRecord() bool {
	return base.ID == primitive.NilObjectID
}

// Reload used to reload record from database
func (base *Base) Reload() error {
	var err error

	// TOOD

	return err
}

// Destroy used to delete record from database
func (base *Base) Destroy() error {
	var err error

	// TOOD

	return err
}

// Store used to add record to list
func (base *Base) Store() {
	base.collection.List().Add(base.item)
}

// StoreInList used to add record to selected list
func (base *Base) StoreInList(list *models.List) {
	list.Add(base.item)
}
