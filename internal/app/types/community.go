package types

type (
	CommunityInfo struct {
		ID            string `json:"id,omitempty" bson:"id,omitempty"`
		Name          string `json:"name,omitempty" bson:"name,omitempty"`
		CreatedByID   string `json:"created_by_id,omitempty" bson:"created_by_id,omitempty"`
		CreatedByName string `json:"created_by_name,omitempty" bson:"created_by_name,omitempty"`
	}
)
