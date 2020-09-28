package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/chonlatee/authserver/pkg/models/mariadb"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type route struct {
	userModel *mariadb.UserModel
}

func (r *route) register(c echo.Context) error {

	var password string

	type UserRegister struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var u UserRegister

	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if len(u.Username) < 6 {
		return c.JSON(http.StatusBadRequest, struct {
			Message string `json:"message"`
		}{"Invalid username length"})
	}

	if len(u.Password) < 6 {
		return c.JSON(http.StatusBadRequest, struct {
			Message string `json:"message"`
		}{"Invalid password length"})
	}

	if len(u.Password) > 50 {
		return c.JSON(http.StatusBadRequest, struct {
			Message string `json:"message"`
		}{"Invalid password length"})
	}

	if !isEmailValid(u.Email) {
		return c.JSON(http.StatusBadRequest, struct {
			Message string `json:"message"`
		}{"Invalid email"})
	}

	password, err := hashPassword(u.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, struct {
			Message string `json:"message"`
		}{"Password is invalid"})
	}

	ctx := c.Request().Context()

	_, err = r.userModel.Insert(ctx, u.Username, u.Email, password)

	if err != nil {
		return c.JSON(http.StatusOK, struct {
			Message string `json:"message"`
		}{fmt.Sprintf("Register error %v", err.Error())})
	}

	return c.JSON(http.StatusOK, struct {
		Message string `json:"message"`
	}{"Register success"})
}

func (r *route) login(c echo.Context) error {

	type Userlogin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var uLogin Userlogin

	if err := c.Bind(&uLogin); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()

	hashPassword, err := r.userModel.GetByEmail(ctx, uLogin.Email)

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(uLogin.Password))

	if err != nil {
		return c.JSON(http.StatusUnauthorized, struct {
			Message string `json:"message"`
		}{"Invalid credentials"})
	}

	token := jwt.New(jwt.GetSigningMethod("RS256"))

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "chonlatee"
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	k, err := ioutil.ReadFile("key/private.pem")

	if err != nil {
		log.Fatalln("Read key error")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(k)

	if err != nil {
		log.Fatalf("Parse rsa private key from pem error: %v", err.Error())
	}

	t, err := token.SignedString(signKey)
	if err != nil {
		log.Fatalf("Sign jwt error: %v", err.Error())
	}

	return c.JSON(http.StatusOK, struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{"Login success", t})
}

func (r *route) getPublic(c echo.Context) error {
	f, err := ioutil.ReadFile("key/public.pem")

	if err != nil {
		return c.JSON(http.StatusBadGateway, struct {
			Message string `json:"message"`
		}{"Can't get public key"})
	}

	return c.JSON(http.StatusOK, struct {
		Message   string `json:"message"`
		PublicKey string `json:"publickey"`
	}{"Get public key success", string(f)})
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}

	if !emailRegex.MatchString(e) {
		return false
	}

	// parts := strings.Split(e, "@")
	// mx, err := net.LookupMX(parts[1])
	// if err != nil || len(mx) == 0 {
	// 	return false
	// }

	return true
}

func hashPassword(password string) (string, error) {
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(passwordHashed), err
}

func openDB() *sql.DB {
	db, err := sql.Open("mysql", "authapp:123456@tcp(172.17.0.2:3306)/authapp")

	log.Println(db, err)
	if err != nil {
		log.Fatalf("Connect db error: %v\n", err)
	}

	return db

}
