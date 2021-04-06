package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/fadhlimulyana20/go_backend/config"
	"github.com/fadhlimulyana20/go_backend/utils/constant"
	"github.com/google/uuid"
)

type ConfirmationUrl struct{}

var ctx = context.Background()

func (c *ConfirmationUrl) Create(userId uint) (string, error) {
	// Create new random token
	token, _ := uuid.NewRandom()

	// Create redis Connection
	rc := &config.RedisConfig{}
	rc.Init()
	rdb := rc.GetConnection()

	// Set a key value of confirmation token
	key := constant.ConfirmUserPrefix + token.String()
	h, _ := time.ParseDuration("48h")
	err := rdb.Set(ctx, key, userId, h).Err()

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("http://127.0.0.1:3000/user/confirmation/%s", token.String())
	return url, nil

}
