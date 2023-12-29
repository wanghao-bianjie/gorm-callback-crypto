package example

import (
	"reflect"
	"strings"
	"testing"

	"github.com/wanghao-bianjie/gorm-callback-crypto/callback"
	"github.com/wanghao-bianjie/gorm-callback-crypto/example/model"
	"github.com/wanghao-bianjie/gorm-callback-crypto/example/repository"
	"github.com/wanghao-bianjie/gorm-callback-crypto/util/aes"
	"gorm.io/gorm"
)

type UserNew struct {
	model.User
	//Id          uint   `gorm:"column:id;primaryKey"`
	//Name        string `gorm:"column:name"`
	//PhoneNumber string `gorm:"column:phone_number"`
	//Address     string `gorm:"column:address"`
	//IdNo        string `gorm:"column:id_no"`
	//UpdateAt    int64  `gorm:"column:update_at"`
	Uid   uint
	UidNo string
}

func CustomizeBeforeFn(str string) (string, error) {
	return "AES[" + str + "]", nil
}

func CustomizeAfterFn(str string) (string, error) {
	if strings.HasPrefix(str, "AES[") {
		str = strings.ReplaceAll(str, "AES[", "")
		str = str[:len(str)-1]
	}
	return str, nil
}

func setUpDB() {
	repository.InitMysqlDB("root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local&time_zone=%27Asia%2FShanghai%27")
	repository.CreateTable()
}

func setUpCallback() {
	setUpCallbackWithAesKey(repository.GetDb(), aesKey)
	//setUpCallbackWithCustomizeFn(repository.GetDb())
}

func setUpCallbackWithAesKey(db *gorm.DB, aesKey string) {
	err := callback.Register(db, []callback.ICryptoModel{
		new(model.User),
	},
		callback.WithDefaultAesFnKey([]byte(aesKey)),
	)
	if err != nil {
		panic(err)
	}
}

func setUpCallbackWithCustomizeFn(db *gorm.DB) {
	err := callback.Register(db, []callback.ICryptoModel{
		new(model.User),
	},
		callback.WithBeforeHandleFn(CustomizeBeforeFn),
		callback.WithAfterHandleFn(CustomizeAfterFn),
	)
	if err != nil {
		panic(err)
	}
}

var user = model.User{
	Id:          0,
	Name:        "1229-tom",
	PhoneNumber: "15651850001",
	Address:     "江苏省南京市建邺区xxxx广场",
	IdNo:        "320322199801010001",
	UpdateAt:    1677662129,
}

var aesKey = "1234567890123456"

func TestBatchCreate(t *testing.T) {
	setUpDB()
	db := repository.GetDb()
	setUpCallback()
	user1 := user
	user2 := user
	//users := []model.User{user1, user2}
	users := []*model.User{&user1, &user2}
	err := db.Save(users).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log(users[0].Id)
	t.Log(users[1].Id)
}

func TestCreate(t *testing.T) {
	setUpDB()
	db := repository.GetDb()
	setUpCallback()
	db = db.Debug()
	//通过结构体指针创建
	var user1 = user
	var err error
	err = db.Create(&user1).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("通过结构体指针创建:", user1)

	//通过map创建
	userMap := map[string]interface{}{
		"name":         user.Name,
		"phone_number": user.PhoneNumber,
		"address":      user.Address,
		"id_no":        user.IdNo,
		"update_at":    user.UpdateAt,
	}
	u1 := model.User{}
	err = db.Model(&u1).Create(userMap).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("通过map创建:", u1)
	t.Log("通过map创建:", userMap)

	u2 := model.User{}
	err = db.Model(&u2).Create(&userMap).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("通过*map创建:", u2)
	t.Log("通过*map创建:", userMap)

	//通过结构体指针批量创建
	var users = []model.User{user, user}
	err = db.Create(users).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("通过结构体指针批量创建:", users)

	//通过map批量创建
	var userMaps = []map[string]interface{}{
		{
			"name":         user.Name,
			"phone_number": user.PhoneNumber,
			"address":      user.Address,
			"id_no":        user.IdNo,
			"update_at":    user.UpdateAt,
		},
		{
			"name":         user.Name,
			"phone_number": user.PhoneNumber,
			"address":      user.Address,
			"id_no":        user.IdNo,
			"update_at":    user.UpdateAt,
		},
	}
	u3 := model.User{}
	err = db.Model(&u3).Create(userMaps).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("通过map批量创建:", u3)
	t.Log("通过map批量创建:", userMaps)

	u4 := model.User{}
	err = db.Model(&u4).Create(&userMaps).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("通过*[]map批量创建:", u4)
	t.Log("通过*[]map批量创建:", userMaps)

	var user2 = user
	err = db.Save(&user2).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("通过结构体指针创建(save):", user2)

	var users2 = []model.User{user, user}
	err = db.Save(users2).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("通过结构体指针批量创建(save):", users2)

	//exec 不会调用 Create() 的 callback 方法，而是自己的 exec 的callback
	/*err = db.Exec("INSERT INTO `user` (`name`,`phone_number`,`address`,`id_no`,`update_at`) VALUES ('tom','15651850001','江苏省南京市建邺区鼓楼创新广场','320322199801010001',1677641264)").Error
	if err != nil {
		t.Fatal(err)
	}*/
}

