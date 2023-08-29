package storagedata

type LoadedItems struct {
	GlobalConfig struct {
		SyncBranch int `json:"sync_branch"`
	} `json:"global_config"`
	Category []Category `json:"category"`
	Data     []Data     `json:"data"`
}
type Category struct {
	Name        string `json:"name"`
	Alias       string `json:"alias"`
	Description string `json:"description"`
	Id          string `json:"id"`
	Rank        int    `json:"rank"`
	Removable   bool   `json:"removable"`
	Renameable  bool   `json:"renameable"`
}
type Data struct {
	Name        string `json:"name"`
	AccountName string `json:"account_name"`
	Password    string `json:"password"`
	Site        string `json:"site"`
	Remark      string `json:"remark"`
	Id          int    `json:"id"`
	CategoryId  string `json:"category_id"`
}
