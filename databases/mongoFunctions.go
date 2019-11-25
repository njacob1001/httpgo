package databases

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetDataByID func
func (mongo *MDB) GetDataByID(articuloID string) (*Product, error) {
	filter := bson.D{{
		"id",
		bson.D{{
			"$in",
			bson.A{articuloID},
		}},
	}}
	var result *Product
	err := mongo.Database("redesuao").Collection("basuras").FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CalculateSum func
func (mongo *MDB) CalculateSum(articles []string) (int64, error) {
	findOptions := options.Find()
	filter := bson.D{{
		"id",
		bson.D{{
			"$in",
			bson.A{articles},
		}},
	}}
	cur, err := mongo.Database("redesuao").Collection("basuras").Find(context.TODO(), filter, findOptions)
	if err != nil {
		return 0, err
	}
	var results []*Product
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Product
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	// return results, nil

	if len(results) > 0 {
		var total int64
		for _, art := range results {
			if art.Precio > 0 {
				total += art.Precio
			}
		}
		return total, nil
	}

	return 0, nil

	// return result, nil
}
