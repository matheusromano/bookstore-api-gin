package controllers

import (
	"context"
	//"encoding/json"
	"gin-mongo-api/database"
	"gin-mongo-api/models"
	"gin-mongo-api/responses"
	"net/http"
	"time"

	//"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"

	//"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/* init authentication with jwt
var (
	users = map[string]string{
		"user1": "password1",
		"user2": "password2",
	}

	// Create the JWT key used to create the signature
	JwtKey = []byte("my_key")
)

// Create a struct to read the username and password from request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT
// We add jwt.StardartClaims as embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// authentication string
type AuthResult struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

 end of jwt authentication */

var bookCollection *mongo.Collection = database.GetCollection(database.DB, "books")

//var validate = validator.New()

func CreateBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		//mdw.Authenticate()

		/* init jwt token auth
		tknStr := c.Request.Header["token"][0]
		if tknStr == "" {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims := helpers.SignedDetails{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return helpers.ValidateToken(tknStr). nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.Writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			c.Writer.WriteHeader(http.StatusForbidden)
		}
		// end jwt token auth */

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var book models.Book
		defer cancel()

		// validate the request body
		if err := c.BindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, responses.BookResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
		}

		// use the validator libraty to validate required fields
		if validationErr := validate.Struct(&book); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.BookResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data:    map[string]interface{}{"data": validationErr.Error()},
			})
			return
		}

		newBook := models.Book{
			//Id:     primitive.NewObjectID(),
			Title:  book.Title,
			Author: book.Author,
			Price:  book.Price,
		}

		result, err := bookCollection.InsertOne(ctx, newBook)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.BookResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}

		c.JSON(http.StatusCreated, responses.BookResponse{
			Status:  http.StatusCreated,
			Message: "success",
			Data:    map[string]interface{}{"data": result},
		})
	}

}
func GetABook() gin.HandlerFunc {
	return func(c *gin.Context) {
		/* init jwt token auth
		tknStr := c.Request.Header["token"][0]
		if tknStr == "" {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.Writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			c.Writer.WriteHeader(http.StatusForbidden)
		}
		// end jwt token auth*/

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		bookId := c.Param("bookId")
		var book models.Book
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(bookId)

		err := bookCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&book)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.BookResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}

		c.JSON(http.StatusOK, responses.BookResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    map[string]interface{}{"data": book},
		})
	}
}

func EditABook() gin.HandlerFunc {
	return func(c *gin.Context) {
		/* init jwt token auth
		tknStr := c.Request.Header["Private-Token"][0]
		if tknStr == "" {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.Writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			c.Writer.WriteHeader(http.StatusForbidden)
		}*/
		// end jwt token auth

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		bookId := c.Param("bookId")
		var book models.Book
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(bookId)

		//validate the request body
		if err := c.BindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, responses.BookResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&book); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.BookResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data:    map[string]interface{}{"data": validationErr.Error()},
			})
			return
		}

		update := bson.M{"title": book.Title, "author": book.Author, "price": book.Price}
		result, err := bookCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.BookResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}

		//get updated book datails
		var updatedBook models.Book
		if result.MatchedCount == 1 {
			err := bookCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedBook)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.BookResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data:    map[string]interface{}{"data": err.Error()},
				})
				return
			}

			c.JSON(http.StatusOK, responses.BookResponse{
				Status:  http.StatusOK,
				Message: "success",
				Data:    map[string]interface{}{"data": updatedBook},
			})
		}
	}
}

func DeleteABook() gin.HandlerFunc {
	return func(c *gin.Context) {
		/* init jwt token auth
		tknStr := c.Request.Header["Private-Token"][0]
		if tknStr == "" {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.Writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			c.Writer.WriteHeader(http.StatusForbidden)
		}
		// end jwt token auth */

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		bookId := c.Param("bookId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(bookId)

		result, err := bookCollection.DeleteOne(ctx, bson.M{"_id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.BookResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, responses.BookResponse{
				Status:  http.StatusNotFound,
				Message: "error",
				Data:    map[string]interface{}{"data": "Book with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusNoContent, responses.BookResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    map[string]interface{}{"data": "Book successfully deleted!"}},
		)

	}
}

func GetAllBooks() gin.HandlerFunc {
	return func(c *gin.Context) {
		/* init jwt token auth
		tknStr := c.Request.Header["Private-Token"][0]
		if tknStr == "" {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.Writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			c.Writer.WriteHeader(http.StatusForbidden)
		}*/
		// end jwt token auth

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var books []models.Book
		defer cancel()

		results, err := bookCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.BookResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}

		// reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleBook models.Book
			if err = results.Decode(&singleBook); err != nil {
				c.JSON(http.StatusInternalServerError, responses.BookResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data:    map[string]interface{}{"data": err.Error()},
				})
				return
			}

			books = append(books, singleBook)
		}

		c.JSON(http.StatusOK, responses.BookResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    map[string]interface{}{"data": books},
		})
	}
}

/*func Signin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds Credentials
		err := json.NewDecoder(c.Request.Body).Decode(&creds)
		// if the structure of the body is wrong, return an HTTP error
		if err != nil {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Get the expected password from our in memory map
		expectedPassword, ok := users[creds.Username]

		// if a password exists for the given user
		// AND, if it is the same as the password we received, the we can move ahead
		// if NOT, then we return an "unauthorized" status
		if !ok || expectedPassword != creds.Password {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Declare the expiration time of the token
		// here, we have kept it as 5 minutes
		expirationTime := time.Now().Add(5 * time.Minute)

		// Create the JWT claims, which includes the username and expiry time
		claims := &Claims{
			Username: creds.Username,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
			},
		}

		// Declare the token with the algorithm used for signing, and the claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Create the JWT string
		tokenString, err := token.SignedString(JwtKey)
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Finally, we set the client cookie for "token" as the JWT we just generated
		// we also set an expiry time which is the same as the token itself
		authResult := &AuthResult{
			Success: true,
			Token:   tokenString,
		}
		_ = json.NewEncoder(c.Writer).Encode(authResult)

	}

}
*/