func TestFind(t *testing.T) {
	setUpDB()
	db := repository.GetDb()
	setUpCallback()
	//db = db.Debug()
	var res []model.User
	var err error

	//err = db.Session(&gorm.Session{SkipHooks: true}).Model(&model.User{}).Where("phone_number = ?", user.PhoneNumber).Find(&res).Error
	//err = db.Session(&gorm.Session{SkipHooks: true}).Model(&model.User{}).Where(User{PhoneNumber: user.PhoneNumber}).Find(&res).Error
	//err = db.Session(&gorm.Session{SkipHooks: true}).Model(&model.User{}).Where(map[string]interface{}{
	//	"phone_number": user.PhoneNumber,
	//}).Find(&res).Error

	//禁用钩子，查询数据库的真实数据
	err = db.Session(&gorm.Session{SkipHooks: true}).Model(&model.User{}).Find(&res).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("skip callback res:")
	for _, re := range res {
		t.Log(re)
	}

	//不禁用钩子，默认使用 callback 查询处理后的数据
	err = db.Model(&model.User{}).Find(&res).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("use callback res:")
	for _, re := range res {
		t.Log(re)
	}

	//使用 map 接收
	var resMap []map[string]interface{}
	err = db.Model(&model.User{}).Find(&resMap).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("use callback res(map):")
	for _, re := range resMap {
		t.Log(re)
	}

	//使用 []string 接收
	var idNos []string
	err = db.Model(&model.User{}).Select("id_no").Find(&idNos).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("use callback res([]string):")
	for _, re := range idNos {
		t.Log(re)
	}

	//使用自定义结构体接收
	var resJoin []UserNew
	err = db.Model(&model.User{}).Select(`
	u.id as uid,
	u.id_no as uid_no,
	user.id,
	user.name,
	user.phone_number,
	user.address,
	user.id_no,
	user.update_at
	`).Joins("left join user as u on user.id = u.id").Find(&resJoin).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("use callback res([]struct join):")
	for _, re := range resJoin {
		t.Log(re)
	}

	//条件查询
	encrypt, err := aes.CBCPKCS7EncryptToBase64([]byte("15651850001"), []byte(aesKey))
	if err != nil {
		t.Fatal(err)
	}
	err = db.Model(&model.User{}).Where("phone_number = ?", encrypt).Find(&res).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("use callback res(filter with encrypt):")
	for _, re := range res {
		t.Log(re)
	}

	err = db.Model(&model.User{}).Where("phone_number = ?", "15651850001").Find(&res).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("use callback res(filter without encrypt):")
	for _, re := range res {
		t.Log(re)
	}
}

