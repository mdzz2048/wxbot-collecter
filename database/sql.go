package database

/* ------------------------ Article ------------------------ */

var CreateTableArticle = `
	CREATE TABLE IF NOT EXISTS article (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		des TEXT,
		url TEXT,
		block_id TEXT,
		createTime INTEGER,
		updateTime INTEGER
	);
`
var InsertArticle = "INSERT INTO article (title, des, url, block_id, createTime, updateTime) VALUES (?, ?, ?, ?, ?, ?)"
var QueryArticle = "SELECT %s FROM article WHERE %s = ?"
var UpdateArticle = "UPDATE article SET %s = ? WHERE %s = ?"
var DeleteArticle = "DELETE FROM article WHERE %s = ?"

/* ------------------------ Cubox ------------------------ */

var CreateTableCubox = `
	CREATE TABLE IF NOT EXISTS cubox (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT,
		bookmark_id TEXT,
		data TEXT
	);
`
var InsertCubox = "INSERT INTO cubox (url, bookmark_id, data) VALUES (?, ?, ?)"
var InsertCuboxID = "INSERT INTO cubox (id) VALUES (?)"
var InsertCuboxURL = "INSERT INTO cubox (url) VALUES (?)"
var InsertCuboxData = "INSERT INTO cubox (data) VALUES (?)"
var QueryCuboxDataByID = "SELECT data FROM cubox WHERE id = ?"
var QueryCuboxDataByURL = "SELECT data FROM cubox WHERE url = ?"
