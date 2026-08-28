package mongoUtil

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/drivertest"
)

type indexedModel struct {
	Email   string `bson:"email,omitempty" mongoIndex:"unique"`
	Age     int    `bson:"age" mongoIndex:"index"`
	Ignored string `bson:"-" mongoIndex:"index"`
	NoBSON  string `mongoIndex:"index"`
	Unknown string `bson:"unknown" mongoIndex:"text"`
}

type retryModel struct {
	Value string `bson:"value" mongoIndex:"index"`
}

func resetIndexCache(t *testing.T) {
	t.Helper()
	indexCacheLock.Lock()
	indexCache = make(map[string]bool)
	indexCacheLock.Unlock()
	t.Cleanup(func() {
		indexCacheLock.Lock()
		indexCache = make(map[string]bool)
		indexCacheLock.Unlock()
	})
}

func mockCollection(t *testing.T, name string, responses []bson.D, commands *[]bson.Raw, commandLock *sync.Mutex) (*mongo.Collection, *drivertest.MockDeployment) {
	t.Helper()
	deployment := drivertest.NewMockDeployment(responses...)
	monitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			commandLock.Lock()
			*commands = append(*commands, append(bson.Raw(nil), evt.Command...))
			commandLock.Unlock()
		},
	}
	clientOptions := options.Client().SetMonitor(monitor)
	clientOptions.Deployment = deployment
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		t.Fatalf("mongo.Connect(mock): %v", err)
	}
	t.Cleanup(func() {
		_ = client.Disconnect(context.Background())
	})
	return client.Database("testdb").Collection(name), deployment
}

func TestEnsureIndexesBuildsExpectedCommandsAndCaches(t *testing.T) {
	resetIndexCache(t)
	var commands []bson.Raw
	var commandLock sync.Mutex
	responses := []bson.D{
		{{Key: "ok", Value: 1}},
		{{Key: "ok", Value: 1}},
	}
	collection, _ := mockCollection(t, "indexed", responses, &commands, &commandLock)

	if err := EnsureIndexes(context.Background(), collection, &indexedModel{}); err != nil {
		t.Fatalf("EnsureIndexes(): %v", err)
	}
	commandLock.Lock()
	gotCommands := append([]bson.Raw(nil), commands...)
	commandLock.Unlock()
	if len(gotCommands) != 2 {
		t.Fatalf("createIndexes command count = %d, want 2", len(gotCommands))
	}
	assertIndexCommand(t, gotCommands[0], "indexed", "email", true)
	assertIndexCommand(t, gotCommands[1], "indexed", "age", false)

	if err := EnsureIndexes(context.Background(), collection, &indexedModel{}); err != nil {
		t.Fatalf("cached EnsureIndexes(): %v", err)
	}
	commandLock.Lock()
	defer commandLock.Unlock()
	if len(commands) != 2 {
		t.Errorf("cached EnsureIndexes() sent %d total commands, want 2", len(commands))
	}
}

func TestEnsureIndexesSkipsModelsWithoutUsableTags(t *testing.T) {
	resetIndexCache(t)
	var commands []bson.Raw
	var commandLock sync.Mutex
	collection, _ := mockCollection(t, "skipped", nil, &commands, &commandLock)

	type skippedModel struct {
		NoIndex string `bson:"noIndex"`
		NoBSON  string `mongoIndex:"index"`
		Ignored string `bson:"-" mongoIndex:"unique"`
		Unknown string `bson:"unknown" mongoIndex:"unsupported"`
	}
	for _, model := range []interface{}{skippedModel{}, &skippedModel{}, "not a struct"} {
		if err := EnsureIndexes(context.Background(), collection, model); err != nil {
			t.Errorf("EnsureIndexes(%T): %v", model, err)
		}
	}
	commandLock.Lock()
	defer commandLock.Unlock()
	if len(commands) != 0 {
		t.Errorf("models without usable tags sent commands: %v", commands)
	}
}

func TestEnsureIndexesNilModelIsNoOp(t *testing.T) {
	resetIndexCache(t)
	var commands []bson.Raw
	var commandLock sync.Mutex
	collection, _ := mockCollection(t, "nil-model", nil, &commands, &commandLock)

	if err := EnsureIndexes(context.Background(), collection, nil); err != nil {
		t.Fatalf("EnsureIndexes(nil): %v", err)
	}

	commandLock.Lock()
	defer commandLock.Unlock()
	if len(commands) != 0 {
		t.Errorf("nil model sent commands: %v", commands)
	}
}

func TestEnsureIndexesCachesValueAndPointerAsSameModel(t *testing.T) {
	resetIndexCache(t)
	var commands []bson.Raw
	var commandLock sync.Mutex
	responses := []bson.D{{{Key: "ok", Value: 1}}}
	collection, _ := mockCollection(t, "normalized-model", responses, &commands, &commandLock)

	if err := EnsureIndexes(context.Background(), collection, retryModel{}); err != nil {
		t.Fatalf("EnsureIndexes(value): %v", err)
	}
	if err := EnsureIndexes(context.Background(), collection, &retryModel{}); err != nil {
		t.Fatalf("EnsureIndexes(pointer): %v", err)
	}

	commandLock.Lock()
	defer commandLock.Unlock()
	if len(commands) != 1 {
		t.Errorf("value and pointer models sent %d commands, want 1", len(commands))
	}
}

