package repository

import (
	"context"
	"event-catch/config"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	config *config.Config

	client *mongo.Client
	db     *mongo.Database

	Tx            *mongo.Collection
	NFT           *mongo.Collection
	NFTCollection *mongo.Collection
}

func NewRepository(config *config.Config) (*Repository, error) {
	r := Repository{
		config: config,
	}

	var err error
	ctx := context.Background()

	mongoConf := config.Database.Mongo

	if r.client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoConf.DataSource)); err != nil {
		return nil, err
	} else if err = r.client.Ping(ctx, nil); err != nil {
		r.db = r.client.Database(mongoConf.DB)

		r.Tx = r.db.Collection(mongoConf.Tx)
		r.NFT = r.db.Collection(mongoConf.NFT)
		r.NFTCollection = r.db.Collection(mongoConf.NFTCollection)
	}

	return &r, nil
}

func (r *Repository) UpsertTxEvent(from, to, sender common.Address, tokenID *big.Int, tx string) error {
	// Transfer 로직만 거르니 하드코딩으로 작성 인터페이스 사용 요망
	opt := options.Update().SetUpsert(true)

	filter := bson.M{
		"tx": tx,
	}
	update := bson.M{"$set": bson.M{
		"sender":  hexutil.Encode(sender[:]),
		"from":    hexutil.Encode(from[:]),
		"to":      hexutil.Encode(to[:]),
		"tokenID": tokenID.Int64(),
	}}

	_, err := r.Tx.UpdateOne(context.Background(), filter, update, opt)

	return err
}

func (r *Repository) UpsertTransferEvent(tokenID *big.Int, to common.Address) error {
	opt := options.Update().SetUpsert(true)

	filter := bson.M{
		"tokenID": tokenID.Int64(),
	}

	update := bson.M{"$set": bson.M{
		"owner": hexutil.Encode(to[:]),
	}}

	_, err := r.NFT.UpdateOne(context.Background(), filter, update, opt)

	return err
}
