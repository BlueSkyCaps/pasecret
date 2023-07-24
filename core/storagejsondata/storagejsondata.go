package storagejson

// GetRelatedDataByCid 根据归类文件夹Id检索它的密码项
func GetRelatedDataByCid(ci Category) *[]Data {
	// new返回指向[]Data的指针
	var relatedData = new([]Data)
	for _, d := range AppRef.LoadedItems.Data {
		if d.CategoryId == ci.Id {
			*relatedData = append(*relatedData, d)
		}
	}
	return relatedData
}
