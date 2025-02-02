package data

import "go.mongodb.org/mongo-driver/bson"

type KeywordFilter struct {
	Keyword string
}

func (keywordFilter *KeywordFilter) ToLogBsonM() bson.M {
	filter := bson.M{}

	if keywordFilter != nil && keywordFilter.Keyword != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"record_id": bson.M{"$regex": keywordFilter.Keyword, "$options": "i"}},
				{"user_id": bson.M{"$regex": keywordFilter.Keyword, "$options": "i"}},
				{"_id": bson.M{"$regex": keywordFilter.Keyword, "$options": "i"}},
			},
		}
	}

	return filter
}
