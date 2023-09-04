package database

import (
	"database/sql"
	"fmt"
	"testing"
	"wxbot/cubox"
	"wxbot/utils"
)

func TestIsDatabaseExist(t *testing.T) {
	dbPath := "wxbot.db"
	isExist := IsDatabaseExist(dbPath)

	if !isExist {
		t.Errorf("数据库路径: %s 不存在", dbPath)
	}
}

func TestInitDatabase(t *testing.T) {
	InitDatabase()
}

func TestConnectDatabase(t *testing.T) {
	db, err := ConnectDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}

	isExist := IsTableExist(db, "article")
	if isExist {
		fmt.Println("数据库正常")
	}
}

func TestIsTableExist(t *testing.T) {
	dbPath := "wxbot.db"
	tableName := "article"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("连接数据库失败: ", err)
		return
	}
	defer db.Close()

	isExist := IsTableExist(db, tableName)
	if !isExist {
		t.Errorf("数据库表 %s 不存在", tableName)
	}
}

func TestArticleAdd(t *testing.T) {
	bookmark := cubox.BookMark{}
	err := utils.JsonUnmarshal([]byte(utils.CuboxData_BookMark), &bookmark)
	if err != nil {
		t.Errorf(err.Error())
	}

	dbPath := "wxbot.db"
	db, _ := sql.Open("sqlite3", dbPath)
	article := ConvertBookMarkToArticle(&bookmark)
	ArticleAdd(db, &article)
}

func TestArticleDelete(t *testing.T) {
	dbPath := "wxbot.db"
	db, _ := sql.Open("sqlite3", dbPath)
	ArticleDelete(db, "id", "2")
}

func TestArticleSearch(t *testing.T) {
	dbPath := "wxbot.db"
	db, _ := sql.Open("sqlite3", dbPath)
	articles := ArticleSearch(db, "title", "url", "https://cloud.tencent.com/developer/article/1849807")

	fmt.Printf("共搜索到: %d 条结果", len(articles))
	for _, article := range articles {
		fmt.Println("Title: ", article.Title)
	}
}

func TestArticleUpdate(t *testing.T) {
	dbPath := "wxbot.db"
	db, _ := sql.Open("sqlite3", dbPath)
	ArticleUpdate(db, "title", "新的标题", "url", "https://cloud.tencent.com/developer/article/1849807")
}

func TestConvertBookMarkToArticle(t *testing.T) {
	bookmark := cubox.BookMark{}
	err := utils.JsonUnmarshal([]byte(utils.CuboxData_BookMark), &bookmark)
	if err != nil {
		t.Errorf(err.Error())
	}
	article := ConvertBookMarkToArticle(&bookmark)

	fmt.Println(article)
}
