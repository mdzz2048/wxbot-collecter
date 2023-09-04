package utils

/* ------------------------ 配置文件 ------------------------ */

type Config struct {
	Cubox  CuboxConfig  `json:"cubox"`
	SiYuan SiYuanConfig `json:"siyuan"`
	WeChat WeChatConfig `json:"wechat"`
}

type CuboxConfig struct {
	CuboxUser     string `json:"cubox_user"`
	CuboxPassword string `json:"cubox_password"`
	CuboxToken    string `json:"cubox_token"`
}

type SiYuanConfig struct {
	SiYuanHost  string `json:"siyuan_host"`
	SiYuanToken string `json:"siyuan_token"`
}

type WeChatConfig struct {
	WeChatName string `json:"wechat_name"`
}

/* ------------------------ 数据库 ------------------------ */

type Article struct {
	ID          int64  `db:"id"`
	Title       string `db:"title"`
	Description string `db:"des"`
	URL         string `db:"url"`
	BlockID     string `db:"block_id"`
	CreateTime  string `db:"createTime"`
	UpdateTime  string `db:"updateTime"`
}
