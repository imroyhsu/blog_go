package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

type LoginController struct {
	beego.Controller
}

type LoginModel struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (this *LoginController) Post() {
	var ob LoginModel
	//map用来表示json格式的数据
	//后面可以考虑封装下result
	result := make(map[string]interface{})

	//使用这个来接收json格式
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &ob); err != nil {
		fmt.Println(err)
		//考虑修改http状态码为3系列
		result["code"] = "003"
		result["msg"] = "fail"
		result["data"] = "登录失败，参数错误"
	} else {
		isLogin := queryLoginValid(ob.Name, ob.Password)
		if isLogin {
			//如果没开启session此处是会报错的
			this.SetSession("isLogin", true)
			result["code"] = "000"
			result["msg"] = "success"
			result["data"] = "登录成功"
		} else {
			result["code"] = "001"
			result["msg"] = "fail"
			result["data"] = "登录失败，重新登录"
		}
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
	}
	this.Data["json"] = string(bytes)
	this.ServeJSON()
}

func queryLoginValid(username string, password string) bool {

	dbhost := beego.AppConfig.String("dbhost")
	dbport := beego.AppConfig.String("dbport")
	dbuser := beego.AppConfig.String("dbuser")
	dbpassword := beego.AppConfig.String("dbpassword")
	dbname := beego.AppConfig.String("dbname")
	dbcharset := beego.AppConfig.String("dbcharset")

	var isLogin bool = false
	db, err := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?"+dbcharset)
	if (err != nil) {
		fmt.Println(err)
	}

	//直接开启mysql数据库去查询
	//后面考虑用orm实现
	var sql string = "SELECT name,password FROM user where name = " + "\"" + username + "\"" + "and password = " + "\"" + password + "\""
	fmt.Println(sql)
	rows, err := db.Query(sql)
	if (err != nil) {
		fmt.Println(err)
	}

	for rows.Next() {
		isLogin = true
	}

	rows.Close()
	db.Close()
	return isLogin
}
