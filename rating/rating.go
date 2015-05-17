package rating

import "gopkg.in/mgo.v2/bson"

type Rating struct {
    ID bson.ObjectId `bson:"_id,omitempty"`
    Tenant string
    Category string
    ItemId string
    Rating uint
    RatingOn uint
}
