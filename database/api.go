package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"wxbot/cubox"
	"wxbot/utils"

	_ "github.com/mattn/go-sqlite3"
)

var dbPath = utils.DB_PATH

/* ------------------------ Database ------------------------ */

func InitDatabase() {
	// 检查数据库是否存在
	if IsDatabaseExist(dbPath) {
		fmt.Println("数据库文件已存在，检查数据库表是否存在")
	}
	fmt.Println("数据库文件不存在，初始化数据库……")

	// 连接到SQLite数据库
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("连接数据库失败: ", err)
		return
	}
	defer db.Close()

	// 创建 Cubox 表、Article 表
	createTable(db, "article", CreateTableArticle)
	createTable(db, "cubox", CreateTableCubox)
}

func ConnectDatabase() (*sql.DB, error) {
	// 检查数据库是否存在
	if !IsDatabaseExist(dbPath) {
		fmt.Println("数据库文件不存在，请先初始化数据库")
		return nil, errors.New("database is not exist")
	}

	// 连接到SQLite数据库
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("连接数据库失败: ", err)
		return nil, err
	}
	return db, nil
}

func IsDatabaseExist(dbPath string) bool {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func IsTableExist(db *sql.DB, tableName string) bool {
	query := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s'", tableName)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("查询 %s 表失败: %s", tableName, err)
		return false
	}
	defer rows.Close()

	return rows.Next()
}

func createTable(db *sql.DB, tableName string, statement string) {
	if IsTableExist(db, tableName) {
		fmt.Printf("表 %s 已存在", tableName)
		return
	}

	_, err := db.Exec(statement)
	if err != nil {
		log.Fatalf("表 %s 创建失败: %s\n", tableName, err)
		return
	}
	fmt.Printf("表 %s 创建成功\n", tableName)
}

/* ------------------------ Article ------------------------ */

func ArticleAdd(db *sql.DB, data *utils.Article) any {
	sql := InsertArticle
	stmt, err := db.Prepare(sql)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer stmt.Close()

	title := data.Title
	des := data.Description
	url := data.URL
	block_id := ""
	createTime := data.CreateTime
	updateTime := data.UpdateTime

	result, err := stmt.Exec(title, des, url, block_id, createTime, updateTime)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("成功插入 %d 行数据\n", rowsAffected)
	return data
}

func ArticleDelete(db *sql.DB, key string, value string) any {
	sql := fmt.Sprintf(DeleteArticle, key)
	stmt, err := db.Prepare(sql)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	result, err := stmt.Exec(value)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("成功删除 %d 行数据", rowsAffected)
	return result
}

func ArticleSearch(db *sql.DB, fieldName string, key string, value string) []utils.Article {
	sql := fmt.Sprintf(QueryArticle, fieldName, key)
	fmt.Println(sql)
	stmt, err := db.Prepare(sql)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	result, err := stmt.Query(value)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer result.Close()

	// todo: 这里解析还需要完善
	articles := []utils.Article{}
	for result.Next() {
		article := utils.Article{}
		err := result.Scan(
			// &article.ID,
			&article.Title,
			// &article.Description,
			// &article.URL,
			// &article.BlockID,
			// &article.CreateTime,
			// &article.UpdateTime,
		)

		if err != nil {
			fmt.Println("Failed to sacn row: ", err)
			return nil
		}

		articles = append(articles, article)
	}
	// rowsAffected, _ := result.Columns()
	// fmt.Printf(rowsAffected[])
	return articles
}

func ArticleUpdate(db *sql.DB, fieldName string, fieldValue string, key string, value string) any {
	sql := fmt.Sprintf(UpdateArticle, fieldName, key)
	stmt, err := db.Prepare(sql)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	result, err := stmt.Exec(fieldValue, value)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("成功更新 %d 行数据", rowsAffected)
	return result
}

/* ------------------------ Message ------------------------ */

func MessageAdd(data any) any {
	return data
}

func MessageDelete(data any) any {
	return data
}

func MessageSearch(data any) any {
	return data
}

func MessageUpdate(data any) any {
	return data
}

/* ------------------------ Cubox ------------------------ */

func CuboxAdd(data any) any {
	return data
}

func CuboxDelete(data any) any {
	return data
}

func CuboxSearch(data any) any {
	return data
}

func CuboxUpdate(data any) any {
	return data
}

func ConvertBookMarkToArticle(bookmark *cubox.BookMark) utils.Article {
	article := utils.Article{
		Title:       bookmark.Title,
		Description: bookmark.Description,
		URL:         bookmark.TargetURL,
		BlockID:     "",
		CreateTime:  bookmark.CreateTime,
		UpdateTime:  bookmark.UpdateTime,
	}
	return article
}
