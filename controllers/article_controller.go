package controllers

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type ArticleController struct {
	beego.Controller
}

type Article struct {
	ArticleId string   `json:"articleId"` //首字母大写表示public，小写表示private，添加tag，使其输出时变成小写
	Title     string   `json:"title"`
	Date      string   `json:"date"`
	Content   string   `json:"content"`
	Gist      string   `json:"gist"`
	Labels    []string `json:"labels"`
}

func (this *ArticleController) Get() {
	fmt.Println("进入articlecontroller的get方法")

	isLogin := this.GetSession("isLogin")
	fmt.Println("准备打印isLogin....")
	fmt.Println(isLogin)

	result := make(map[string]interface{})
	if isLogin == nil {
		result["code"] = "001"
		result["msg"] = "fail"
		result["data"] = "登录失效，请重新登录"

	} else {
		result["code"] = "000"
		result["msg"] = "success"
		var articleList []Article = queryArticleList()
		result["data"] = articleList
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
	}

	this.Data["json"] = string(bytes)
	this.ServeJSON()
}

func (this *ArticleController) Post() {
	isLogin := this.GetSession("isLogin")
	fmt.Println("保存文章接口...")
	fmt.Println(isLogin)

	result := make(map[string]interface{})
	if isLogin == nil {
		result["code"] = "001"
		result["msg"] = "fail"
		result["data"] = "登录失效，请重新登录"

	} else {
		var ob Article
		json.Unmarshal(this.Ctx.Input.RequestBody, &ob)
		fmt.Println(ob)

		fmt.Println("fuck you...")
		fmt.Println(ob.Labels)
		var isSaveSuccess bool = saveArticle(ob.Title, ob.Date, ob.Content, ob.Gist, ob.Labels)
		result["msg"] = "success"
		if isSaveSuccess {
			result["code"] = "000"
			result["data"] = "文章保存成功"
		} else {
			result["code"] = "002"
			result["data"] = "文章保存失败"
		}
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
	}
	this.Data["json"] = string(bytes)
	this.ServeJSON()
}

func (this *ArticleController) Put() {
	isLogin := this.GetSession("isLogin")
	fmt.Println("准备打印isLogin....")
	fmt.Println(isLogin)

	result := make(map[string]interface{})
	if isLogin == nil {
		result["code"] = "001"
		result["msg"] = "fail"
		result["data"] = "登录失效，请重新登录"

	} else {
		var ob Article
		json.Unmarshal(this.Ctx.Input.RequestBody, &ob)
		labelsString := ""
		for i := 0; i < len(ob.Labels); i++ {
			labelsString += "," + ob.Labels[i]
		}
		labelsString = string(labelsString[1:])
		fmt.Println(labelsString)
		var isUpdateSuccess bool = updateArticle(ob.ArticleId, ob.Title, ob.Date, ob.Content, ob.Gist, labelsString)
		result["msg"] = "success"
		if isUpdateSuccess {
			result["code"] = "000"
			result["data"] = "文章更新成功"
		} else {
			result["code"] = "002"
			result["data"] = "文章更新失败"
		}
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
	}
	this.Data["json"] = string(bytes)
	this.ServeJSON()
}

func (this *ArticleController) Delete() {
	isLogin := this.GetSession("isLogin")
	fmt.Println("准备打印isLogin....")
	fmt.Println(isLogin)

	result := make(map[string]interface{})
	if isLogin == nil {
		result["code"] = "001"
		result["msg"] = "fail"
		result["data"] = "登录失效，请重新登录"

	} else {

		var isDeleteSuccess bool = deleteArticle(this.GetString("articleId"))
		fmt.Println("删除成功没", isDeleteSuccess)
		if isDeleteSuccess {
			result["code"] = "000"
			result["msg"] = "success"
			result["data"] = "文章删除成功"
		} else {
			result["code"] = "002"
			result["msg"] = "faile"
			result["data"] = "文章删除失败"
		}

	}
	bytes, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
	}
	this.Data["json"] = string(bytes)
	this.ServeJSON()
}

// 删除文章
func deleteArticle(articleId string) bool {
	var dbhost string = beego.AppConfig.String("dbhost")
	var dbport string = beego.AppConfig.String("dbport")
	var dbuser string = beego.AppConfig.String("dbuser")
	var dbpassword string = beego.AppConfig.String("dbpassword")
	var dbname string = beego.AppConfig.String("dbname")
	var dbcharset string = beego.AppConfig.String("dbcharset")

	var isDeleteSuccess bool = false
	db, err := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?"+dbcharset)
	if err != nil {
		fmt.Println(err)
	}
	var sql string = "delete FROM article where articleId = \"" + articleId + "\""
	fmt.Println(sql)
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println(err)
	} else {
		isDeleteSuccess = true
	}

	rows.Close()
	db.Close()
	return isDeleteSuccess
}

