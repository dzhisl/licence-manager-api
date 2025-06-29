package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/dzhisl/license-api/pkg/config"
	"github.com/dzhisl/license-api/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

var (
	databaseName   = "license-manager"
	collectionName = "users"
)

var (
	connector Connector
)

type Connector struct {
	userCollection *mongo.Collection
}

func GetConnector() Connector {
	return connector
}

func InitStorage(ctx context.Context) {
	client, err := mongo.Connect(options.Client().ApplyURI(config.AppConfig.MongoHost))
	if err != nil {
		logger.Fatal(ctx, "failed to connect to mongoDB", zap.Error(err))
	}

	userColl := client.Database(databaseName).Collection(collectionName)
	for i := 0; i < 3; i++ {
		err = userColl.Database().Client().Ping(context.Background(), nil)
		if err != nil {
			logger.Warn(ctx, "failed to ping mongoDB", zap.Error(err))
			time.Sleep(5 * time.Second)
			continue
		}
	}
	if err != nil {
		logger.Fatal(ctx, "failed to ping mongoDB after 3 attempts")
	}

	connector.userCollection = userColl

	logger.Info(ctx, "connected to MONGO DB")
}

func (c *Connector) CreateUser(ctx context.Context, u User) error {
	_, err := c.userCollection.InsertOne(ctx, u)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connector) DeleteUser(ctx context.Context, userId int) (deletedCount int64, err error) {
	filter := bson.M{"_id": userId}
	res, err := c.userCollection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

func (c *Connector) GetUser(ctx context.Context, params GetUserParams) (user *User, err error) {

	var filter primitive.M

	switch {
	case params.UserId != 0:
		filter = bson.M{"_id": params.UserId}
	case params.TelegramId != 0:
		filter = bson.M{"telegramId": params.TelegramId}
	case params.DiscordId != 0:
		filter = bson.M{"discordId": params.DiscordId}
	case params.License != "":
		filter = bson.M{"license.key": params.License}
	default:
		return nil, fmt.Errorf("at least one param must be provided")
	}

	var u *User

	err = c.userCollection.FindOne(ctx, filter).Decode(&u)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("record for user wasn't found")
		}
		return nil, err
	}
	return u, nil
}

func (c *Connector) GetAllUsers(ctx context.Context) (user []*User, err error) {
	var u []*User

	cursor, err := c.userCollection.Find(ctx, nil)

	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &u); err != nil {
		return nil, fmt.Errorf("failed to unpack users to struct:%w", err)
	}
	return u, nil
}

func (c *Connector) AddHwidSession(ctx context.Context, userId int, hwid string) error {
	filter := bson.M{"_id": userId}

	user, err := c.GetUser(ctx, GetUserParams{UserId: userId})
	if err != nil {
		return err
	}
	if len(user.License.Devices) >= user.License.MaxActivations {
		return fmt.Errorf("user have maximum allowed activations: %d", user.License.MaxActivations)
	}

	user.License.Devices = append(user.License.Devices, hwid)
	update := bson.M{"$set": bson.M{"license.devices": user.License.Devices}}
	res, err := c.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user devices: %w", err)
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func (c *Connector) DeleteHwidSession(ctx context.Context, userId int, hwid string) error {
	filter := bson.M{"_id": userId}

	user, err := c.GetUser(ctx, GetUserParams{UserId: userId})
	if err != nil {
		return err
	}

	var newHwidSessions []string
	for _, session := range user.License.Devices {
		if session != hwid {
			newHwidSessions = append(newHwidSessions, session)
		}
	}
	update := bson.M{"$set": bson.M{"license.devices": newHwidSessions}}
	res, err := c.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user devices: %w", err)
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func (c *Connector) ResetHwidSessions(ctx context.Context, userId int) error {
	filter := bson.M{"_id": userId}
	update := bson.M{"$set": bson.M{"license.devices": nil}}
	res, err := c.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to reset user devices: %w", err)
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func (c *Connector) ChangeLicenseStatus(ctx context.Context, userId int, status LicenseStatus) error {
	filter := bson.M{"_id": userId}
	update := bson.M{"$set": bson.M{"license.status": status}}
	res, err := c.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update license status: %w", err)
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func (c *Connector) UpdateLicense(ctx context.Context, userId int, license License) error {
	filter := bson.M{"_id": userId}
	update := bson.M{"$set": bson.M{"license": license}}
	res, err := c.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update license: %w", err)
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func (c *Connector) UpdateHwidLimit(ctx context.Context, userId, newLimit int) error {
	filter := bson.M{"_id": userId}
	update := bson.M{"$set": bson.M{"license.maxActivations": newLimit}}
	res, err := c.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update license hwid limits: %w", err)
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func (c *Connector) RenewLicense(ctx context.Context, userId int, expiresAt Timestamp) error {
	filter := bson.M{"_id": userId}
	update := bson.M{"$set": bson.M{"license.expiresAt": expiresAt}}
	res, err := c.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to renew license: %w", err)
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func (c *Connector) BindDiscord(ctx context.Context, userId, discordId int) error {
	filter := bson.M{"_id": userId}
	update := bson.M{"$set": bson.M{"discordId": discordId}}
	res, err := c.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to bind discord for user: %w", err)
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func (c *Connector) BindTelegram(ctx context.Context, userId, telegramId int) error {
	filter := bson.M{"_id": userId}
	update := bson.M{"$set": bson.M{"telegramId": telegramId}}
	res, err := c.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to bind telegram for user: %w", err)
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}
