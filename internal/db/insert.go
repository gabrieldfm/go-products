package db

import (
	"context"
	"fmt"
)

func Insert(collection string, data interface{}) error {
	client, ctx := getConnection()
	defer client.Disconnect(ctx)

	c := client.Database("crawler").Collection(collection)

	_, err := c.InsertOne(context.Background(), data)

	if err != nil {
		fmt.Println("error: ", err)
	}

	return err
}
