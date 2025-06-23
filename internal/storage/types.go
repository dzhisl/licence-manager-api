package storage

type LicenseStatus string
type Timestamp int

const (
	Frozen LicenseStatus = "frozen"
	Active LicenseStatus = "active"
	Burned LicenseStatus = "burned"
)

type User struct {
	Id         int       `bson:"_id" json:"id"`
	TelegramId int       `bson:"telegramId" json:"telegramId"`
	DiscordId  int       `bson:"discordId" json:"discordId"`
	License    License   `bson:"license" json:"license"`
	CreatedAt  Timestamp `bson:"createdAt" json:"createdAt"`
}

type License struct {
	Key            string        `bson:"key" json:"key"`
	MaxActivations int           `bson:"maxActivations" json:"maxActivations"`
	Devices        []Device      `bson:"devices" json:"devices"`
	IssuedAt       Timestamp     `bson:"issuedAt" json:"issuedAt"`
	ExpiresAt      Timestamp     `bson:"expiresAt" json:"expiresAt"`
	Status         LicenseStatus `bson:"status" json:"status"`
}

type Device struct {
	HWID string `bson:"hwid" json:"hwid"`
}

type GetUserParams struct {
	UserId     int
	TelegramId int
	DiscordId  int
	License    string
}