func TestName(t *testing.T) {
	setUpDB()
	db := repository.GetDb()
	setUpCallback()
	db = db.Debug()
	user1 := user
	err := db.Create(&user1).Error
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Where(&model.User{Id: user1.Id}).Updates(model.User{
		PhoneNumber: "123",
		IdNo:        "qwe",
	}).Error; err != nil {
		t.Fatal(err)
	}
	var res model.User
	err = db.Last(&res).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestUpdate(t *testing.T) {
	setUpDB()
	db := repository.GetDb()
	setUpCallback()
	db = db.Debug()

	var users []model.User
	for i := 0; i < 11; i++ {
		users = append(users, user)
	}
	err := db.Create(users).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("批量创建:")
	for _, u := range users {
		t.Log(u)
	}

	//save
	user1 := users[0]
	user1.Name = "new1-" + user.Name
	user1.PhoneNumber = "new1-" + user.PhoneNumber
	err = db.Save(&user1).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("save user1:")
	t.Log(user1)

	//update
	user2 := users[1]
	err = db.Model(&user2).Update("phone_number", "new2-"+user2.PhoneNumber).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("update user2:")
	t.Log(user2)

	//updates
	user3 := users[2]
	u3 := model.User{
		Name:        "new3-p-sp-" + user.Name,
		PhoneNumber: "new3-p-sp-" + user.PhoneNumber,
	}
	err = db.Model(&user3).Updates(&u3).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("updates user3-p-sp:")
	t.Log(user3)
	t.Log(u3)

	user4 := users[3]
	u4 := model.User{
		Id:          user4.Id,
		Name:        "new4-t-sp-" + user.Name,
		PhoneNumber: "new4-t-sp-" + user.PhoneNumber,
	}
	err = db.Table(user4.TableName()).Updates(&u4).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("updates user4-t-sp:")
	t.Log(user4)
	t.Log(u4)

	user5 := users[4]
	user5.Name = "new5-p-" + user.Name
	user5.PhoneNumber = "new5-p-" + user.PhoneNumber
	err = db.Updates(&user5).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("updates user5-p:")
	t.Log(user5)

	//updates map
	user6 := users[5]
	m6 := map[string]interface{}{
		"name":         "new6-p-m-" + user.Name,
		"phone_number": "new6-p-m-" + user.PhoneNumber,
	}
	err = db.Model(&user6).Updates(m6).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("updates user6-p-m:")
	t.Log(user6)
	t.Log(m6)

	user7 := users[6]
	m7 := map[string]interface{}{
		"name":         "new7-p-mp-" + user.Name,
		"phone_number": "new7-p-mp-" + user.PhoneNumber,
	}
	err = db.Model(&user7).Updates(&m7).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("updates user7-p-mp:")
	t.Log(user7)
	t.Log(m7)

	user8 := users[7]
	m8 := map[string]interface{}{
		"name":         "new8-t-m-" + user.Name,
		"phone_number": "new8-t-m-" + user.PhoneNumber,
	}
	err = db.Table(user8.TableName()).Where(model.User{Id: user8.Id}).Updates(m8).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("updates user8-t-m:")
	t.Log(user8)
	t.Log(m8)

	user9 := users[8]
	m9 := map[string]interface{}{
		"name":         "new9-t-mp-" + user.Name,
		"phone_number": "new9-t-mp-" + user.PhoneNumber,
	}
	err = db.Table(user9.TableName()).Where(model.User{Id: user9.Id}).Updates(&m9).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("updates user9-t-mp:")
	t.Log(user9)
	t.Log(m9)

	user10 := users[9]
	u10 := model.User{
		Id:          user10.Id,
		Name:        "new10-t-s-" + user.Name,
		PhoneNumber: "new10-t-s-" + user.PhoneNumber,
	}
	err = db.Table(user10.TableName()).Updates(u10).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("updates user10-t-s:")
	t.Log(user10)
	t.Log(u10)

	user11 := users[10]
	u11 := model.User{
		Name:        "new11-p-s-" + user.Name,
		PhoneNumber: "new11-p-s-" + user.PhoneNumber,
	}
	err = db.Model(&user11).Updates(u11).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log("updates user11-p-s:")
	t.Log(user11)
	t.Log(u11)
}

func TestReflect(t *testing.T) {
	var v interface{}
	v = user

	reflectValue := reflect.ValueOf(v)
	if !reflectValue.CanSet() {
		newStruct := reflect.New(reflectValue.Type())
		//for i := 0; i < reflectValue.Type().NumField(); i++ {
		//	newStruct.Elem().Field(i).Set(reflectValue.Field(i))
		//}
		newStruct.Elem().Set(reflectValue)
		t.Log(newStruct.Elem().Interface())
	}
}
