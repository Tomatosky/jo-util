package mongoUtil

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	indexCache     = make(map[string]bool)
	indexCacheLock sync.Mutex
)

// EnsureIndexes 确保集合的索引已创建
func EnsureIndexes(ctx context.Context, collection *mongo.Collection, model interface{}) error {
	t := reflect.TypeOf(model)
	if t == nil {
		return nil
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	cacheKey := makeIndexCacheKey(collection, t)

	indexCacheLock.Lock()
	defer indexCacheLock.Unlock()

	// 检查是否已经创建过索引
	if indexCache[cacheKey] {
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
		if bsonName == "" {
			bsonName = strings.ToLower(field.Name)
		}

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

func makeIndexCacheKey(collection *mongo.Collection, modelType reflect.Type) string {
	database := collection.Database()
	return fmt.Sprintf("%p:%s:%s:%s", database.Client(), database.Name(), collection.Name(), modelType.String())
}
