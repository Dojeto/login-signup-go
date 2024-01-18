package controllers

import (
	"fmt"
	"go-jwt/initializers"
	"go-jwt/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(ctx *gin.Context) {

	//get the data from body
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// fmt.Println(body)
	// ctx.Bind(&body)
	if err := ctx.BindJSON(&body); err != nil {
		fmt.Println("Error binding JSON:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	if body.Email == "" || body.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "email and password are required fields",
		})
		return
	}

	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// create a user in database

	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed To Create User",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User Created",
	})
}

func Login(ctx *gin.Context) {
	// get the email and pass off req body

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		fmt.Println("Error binding JSON:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email",
		})
		return
	}

	passwordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if passwordErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 36).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		fmt.Println("Error creating JWT:", err.Error()) // Print the actual error
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "something", // Consider providing a more informative error message
		})
		return
	}

	// Now you can use the 'tokenString' for further processing.
	fmt.Println("Generated JWT:", tokenString)

	ctx.SetSameSite(http.SameSiteLaxMode)

	ctx.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})

}

func Validate(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
