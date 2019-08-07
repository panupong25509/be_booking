package repositories

import (
	"encoding/base64"
	"reflect"
	"unsafe"

	"github.com/labstack/echo"
	"github.com/panupong25509/be_booking_sign/db"
	"github.com/panupong25509/be_booking_sign/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(c echo.Context, data map[string]interface{}) (interface{}, interface{}) {
	db := db.DbManager()
	user := models.User{}
	if !user.CheckParams(data) {
		return nil, models.Error{400, "กรอกข้อมูลไม่ครบ"}
	}
	_, err := GetUserByUsername(c, data)
	if err == nil {
		return nil, models.Error{500, "Username นี้มีผู้ใช้แล้ว"}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(data["password"].(string)), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	_ = user.CreateModel(data, string(hash))
	db.NewRecord(user)
	db.Create(&user)
	return models.Success{200, "success"}, nil
}

func Login(c echo.Context, data map[string]interface{}) (interface{}, interface{}) {
	db := db.DbManager()
	username := data["username"].(string)
	password := data["password"].(string)
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
	if CheckPasswordHash(BytesToString(hashBytes), user.Password) {
		jwt := EncodeJWT(user)
		return models.JWT{jwt}, nil
	}
	return nil, models.Error{400, "username or password incorrect"}
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

func GetUserByUsername(c echo.Context, data map[string]interface{}) (interface{}, interface{}) {
	db := db.DbManager()
	if data["username"] == nil {
		return nil, models.Error{400, "ไม่มี username"}
	}
	username := data["username"].(string)
	user := models.Users{}
	db.First(&user, "username = (?)", username)
	if len(user) == 0 {
		return nil, models.Error{400, "ไม่มี username"}
	}
	return user[0], nil
}

func GetUserById(c echo.Context) (interface{}, interface{}) {
	jwtReq := c.Request().Header.Get("Authorization")
	tokens, err := DecodeJWT(jwtReq)
	if err != nil {
		return nil, err
	}
	db := db.DbManager()
	user := models.User{}
	err = db.Find(&user, tokens["UserID"])
	if err != nil {
		return nil, models.Error{400, "ไม่มีผู้ใช้นี้ใน database"}
	}
	return user, nil
}

// fileBytes, err := ioutil.ReadAll(f)
// tempFile.Write(fileBytes)
// in, err := os.Open(tempFile.Name())
// defer in.Close()
// _, file := filepath.Split(tempFile.Name())
// out, err := os.Create(`D:\fe_booking_sign\public\img\` + file)
// if _, err = io.Copy(out, in); err != nil {
// 	log.Print(err)
// 	return nil, models.Error{500, "can't add sign"}
// }
// return file, nil

// }
