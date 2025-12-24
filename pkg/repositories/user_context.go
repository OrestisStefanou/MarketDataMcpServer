package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"market_data_mcp_server/pkg/domain"
	market_data_mcp_serverErr "market_data_mcp_server/pkg/errors"
	"time"

	"github.com/dgraph-io/badger/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserContextRepository struct {
	db *badger.DB
}

func NewUserContextRepository(db *badger.DB) (*UserContextRepository, error) {
	return &UserContextRepository{db: db}, nil
}

func (r *UserContextRepository) GetUserContext(userID string) (domain.UserContext, error) {
	var userContext domain.UserContext
	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(userID))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return market_data_mcp_serverErr.UserContextNotFoundError{UserID: userID}
			}
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &userContext)
		})
	})

	return userContext, err
}

func (r *UserContextRepository) InsertUserContext(userContext domain.UserContext) error {
	err := r.db.Update(func(txn *badger.Txn) error {
		userContextBytes, err := json.Marshal(userContext)
		if err != nil {
			return err
		}

		return txn.Set([]byte(userContext.UserID), userContextBytes)
	})

	return err
}

func (r *UserContextRepository) UpdateUserContext(userContext domain.UserContext) error {
	err := r.db.Update(func(txn *badger.Txn) error {
		userContextBytes, err := json.Marshal(userContext)
		if err != nil {
			return err
		}

		return txn.Set([]byte(userContext.UserID), userContextBytes)
	})

	return err
}

type UserContextMongoRepo struct {
	client         *mongo.Client
	dbName         string
	collectionName string
}

func NewUserContextMongoRepo(client *mongo.Client, dbName, collectionName string) (*UserContextMongoRepo, error) {
	return &UserContextMongoRepo{
		client:         client,
		dbName:         dbName,
		collectionName: collectionName,
	}, nil
}

func (r *UserContextMongoRepo) GetUserContext(userID string) (domain.UserContext, error) {
	var userContext domain.UserContext
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := r.client.Database(r.dbName).Collection(r.collectionName)
	err := collection.FindOne(ctx, bson.M{"userid": userID}).Decode(&userContext)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.UserContext{}, market_data_mcp_serverErr.UserContextNotFoundError{UserID: userID}
		}
		return domain.UserContext{}, err
	}

	return userContext, nil
}

func (r *UserContextMongoRepo) InsertUserContext(userContext domain.UserContext) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := r.client.Database(r.dbName).Collection(r.collectionName)
	_, err := collection.InsertOne(ctx, userContext)
	return err
}

func (r *UserContextMongoRepo) UpdateUserContext(userContext domain.UserContext) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": userContext}
	collection := r.client.Database(r.dbName).Collection(r.collectionName)
	res, err := collection.UpdateOne(ctx, bson.M{"userid": userContext.UserID}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return market_data_mcp_serverErr.UserContextNotFoundError{UserID: userContext.UserID}
	}
	return nil
}
