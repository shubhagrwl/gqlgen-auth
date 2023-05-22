package util

import (
	"context"
	"strings"
	userDBmodel "todo/internal/app/db/dto/user"
	"todo/internal/app/service/graph/model"
	"todo/internal/app/service/logger"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

func ValidatePassword(ctx context.Context, password string, hashedPassword string) bool {
	log := logger.Logger(ctx)
	// Comparing the password with the hash
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Errorf("error while decoding hash", err)
		return false
	}
	return true
}

func GenerateHash(ctx context.Context, password string) (string, error) {
	log := logger.Logger(ctx)

	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("error while encoding hash", err)
		return "", err
	}
	return string(hash), nil
}

func SendMail(ctx context.Context, sub, name, email, templateID, otp string) error {
	log := logger.Logger(ctx)

	body := "YOUR OTP :" + otp
	subject := "todo: Welcome Onboard !!"

	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@orgosys.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("us2.smtp.mailhostbox.com", 587, "no-reply@orgosys.com", "NMvBLXt7")
	if err := d.DialAndSend(m); err != nil {
		log.Error(err)
		return err

	}
	log.Infof("Mail sent")

	return nil

	// log := logger.Logger(ctx)
	// from := mail.NewEmail("Shubham Agarwal", "admin@nomizo.io")
	// subject := sub
	// to := mail.NewEmail(name, email)

	// content := mail.NewContent("text/html", "I'm replacing the <strong>body tag</strong>")

	// m := mail.NewV3MailInit(from, subject, to, content)

	// m.Personalizations[0].DynamicTemplateData["token"] = otp
	// m.SetTemplateID(templateID)

	// client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	// response, err := client.Send(m)
	// if err != nil || response.StatusCode != http.StatusOK {
	// 	log.Errorf("unable to send email response : ", response, "err:", err)
	// 	return fmt.Errorf("unable to send")
	// }

	// log.Infof("Body : ", response.Body, "Headers :", response.Headers, "Status Code :", response.StatusCode)

	// log.Info("Mail sent")

	// return nil
}

func UserDetailToUser(user userDBmodel.UserDetail) *model.User {
	firstName := *user.FirstName
	lastName := *user.LastName
	fullName := strings.Join([]string{firstName, lastName}, " ")
	dob := user.Dob.String()
	return &model.User{
		ID:              user.ID.String(),
		FullName:        &fullName,
		DateOfBirth:     &dob,
		Email:           user.Email,
		IsEmailVerified: user.IsEmailVerified,
		// ImageURL:         user.Picture,
		// MobileNumber:     user.Mobile,
		// IsMobileVerified: user.IsMobileVerified,
		// Gender:           (*model.Gender)(user.Gender),
		LastLoginAt: user.LastLoginAt,
		Active:      user.Active,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

// Bool stores v in a new bool value and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int32 stores v in a new int value and returns a pointer to it.
func Int(v int) *int { return &v }

// Int32 stores v in a new int32 value and returns a pointer to it.
func Int32(v int32) *int32 { return &v }

// Int64 stores v in a new int64 value and returns a pointer to it.
func Int64(v int64) *int64 { return &v }

// Float32 stores v in a new float32 value and returns a pointer to it.
func Float32(v float32) *float32 { return &v }

// Float64 stores v in a new float64 value and returns a pointer to it.
func Float64(v float64) *float64 { return &v }

// Uint32 stores v in a new uint32 value and returns a pointer to it.
func Uint32(v uint32) *uint32 { return &v }

// Uint64 stores v in a new uint64 value and returns a pointer to it.
func Uint64(v uint64) *uint64 { return &v }

// String stores v in a new string value and returns a pointer to it.
func String(v string) *string { return &v }

// Array stores v in a new slice and returns a pointer to it.
func Array(v []string) *[]string { return &v }
