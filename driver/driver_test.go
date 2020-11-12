package driver

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
	"time"
)

//=== RUN   TestNewMongoDB
//map[_id:ObjectID("5faccc5d68e65ec0427ded3b") age:18 name:猪大肠 status:pering]
//map[_id:ObjectID("5faccc7368e65ec0427ded3c") age:18 name:猪大肠 status:pering]
//map[_id:ObjectID("5faccc7768e65ec0427ded3d") age:18 name:猪大肠 status:pering]
//map[_id:ObjectID("5faccc7768e65ec0427ded3e") age:18 name:猪大肠 status:pering]
//map[_id:ObjectID("5faccc7868e65ec0427ded3f") age:18 name:猪大肠 status:pering]
//map[_id:ObjectID("5faccc7968e65ec0427ded40") age:18 name:猪大肠 status:pering]
//map[_id:ObjectID("5faccc7a68e65ec0427ded41") age:18 name:猪大肠 status:pering]
//map[_id:ObjectID("5faccc7c68e65ec0427ded42") age:18 name:猪大肠 status:pering]
//map[_id:ObjectID("5faccca068e65ec0427ded43") age:18 name:猪大肠 status:pering]
//map[_id:ObjectID("5faccca168e65ec0427ded44") age:18 name:猪大肠 status:pering]
//--- PASS: TestNewMongoDB (0.62s)
//PASS
func TestNewMongoDB(t *testing.T) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), 4*time.Second)
		opt         = options.Find()

		//auth = options.Credential{
		//	AuthMechanism:           "GSSAPI",
		//	AuthMechanismProperties: nil,
		//	AuthSource:              "",
		//	Username:                "test",
		//	Password:                "123456",
		//	PasswordSet:             true,
		//}
	)

	mongoClient, err := NewMongoDB("test", options.Client().ApplyURI("mongodb://test:123456@localhost:27017/?serverSelectionTimeoutMS=5000&connectTimeoutMS=10000&authSource=test"), "users")
	if err != nil {
		log.Fatal(err)
	}
	cursor, err := mongoClient.Find(ctx, bson.D{}, opt.SetAllowDiskUse(true), opt.SetLimit(10), opt.SetSort(bson.M{
		"_id": -1,
	}))
	if err != nil {
		log.Fatal(err)
	}
	for cursor.Next(ctx) {
		var res bson.M
		if err := cursor.Decode(&res); err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(res)
	}
	defer func() {
		cursor.Close(ctx)
		cancel()
	}()
}
