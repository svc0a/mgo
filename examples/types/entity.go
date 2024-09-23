package types

// Entity @generated
type Entity struct {
	ID        int64 `json:"id" bson:"_id"`
	CreatedAt int64 `json:"createdAt" bson:"createdAt"`
	UpdatedAt int64 `json:"updatedAt" bson:"updatedAt"`
	Version   int64 `json:"version" bson:"version"`
}
