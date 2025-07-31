package mongoUtil

import (
	"context"
	"reflect"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	indexCache     = make(map[string]bool)
	indexCacheLock sync.Mutex
)

// EnsureIndexes 确保集合的索引已创建
func EnsureIndexes(ctx context.Context, collection *mongo.Collection, model interface{}) error {
	typeName := reflect.TypeOf(model).String()
	cacheKey := collection.Name() + ":" + typeName

	indexCacheLock.Lock()
	defer indexCacheLock.Unlock()

	// 检查是否已经创建过索引
	if indexCache[cacheKey] {
		return nil
	}

	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("mongoIndex")
		if tag == "" {
			continue
		}

		bsonTag := field.Tag.Get("bson")
		if bsonTag == "" || bsonTag == "-" {
			continue
		}

		bsonName := strings.Split(bsonTag, ",")[0]

		indexOptions := options.Index()

		switch tag {
		case "unique":
			indexOptions.SetUnique(true)
			fallthrough
		case "index":
			_, err := collection.Indexes().CreateOne(
				ctx,
				mongo.IndexModel{
					Keys:    bson.D{{Key: bsonName, Value: 1}},
					Options: indexOptions,
				},
			)
			if err != nil {
				return err
			}
		}
	}

	// 标记为已创建
	indexCache[cacheKey] = true
	return nil
}
