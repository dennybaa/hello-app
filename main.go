package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserReq struct defines a user
type UserReq struct {
	UserName    string `json:"username"`
	DateOfBirth string `json:"dateofbirth"`
}

// UserModel struct defines a user model
type UserModel struct {
	UserName    string    `json:"username"`
	DateOfBirth time.Time `json:"dateofbirth"`
}

// DaysToBirthday function evaluates days left to your birthday
// Born during a leap year on 29 Feb, are admitted to have birthday on 1 March
// during non-leap years.
func DaysToBirthday(birthday, today time.Time) (int, error) {
	var daysToBirthday int

	if today.Before(birthday) {
		return -1, errors.New("Not born yet :(")
	}

	y, m, d := today.Date()
	todayDate := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	_, m, d = birthday.Date()
	nextBirthdayDate := time.Date(y, m, d, 23, 59, 0, 0, time.UTC)

	// time duration, time until the upcoming birthday
	if nextBirthdayDate.Before(todayDate) {
		nextBirthdayDate = time.Date(y+1, m, d, 23, 59, 0, 0, time.UTC)
		daysToBirthday = int(nextBirthdayDate.Sub(todayDate).Hours() / 24)
	} else {
		nextBirthdayDate = time.Date(y, m, d, 23, 59, 0, 0, time.UTC)
		daysToBirthday = int(nextBirthdayDate.Sub(today).Hours() / 24)
	}

	return daysToBirthday, nil
}

func (user UserModel) helloMessage() (string, error) {
	var msg string
	days, err := DaysToBirthday(user.DateOfBirth, time.Now())

	if err != nil {
		return "", err
	} else if days == 0 {
		msg = fmt.Sprintf("Hello, %s! Happy birthday!", user.UserName)
	} else {
		msg = fmt.Sprintf("Hello, %s! Your birthday is in %d day(s)", user.UserName, days)
	}

	return msg, nil
}

// putHelloUsername creates or updates a user
func putHelloUsername(app *App) gin.HandlerFunc {
	db, client := app.databaseName, app.client

	return func(c *gin.Context) {
		var json UserReq
		var user UserModel

		userName := c.Param("username")

		err := c.BindJSON(&json)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		// We can't create or update without datetime provided
		dateOfBirth, err := time.Parse("2006-01-02", json.DateOfBirth)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid dateofbirth datetime format YYYY-MM-DD expected"},
			)
			return
		}

		// find and update
		col := client.Database(db).Collection("users")
		opts := options.FindOneAndUpdate()
		opts.Upsert = new(bool)
		*opts.Upsert = true

		res := col.FindOneAndUpdate(
			context.Background(),
			bson.M{"username": userName},
			bson.M{
				"$set": bson.M{"username": userName, "dateofbirth": dateOfBirth},
			},
			opts,
		)

		if err := res.Decode(&user); err != nil && err != mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// respond with 204 if OK
		c.JSON(http.StatusNoContent, gin.H{})
	}
}

// getHelloUsername outputs hello message containing amount of days to wait for the birthday
func getHelloUsername(app *App) gin.HandlerFunc {
	db, client := app.databaseName, app.client

	return func(c *gin.Context) {
		userName := c.Param("username")

		col := client.Database(db).Collection("users")
		res := col.FindOne(
			context.Background(),
			bson.M{"username": userName},
		)

		var user UserModel

		// user not found
		if err := res.Decode(&user); err == mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{})
			return
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// days to birthday (hello message)
		msg, err := user.helloMessage()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"message": msg})
	}
}

// Start the appliction
func main() {
	app := new(App)
	app.initConf()
	app.dbConnect()
	app.serve()
}