func TestEnsureIndexesCacheIsScopedByDatabaseAndClient(t *testing.T) {
	t.Run("different databases on one client", func(t *testing.T) {
		resetIndexCache(t)
		var commands []bson.Raw
		var commandLock sync.Mutex
		responses := []bson.D{{{Key: "ok", Value: 1}}, {{Key: "ok", Value: 1}}}
		first, _ := mockCollection(t, "shared", responses, &commands, &commandLock)
		second := first.Database().Client().Database("otherdb").Collection("shared")

		if err := EnsureIndexes(context.Background(), first, retryModel{}); err != nil {
			t.Fatalf("EnsureIndexes(first database): %v", err)
		}
		if err := EnsureIndexes(context.Background(), second, retryModel{}); err != nil {
			t.Fatalf("EnsureIndexes(second database): %v", err)
		}
		commandLock.Lock()
		defer commandLock.Unlock()
		if len(commands) != 2 {
			t.Fatalf("different databases sent %d commands, want 2", len(commands))
		}
		if got := commands[0].Lookup("$db").StringValue(); got != "testdb" {
			t.Errorf("first command database = %q, want testdb", got)
		}
		if got := commands[1].Lookup("$db").StringValue(); got != "otherdb" {
			t.Errorf("second command database = %q, want otherdb", got)
		}
	})

	t.Run("different clients", func(t *testing.T) {
		resetIndexCache(t)
		var firstCommands, secondCommands []bson.Raw
		var firstLock, secondLock sync.Mutex
		first, _ := mockCollection(t, "shared", []bson.D{{{Key: "ok", Value: 1}}}, &firstCommands, &firstLock)
		second, _ := mockCollection(t, "shared", []bson.D{{{Key: "ok", Value: 1}}}, &secondCommands, &secondLock)

		if err := EnsureIndexes(context.Background(), first, retryModel{}); err != nil {
			t.Fatalf("EnsureIndexes(first client): %v", err)
		}
		if err := EnsureIndexes(context.Background(), second, retryModel{}); err != nil {
			t.Fatalf("EnsureIndexes(second client): %v", err)
		}
		firstLock.Lock()
		firstCount := len(firstCommands)
		firstLock.Unlock()
		secondLock.Lock()
		secondCount := len(secondCommands)
		secondLock.Unlock()
		if firstCount != 1 || secondCount != 1 {
			t.Errorf("different clients command counts = %d/%d, want 1/1", firstCount, secondCount)
		}
	})
}

func TestEnsureIndexesFailureIsRetryable(t *testing.T) {
	resetIndexCache(t)
	var commands []bson.Raw
	var commandLock sync.Mutex
	responses := []bson.D{
		{{Key: "ok", Value: 0}, {Key: "errmsg", Value: "create failed"}, {Key: "code", Value: 123}},
		{{Key: "ok", Value: 1}},
	}
	collection, _ := mockCollection(t, "retry", responses, &commands, &commandLock)

	if err := EnsureIndexes(context.Background(), collection, retryModel{}); err == nil {
		t.Fatal("EnsureIndexes() mock failure returned nil error")
	}
	cacheKey := makeIndexCacheKey(collection, reflect.TypeOf(retryModel{}))
	indexCacheLock.Lock()
	cachedAfterFailure := indexCache[cacheKey]
	indexCacheLock.Unlock()
	if cachedAfterFailure {
		t.Error("failed index creation was cached")
	}

	if err := EnsureIndexes(context.Background(), collection, retryModel{}); err != nil {
		t.Fatalf("EnsureIndexes() retry: %v", err)
	}
	indexCacheLock.Lock()
	cachedAfterSuccess := indexCache[cacheKey]
	indexCacheLock.Unlock()
	if !cachedAfterSuccess {
		t.Error("successful index creation was not cached")
	}
	commandLock.Lock()
	defer commandLock.Unlock()
	if len(commands) != 2 {
		t.Errorf("failure and retry sent %d commands, want 2", len(commands))
	}
}

func assertIndexCommand(t *testing.T, command bson.Raw, collectionName, keyName string, unique bool) {
	t.Helper()
	if got := command.Lookup("createIndexes").StringValue(); got != collectionName {
		t.Errorf("createIndexes collection = %q, want %q", got, collectionName)
	}
	indexValues, err := command.Lookup("indexes").Array().Values()
	if err != nil {
		t.Fatalf("decode indexes array: %v", err)
	}
	if len(indexValues) != 1 {
		t.Fatalf("indexes count = %d, want 1", len(indexValues))
	}
	index := indexValues[0].Document()
	keys := index.Lookup("key").Document()
	keyElements, err := keys.Elements()
	if err != nil {
		t.Fatalf("decode index key: %v", err)
	}
	if len(keyElements) != 1 || keyElements[0].Key() != keyName || keyElements[0].Value().Int32() != 1 {
		t.Errorf("index key = %v, want {%s: 1}", keys, keyName)
	}
	uniqueValue, ok := index.Lookup("unique").BooleanOK()
	if unique {
		if !ok || !uniqueValue {
			t.Errorf("unique option = (%v, %v), want true", uniqueValue, ok)
		}
	} else if ok {
		t.Errorf("non-unique index unexpectedly encoded unique option = %v", uniqueValue)
	}
}