//查询文章列表
func updateArticle(articleId string, title string, date string, content string, gist string, lables string) bool {
	var dbhost string = beego.AppConfig.String("dbhost")
	var dbport string = beego.AppConfig.String("dbport")
	var dbuser string = beego.AppConfig.String("dbuser")
	var dbpassword string = beego.AppConfig.String("dbpassword")
	var dbname string = beego.AppConfig.String("dbname")
	var dbcharset string = beego.AppConfig.String("dbcharset")

	var isUpdateSuccess bool = false

	db, err := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?"+dbcharset)
	if err != nil {
		fmt.Println(err)
	}

	content = strings.Replace(content, `"`, `\"`, -1)
	content = strings.Replace(content, `'`, `\'`, -1)
	content = strings.Replace(content, `\`, `\\\`, -1)
	var sql string = "update article set title = " + "\"" + title + "\", content = " + "\"" + content + "\", date = " + "\"" + date + "\", gist = " + "\"" + gist + "\", labels = " + "\"" + lables + "\" where articleId = \"" + articleId + "\""
	fmt.Println(sql)
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {

		isUpdateSuccess = true
	}
	rows.Close()
	db.Close()
	return isUpdateSuccess
}
func (this *ArticleController) GetDetails() {

	fmt.Println("进入articlecontroller的getdetails方法")

	isLogin := this.GetSession("isLogin")
	fmt.Println("准备打印isLogin....")
	fmt.Println(isLogin)

	result := make(map[string]interface{})
	if isLogin == nil {
		result["code"] = "001"
		result["msg"] = "fail"
		result["data"] = "登录失效，请重新登录"

	} else {
		var articleDetail Article = queryArticleDetail(this.GetString("articleId"))
		result["code"] = "000"
		result["msg"] = "success"
		result["data"] = articleDetail
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
	}
	this.Data["json"] = string(bytes)
	this.ServeJSON()
}

//查询文章列表
func queryArticleList() []Article {
	var dbhost string = beego.AppConfig.String("dbhost")
	var dbport string = beego.AppConfig.String("dbport")
	var dbuser string = beego.AppConfig.String("dbuser")
	var dbpassword string = beego.AppConfig.String("dbpassword")
	var dbname string = beego.AppConfig.String("dbname")
	var dbcharset string = beego.AppConfig.String("dbcharset")

	db, err := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?"+dbcharset)
	if err != nil {
		return nil
	}

	sqlstr := "SELECT articleId,title,date,gist,labels FROM article"
	fmt.Println(sqlstr)
	rows, err := db.Query(sqlstr)
	if err != nil {
		return nil
	}

	var articleResults []Article
	for rows.Next() {

		var article Article
		var lableStr string

		rows.Columns()
		err = rows.Scan(&article.ArticleId, &article.Title, &article.Date, &article.Gist, &lableStr)
		if err != nil {
			return nil
		}

		//删除像框框一样的字符，还不能写在注释里很奇怪
		var pureLabelString string = strings.Replace(lableStr, "\u0000", "", -1)
		var labels [] string = strings.Split(pureLabelString, ",")
		article.Labels = labels
		fmt.Println(article)
		articleResults = append(articleResults, article)
	}
	rows.Close()
	db.Close()
	return articleResults
}

//查询文章详情
func queryArticleDetail(articleId string) Article {
	var dbhost string = beego.AppConfig.String("dbhost")
	var dbport string = beego.AppConfig.String("dbport")
	var dbuser string = beego.AppConfig.String("dbuser")
	var dbpassword string = beego.AppConfig.String("dbpassword")
	var dbname string = beego.AppConfig.String("dbname")
	var dbcharset string = beego.AppConfig.String("dbcharset")

	var article Article
	db, err := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?"+dbcharset)
	if err != nil {
		return article
	}

	var sql string = "SELECT articleId,title,content,date,gist,labels FROM article where articleId = \"" + articleId + "\""
	fmt.Println(sql)
	rows, err := db.Query(sql)
	if err != nil {
		return article
	}

	var labelStr string
	for rows.Next() {

		rows.Columns()
		err = rows.Scan(&article.ArticleId, &article.Title, &article.Content, &article.Date, &article.Gist, &labelStr)
		if err != nil {
			return article
		}

	}
	var content string = article.Content
	content = strings.Replace(content, `\"`, `"`, -1)
	content = strings.Replace(content, `\'`, `'`, -1)
	content = strings.Replace(content, `\\\`, `\`, -1)
	article.Content = content

	var pureLabelString string = strings.Replace(labelStr, "\u0000", "", -1)
	var labels [] string = strings.Split(pureLabelString, ",")
	article.Labels = labels

	rows.Close()
	db.Close()
	return article
}

//保存文章
func saveArticle(title string, date string, content string, gist string, labels []string) bool {
	var dbhost string = beego.AppConfig.String("dbhost")
	var dbport string = beego.AppConfig.String("dbport")
	var dbuser string = beego.AppConfig.String("dbuser")
	var dbpassword string = beego.AppConfig.String("dbpassword")
	var dbname string = beego.AppConfig.String("dbname")
	var dbcharset string = beego.AppConfig.String("dbcharset")

	fmt.Println(len(labels))
	var labelsString string = ""
	for i := 0; i < len(labels); i++ {
		labelsString += "," + labels[i]
	}
	rs := []rune(labelsString)
	labelsString = string(rs[1:len(labelsString)])

	fmt.Println(22222)
	var isSaveSuccess bool = false
	db, err := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?"+dbcharset)
	if (err != nil) {
		fmt.Println(err)
	}

	h := md5.New()
	h.Write([]byte(date))
	if err != nil {

	}
	fmt.Println(3333)
	// articleId用日期md5加密后的数据
	var articleId string = hex.EncodeToString(h.Sum(nil))
	content = strings.Replace(content, `"`, `\"`, -1)
	content = strings.Replace(content, `'`, `\'`, -1)
	content = strings.Replace(content, `\`, `\\\`, -1)
	fmt.Println(4444)
	var sql string = "insert into article(articleId,title,date,content,gist,labels) values(" + "\"" + articleId + "\",\"" + title + "\",\"" + date + "\",\"" + content + "\",\"" + gist + "\",\"" + labelsString + "\"" + ")"
	fmt.Println(sql)
	rows, err := db.Query(sql)
	if (err != nil) {
		fmt.Println(err)
	}

	for rows.Next() {

		isSaveSuccess = true
	}
	rows.Close()
	db.Close()
	return isSaveSuccess

}
