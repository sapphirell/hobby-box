package model

func FlushCache(cacheKey string) {

}

var (
	cachePre         = "box_"
	ItemListCacheKey = cachePre + "item_list;type:%s;page:%d,uid:%d"
)
