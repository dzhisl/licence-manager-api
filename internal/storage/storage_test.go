package storage

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/dzhisl/license-api/pkg/config"
	"github.com/dzhisl/license-api/pkg/logger"
	"github.com/google/go-cmp/cmp"
)

var testCtx = context.Background()
var (
	user = User{
		Id:         1,
		DiscordId:  3,
		TelegramId: 120938,
		License: License{
			Key:            "initialKey",
			MaxActivations: 3,
			Devices: []Device{
				{HWID: "device_1"},
				{HWID: "device_2"},
			},
			IssuedAt:  Timestamp(1234567890),
			ExpiresAt: Timestamp(1234567890 + 30*24*3600),
			Status:    Frozen,
		},
		CreatedAt: Timestamp(234567890),
	}
)

func TestMain(m *testing.M) {
	config.InitConfig()
	logger.InitLogger()

	InitStorage(testCtx)
	code := m.Run()
	os.Exit(code)
}

func TestUserFlow(t *testing.T) {
	t.Run("CreateUser", func(t *testing.T) {
		err := connector.CreateUser(testCtx, user)
		if err != nil {
			t.Errorf("failed to create user in database: %v", err)
		}
	})

	t.Run("FindUser", func(t *testing.T) {
		testCases := []struct {
			name   string
			params GetUserParams
		}{
			{"Find by UserId", GetUserParams{UserId: user.Id}},
			{"Find by TelegramId", GetUserParams{TelegramId: user.TelegramId}},
			{"Find by DiscordId", GetUserParams{DiscordId: user.DiscordId}},
			{"Find by License", GetUserParams{License: user.License.Key}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				got, err := connector.GetUser(testCtx, tc.params)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if diff := cmp.Diff(&user, got); diff != "" {
					t.Errorf("User mismatch (-want +got):\n%v", diff)
				}
			})
		}
	})

	t.Run("UpdateDevices", func(t *testing.T) {
		err := connector.AddHwidSession(testCtx, user.Id, "device_3")
		if err != nil {
			t.Fatalf("failed to add hwid session: %v", err)
		}
		userObj, err := connector.GetUser(testCtx, GetUserParams{UserId: user.Id})
		if err != nil {
			t.Fatalf("failed to find user from database: %v", err)
		}
		if len(userObj.License.Devices) != 3 {
			spew.Dump(userObj.License.Devices)
			t.Fatalf("devices length doesn't match: want 3 got: %d", len(userObj.License.Devices))
		}

		if strings.Compare(userObj.License.Devices[2].HWID, "device_3") != 0 {
			t.Errorf("device name doesn't match: device_3 %s", userObj.License.Devices[2].HWID)
		}
	})

	t.Run("UpdateLimitDevices", func(t *testing.T) {
		// Attempt to add a 4th device, which should fail as MaxActivations is 3.
		err := connector.AddHwidSession(testCtx, user.Id, "device_4")
		if err == nil {
			t.Fatalf("added session, but should've had an error as limit was reached")
		}

	})

	t.Run("RemoveDevice", func(t *testing.T) {
		err := connector.DeleteHwidSession(testCtx, user.Id, "device_3")
		if err != nil {
			t.Fatalf("failed to delete hwid session: %v", err)
		}
		userObj, err := connector.GetUser(testCtx, GetUserParams{UserId: user.Id})
		if err != nil {
			t.Fatalf("failed to find user from database: %v", err)
		}
		if len(userObj.License.Devices) != 2 {
			spew.Dump(userObj.License.Devices)
			t.Fatalf("devices length doesn't match: want 2 got: %d", len(userObj.License.Devices))
		}

	})

	t.Run("ResetDevices", func(t *testing.T) {
		err := connector.ResetHwidSessions(testCtx, user.Id)
		if err != nil {
			t.Fatalf("failed to reset user's devices: %v", err)
		}

		user, err := connector.GetUser(testCtx, GetUserParams{UserId: user.Id})
		if err != nil {
			t.Fatalf("failed to find user from database: %v", err)
		}
		if len(user.License.Devices) != 0 {
			t.Fatalf("devices length doesn't match: want 0 got: %d", len(user.License.Devices))
		}
	})

	t.Run("UpdateLicenseStatus", func(t *testing.T) {
		newStatus := Active
		err := connector.ChangeLicenseStatus(testCtx, user.Id, newStatus)
		if err != nil {
			t.Fatalf("failed to change license status for user: %v", err)
		}
		user, err := connector.GetUser(testCtx, GetUserParams{UserId: user.Id})
		if err != nil {
			t.Fatalf("failed to find user from database: %v", err)
		}
		if user.License.Status != newStatus {
			t.Fatalf("license status doesn't match: want %s got: %s", newStatus, user.License.Status)
		}
	})

	t.Run("UpdateLicense", func(t *testing.T) {
		newLicense := License{
			Key:            "NewKey",
			MaxActivations: 10,
			IssuedAt:       Timestamp(1234567892),
			ExpiresAt:      Timestamp(1234567892 + 30*24*3600),
			Status:         Burned,
		}
		err := connector.UpdateLicense(testCtx, user.Id, newLicense)
		if err != nil {
			t.Fatalf("failed to update license for user: %v", err)
		}
		userObj, err := connector.GetUser(testCtx, GetUserParams{UserId: user.Id})
		if err != nil {
			t.Fatalf("failed to find user from database: %v", err)
		}

		if diff := cmp.Diff(newLicense, userObj.License); diff != "" {
			t.Errorf("License mismatch (-want +got):\n%v", diff)
		}
	})

	t.Run("UpdateHwidLimit", func(t *testing.T) {
		err := connector.UpdateHwidLimit(testCtx, user.Id, 5)
		if err != nil {
			t.Fatalf("failed to update hwid limit: %v", err)
		}
	})

	t.Run("RenewLicense", func(t *testing.T) {
		newTimestamp := Timestamp(9999999999)
		err := connector.RenewLicense(testCtx, user.Id, newTimestamp)
		if err != nil {
			t.Fatalf("failed to renew license: %v", err)
		}
	})

	t.Run("BindDiscord", func(t *testing.T) {
		newDiscordID := 98765
		err := connector.BindDiscord(testCtx, user.Id, newDiscordID)
		if err != nil {
			t.Fatalf("failed to bind discord for user: %v", err)
		}
		user, err := connector.GetUser(testCtx, GetUserParams{UserId: user.Id})
		if err != nil {
			t.Fatalf("failed to find user from database: %v", err)
		}
		if user.DiscordId != newDiscordID {
			t.Fatalf("discord id doesn't match: want %d got: %d", newDiscordID, user.DiscordId)
		}
	})

	t.Run("BindTelegram", func(t *testing.T) {
		newTelegramID := 12345
		err := connector.BindTelegram(testCtx, user.Id, newTelegramID)
		if err != nil {
			t.Fatalf("failed to bind telegram for user: %v", err)
		}
		user, err := connector.GetUser(testCtx, GetUserParams{UserId: user.Id})
		if err != nil {
			t.Fatalf("failed to find user from database: %v", err)
		}
		if user.TelegramId != newTelegramID {
			t.Fatalf("telegram id doesn't match: want %d got: %d", newTelegramID, user.TelegramId)
		}
	})

	t.Run("DeleteUser", func(t *testing.T) {
		count, err := connector.DeleteUser(testCtx, user.Id)
		if err != nil {
			t.Errorf("failed to delete user from database: %v", err)
			return
		}
		var want int64 = 1
		if count != want {
			t.Errorf("got %d, wanted %d", count, want)
		}

		// Verify user is deleted
		_, err = connector.GetUser(testCtx, GetUserParams{UserId: user.Id})
		if err == nil {
			t.Errorf("user should have been deleted, but was found")
		}
	})
}

func TestGetUserNotFound(t *testing.T) {
	// a user that is not in the database
	nonExistentUserID := 999
	_, err := connector.GetUser(testCtx, GetUserParams{UserId: nonExistentUserID})
	if err == nil {
		t.Errorf("expected an error when getting a non-existent user, but got nil")
	}
}

// bind a discord ID to user which is not in DB
func TestBindDiscordUserNotFound(t *testing.T) {
	nonExistentUserID := 999
	err := connector.BindDiscord(testCtx, nonExistentUserID, 2345678)
	if err == nil {
		t.Errorf("expected an error when setting a discord ID for a non-existent user, but got nil")
	}
}
