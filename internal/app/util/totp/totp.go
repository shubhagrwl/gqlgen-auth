package totp

import (
	"context"
	"encoding/base32"
	"time"
	"todo/internal/app/constants"
	"todo/internal/app/service/logger"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/viper"
)

// Provider describes a Provider which can generate and verify a otp
type ITOTP interface {
	GetCode(ctx context.Context, identifier string, period *int) (string, error)
	VerifyCode(ctx context.Context, identifier, code string, period *int) (bool, error)
}

// TOTP is used to create custom time based otp
type TOTP struct {
	secretStr       string
	periodInSeconds uint
	digits          otp.Digits
	algorithm       otp.Algorithm
	godMode         bool
}

// New creates an instance on TotpProvides and returns it.
// default value of dur is 60, if the passed in value is less than 1 second
func New(ctx context.Context, projectEnv constants.ProjectEnvironment) *TOTP {
	log := logger.Logger(ctx)

	dur := time.Duration(viper.GetInt(constants.OTPDURATION)) * time.Second
	secret := viper.GetString(constants.OTPSECRETKEY)

	periodInSeconds := uint(dur.Seconds())

	if periodInSeconds == 0 {
		periodInSeconds = 60
	}

	var godMode bool

	if projectEnv == constants.ProjectEnvDevelopment {
		godMode = true
	}

	log.Info("Successfully TOTP started...")

	return &TOTP{
		secretStr:       secret,
		periodInSeconds: periodInSeconds,
		digits:          otp.DigitsSix,
		algorithm:       otp.AlgorithmSHA256,
		godMode:         godMode,
	}
}

func (tp *TOTP) opts() totp.ValidateOpts {
	return totp.ValidateOpts{
		Period:    tp.periodInSeconds,
		Skew:      1,
		Digits:    tp.digits,
		Algorithm: tp.algorithm,
	}
}

func (tp *TOTP) secret(id string) string {
	return base32.StdEncoding.EncodeToString([]byte(id + tp.secretStr))
}

// GetCode gets the TOTP code that expires after configured time
func (tp *TOTP) GetCode(ctx context.Context, identifier string, period *int) (string, error) {
	if period != nil {
		period := time.Duration(*period)
		tp.periodInSeconds = uint(period)
	}
	s, err := totp.GenerateCodeCustom(tp.secret(identifier), time.Now(), tp.opts())
	return s, err
}

// VerifyCode validates the received OTP code
func (tp *TOTP) VerifyCode(ctx context.Context, identifier, code string, period *int) (bool, error) {
	if period != nil {
		period := time.Duration(*period)
		tp.periodInSeconds = uint(period)
	}

	if tp.godMode {
		if code == "111111" {
			return true, nil
		}
	}

	return totp.ValidateCustom(code, tp.secret(identifier), time.Now(), tp.opts())
}
