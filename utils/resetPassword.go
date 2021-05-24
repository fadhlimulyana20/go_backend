package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/fadhlimulyana20/go_backend/config"
	"github.com/fadhlimulyana20/go_backend/utils/constant"
	"github.com/google/uuid"
)

type ResetPassword struct {
	Url string
}

func (r *ResetPassword) Init(userId uint) error {
	// Create new random token
	token, _ := uuid.NewRandom()

	// Create redis Connection
	rc := &config.RedisConfig{}
	rc.Init()
	rdb := rc.GetConnection()

	// Set a key value of rest password token
	key := constant.ResetPasswordPrefix + token.String()
	h, _ := time.ParseDuration("48h")
	err := rdb.Set(context.Background(), key, userId, h).Err()

	if err != nil {
		return err
	}

	r.Url = fmt.Sprintf("http://127.0.0.1:3000/user/reset_password/%s", token.String())
	return nil
}
