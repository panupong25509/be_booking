package repositories

import (
	"encoding/base64"
	"log"
	"reflect"
	"unsafe"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo"
	"github.com/panupong25509/be_booking_sign/db"
	"github.com/panupong25509/be_booking_sign/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(c echo.Context) (interface{}, interface{}) {
	db := db.DbManager()
	user := models.User{}
	if !user.CheckParams(c) {
		return nil, models.Error{400, "กรอกข้อมูลไม่ครบ"}
	}
	_, err := GetUserByUsername(c)
	if err == nil {
		return nil, models.Error{500, "Username นี้มีผู้ใช้แล้ว"}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(c.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	_ = user.CreateModel(c, string(hash))
	db.NewRecord(user)
	db.Create(&user)
	return models.Success{200, "success"}, nil
}

func Login(c echo.Context) (interface{}, interface{}) {
	db := db.DbManager()
	username := c.FormValue("username")
	password := c.FormValue("password")
	if username == "" {
		return nil, models.Error{400, "ไม่มี username"}
	}
	if password == "" {
		return nil, models.Error{400, "ไม่มี password"}
	}
	hashBytes, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		return nil, err
	}
	user := models.User{}
	db.Where("username = (?)", username).Find(&user)
	if !CheckPasswordHash(BytesToString(hashBytes), user.Password) {
		return nil, models.Error{400, "username or password incorrect"}
	}
	jwt := EncodeJWT(user)
	return models.JWT{jwt}, nil
}

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserByUsername(c echo.Context) (interface{}, interface{}) {
	db := db.DbManager()
	if c.FormValue("username") == "" {
		return nil, models.Error{400, "ไม่มี username"}
	}
	username := c.FormValue("username")
	user := models.Users{}
	db.Where("username = (?)", username).First(&user)
	if len(user) == 0 {
		return nil, models.Error{400, "ไม่มี username"}
	}
	return user[0], nil
}

func GetUserById(c echo.Context) (interface{}, interface{}) {
	jwtReq, err := GetJWT(c)
	if err != nil {
		return nil, err
	}
	tokens, err := DecodeJWT(jwtReq.(string))
	if err != nil {
		return nil, err
	}
	log.Print(tokens["UserID"])
	db := db.DbManager()
	user := models.User{}
	db.Where("id like (?)", tokens["UserID"]).Find(&user)
	return user, nil
}

func GetUserByIduuid(c echo.Context, id uuid.UUID) (interface{}, interface{}) {
	db := db.DbManager()
	user := models.User{}
	db.Where("id like (?)", id).Find(&user)
	return user, nil
}
