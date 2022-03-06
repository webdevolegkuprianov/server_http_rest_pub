package sqlstore

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	logger "github.com/webdevolegkuprianov/server_http_rest/app/apiserver/logger"
	"github.com/webdevolegkuprianov/server_http_rest/app/apiserver/model"
	"github.com/webdevolegkuprianov/server_http_rest/app/apiserver/store"
)

//User repository
type UserRepository struct {
	store *Store
}

type AccessDetails struct {
	UserId uint64
	Exp    uint64
}

//Find jwt email password (create token)
func (r *UserRepository) FindUser(email string, password string) (*model.User1, error) {
	u := &model.User1{}
	if err := r.store.dbPostgres.QueryRow(context.Background(),
		"SELECT id, email, password FROM users WHERE email = $1 AND password = $2",
		email, password).Scan(&u.ID, &u.Email, &u.Password); err != nil {
		if err == sql.ErrNoRows {
			logger.ErrorLogger.Println(err)
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}

//Find jwt user id (verify token)
func (r *UserRepository) FindUserid(userid uint64) error {
	u := &model.User2{}

	if err := r.store.dbPostgres.QueryRow(context.Background(),
		"SELECT id FROM users WHERE id = $1",
		userid).Scan(&u.UserID); err != nil {
		if err == sql.ErrNoRows {
			logger.ErrorLogger.Println(err)
			return store.ErrRecordNotFound
		}

		return err
	}
	return nil
}

//creating token
//create token
func (r *UserRepository) CreateToken(userid uint64, config *model.Service) (string, time.Time, error) {
	var err error

	datetime_exp_unix := time.Now().Add(time.Hour * 24 * time.Duration(config.Spec.Jwt.LifeTerm)).Unix()
	datetime_exp := time.Unix(datetime_exp_unix, 0)
	t := new(time.Time)

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = datetime_exp_unix
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(config.Spec.Jwt.TokenDecode))
	if err != nil {
		logger.ErrorLogger.Println(err)
		return "", *t, err
	}

	return token, datetime_exp, nil
}

//extract token from header
func (r *UserRepository) ExtractToken(req *http.Request) string {
	bearToken := req.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return bearToken
}

//verify token
func (r *UserRepository) VerifyToken(req *http.Request, config *model.Service) (*jwt.Token, error) {
	tokenString := r.ExtractToken(req)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Spec.Jwt.TokenDecode), nil
	})
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}
	return token, nil
}

//extract data from token
func (r *UserRepository) ExtractTokenMetadata(req *http.Request, config *model.Service) (*model.AccessDetails, error) {

	//var accessDetails model.AccessDetails

	token, err := r.VerifyToken(req, config)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return nil, err
		}
		return &model.AccessDetails{
			UserId: userId,
		}, nil
	}
	return nil, err
}
