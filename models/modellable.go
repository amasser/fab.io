package models

// Modellable is the interface for all models implementing model's functionalities
type Modellable interface {
	InitializeBase(*Context)
	GetID() string
	Save() error
	GetHooksHandler() *HooksHandler
	Store()
	StoreInList(*List)
}
