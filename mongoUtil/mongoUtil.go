package mongoUtil

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"time"
)

type MongoUtil struct {
	client   *mongo.Client
	database *mongo.Database
	config   *Config
	Ctx      context.Context
}

type Config struct {
	URI         string
	Database    string
	MaxPoolSize uint64
	MinPoolSize uint64
}

func New(config *Config) (*MongoUtil, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(config.URI).
		SetMaxPoolSize(config.MaxPoolSize).
		SetMinPoolSize(config.MinPoolSize)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %v", err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping failed: %v", err)
	}
	return &MongoUtil{
		client:   client,
		database: client.Database(config.Database),
		config:   config,
	}, nil
}

// Collection 获取集合
func (mt *MongoUtil) Collection(collectionName string) *mongo.Collection {
	return mt.database.Collection(collectionName)
}

// Save 插入文档
func (mt *MongoUtil) Save(document interface{}) (interface{}, error) {
	collectionName := getCollectionName(document)
	insertResult, err := mt.database.Collection(collectionName).InsertOne(mt.Ctx, document)
	if err != nil {
		return nil, err
	}
	return insertResult.InsertedID, nil
}

// GetById 根据ID查询文档
func (mt *MongoUtil) GetById(id interface{}, document interface{}) error {
	// 参数有效性检查
	if document == nil {
		return errors.New("document must be a non-nil pointer")
	}
	val := reflect.ValueOf(document)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("document must be a non-nil pointer")
	}
	// 获取集合名称
	collectionName := getCollectionName(document)
	// 转换ID为ObjectID
	objectID, err := convertToObjectID(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}
	// 构造查询条件
	filter := bson.M{"_id": objectID}
	// 执行查询并解码结果
	if err = mt.database.Collection(collectionName).FindOne(mt.Ctx, filter).Decode(document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return mongo.ErrNoDocuments
		}
		return fmt.Errorf("document decoding failed: %w", err)
	}
	return nil
}

// UpdateById 根据文档的ID字段更新整个文档
func (mt *MongoUtil) UpdateById(document interface{}) error {
	// 获取集合名称
	collectionName := getCollectionName(document)
	// 反射解析ID字段
	idValue, err := extractIDValue(document)
	if err != nil {
		return fmt.Errorf("failed to get document ID: %w", err)
	}
	// 构建过滤器
	filter := bson.M{"_id": idValue}
	// 执行更新操作（使用替换整个文档的方式）
	result, err := mt.database.Collection(collectionName).ReplaceOne(mt.Ctx, filter, document)
	if err != nil {
		return fmt.Errorf("update operation failed: %w", err)
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// UpdateFieldById 根据ID部分更新文档字段
func (mt *MongoUtil) UpdateFieldById(id interface{}, document interface{}, field string, value interface{}) error {
	collectionName := getCollectionName(document)
	// 参数校验
	if field == "" {
		return errors.New("field name cannot be empty")
	}
	// 转换ID为ObjectID
	objectID, err := convertToObjectID(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}
	// 构造更新文档
	update := bson.M{"$set": bson.M{field: value}}
	// 执行更新
	result, err := mt.database.Collection(collectionName).UpdateByID(mt.Ctx, objectID, update)
	if err != nil {
		return fmt.Errorf("update operation failed: %w", err)
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (mt *MongoUtil) DeleteById(id interface{}, document interface{}) error {
	collectionName := getCollectionName(document)
	// 参数校验
	if collectionName == "" {
		return errors.New("collection name cannot be empty")
	}
	// 转换ID为ObjectID
	objectID, err := convertToObjectID(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}
	// 执行删除操作
	_, err = mt.database.Collection(collectionName).DeleteOne(mt.Ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("delete operation failed: %w", err)
	}
	return nil
}

// 辅助方法 -----------------------------------------------------
// getCollectionName 通过反射自动获取集合名称
func getCollectionName(document interface{}) string {
	t := reflect.TypeOf(document)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

// extractIDValue 从文档中提取ID值（支持多种ID类型）
func extractIDValue(document interface{}) (interface{}, error) {
	s := structs.New(document)
	// 查找带有 _id 标签的字段
	for _, field := range s.Fields() {
		if tag := field.Tag("bson"); tag == "_id" || tag == "_id,omitempty" {
			return convertToPrimitiveID(field.Value())
		}
	}
	// 查找名为 ID/Id/id 的字段
	for _, fieldName := range []string{"ID", "Id", "id"} {
		if field, ok := s.FieldOk(fieldName); ok {
			return convertToPrimitiveID(field.Value())
		}
	}
	return nil, fmt.Errorf("no valid ID field found in document")
}

// convertToObjectID 统一处理不同类型的ID
func convertToObjectID(id interface{}) (primitive.ObjectID, error) {
	switch v := id.(type) {
	case primitive.ObjectID:
		return v, nil
	case string:
		return primitive.ObjectIDFromHex(v)
	case []byte:
		if len(v) != 12 {
			return primitive.NilObjectID, errors.New("invalid byte length for ObjectID")
		}
		return primitive.ObjectID(v), nil
	default:
		return primitive.NilObjectID, fmt.Errorf("unsupported ID type: %T", id)
	}
}

// convertToPrimitiveID 转换为MongoDB支持的ID类型
func convertToPrimitiveID(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case primitive.ObjectID:
		return v, nil
	case string:
		if oid, err := primitive.ObjectIDFromHex(v); err == nil {
			return oid, nil
		}
		return v, nil // 支持字符串ID
	case []byte:
		if len(v) == 12 {
			return primitive.ObjectID(v), nil
		}
	}
	// 尝试反射转换
	if reflect.TypeOf(value).ConvertibleTo(typeObjectID) {
		return reflect.ValueOf(value).Convert(typeObjectID).Interface(), nil
	}
	return nil, fmt.Errorf("unsupported ID type: %T", value)
}

var typeObjectID = reflect.TypeOf(primitive.ObjectID{})
