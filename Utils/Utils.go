package Utils

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
	"net/smtp"
	"os"
	"strconv"
)

func SendEmailToAdmin(subject string, message string) error {
	msg := []byte("To: mahendrakrs448@gmail.com\r\n" + "Subject:" + subject + "\r\n" + "\r\n" + message + "\r\n")

	auth := smtp.PlainAuth(
		"",
		os.Getenv("SMTP_EMAIL"),
		os.Getenv("SMTP_PASSWORD"),
		"smtp.gmail.com")

	if err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		os.Getenv("SMTP_EMAIL"),
		[]string{"mahendrakrs448@gmail.com"},
		msg,
	); err != nil {
		return err
	}

	return nil
}

func SendEmail(to string, subject string, message string) error {
	msg := []byte("To:" + to + "\r\n" + "Subject:" + subject + "\r\n" + "\r\n" + message + "\r\n")

	auth := smtp.PlainAuth(
		"",
		os.Getenv("SMTP_EMAIL"),
		os.Getenv("SMTP_PASSWORD"),
		"smtp.gmail.com")

	if err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		os.Getenv("SMTP_EMAIL"),
		[]string{to},
		msg,
	); err != nil {
		return err
	}

	return nil
}

func MonthToRoman(month int) string {
	months := map[int]string{
		1:  "I",
		2:  "II",
		3:  "III",
		4:  "IV",
		5:  "V",
		6:  "VI",
		7:  "VII",
		8:  "VIII",
		9:  "IX",
		10: "X",
		11: "XI",
		12: "XII",
	}

	roman, exists := months[month]
	if !exists {
		return ""
	}

	return roman
}

func Paginate(ctx *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		if page <= 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		dbClone := db.Session(&gorm.Session{})
		var total int64
		dbClone.Count(&total)

		totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

		ctx.Set("total_pages", totalPages)
		ctx.Set("page_size", pageSize)
		ctx.Set("page", page)
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
