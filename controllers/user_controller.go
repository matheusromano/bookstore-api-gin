package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"gin-mongo-api/database"
	helper "gin-mongo-api/helpers"
	"gin-mongo-api/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.GetCollection(database.DB, "user")
var validate = validator.New()

//HashPassword is used to encrypt the password
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

//VerifyPassword checks
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or passowrd is incorrect")
		check = false
	}

	return check, msg
}

//sign up user
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "data": nil, "message": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,
				gin.H{"success": false, "data": nil, "message": "error occured while checking for the email"})

			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		if count > 0 {
			c.JSON(http.StatusInternalServerError,
				gin.H{"success": false, "data": nil, "message": "email ready exists"})

			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email)
		user.Token = &token
		user.Refresh_token = &refreshToken

		insertedId, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError,
				gin.H{"success": false, "data": nil, "message": msg})

			return
		}
		defer cancel()

		c.JSON(http.StatusOK, gin.H{"success": true, "data": insertedId, "message": "user signup sucess"})

	}
}

//Login user
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User
		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email)

		helper.UpdateAllTokens(token, refreshToken, *foundUser.Email)
		err = userCollection.FindOne(ctx, bson.M{"email": foundUser.Email}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{
			"token":         foundUser.Token,
			"refresh_token": foundUser.Refresh_token},
			"message": "return successfully"})

	}
}
