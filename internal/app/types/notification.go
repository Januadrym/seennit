package types

type (
	PostCreateNotification struct {
		Post PostInfo `json:"post"`
	}

	PushNotification struct {
		ID      string   `json:"id,omitempty" bson:"id,omitempty"`
		Message string   `json:"message,omitempty" bson:"message,omitempty"`
		PostID  string   `json:"post_id,omitempty" bson:"post_id,omitempty"`
		UserID  []string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	}
)
