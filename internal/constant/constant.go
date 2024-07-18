package constant

import "time"

// Fields 从文章获取的字段
var Fields = [9]string{"source_appid", "countryName", "provinceName", "user_name", "createTime", "biz", "msg_title", "msg_link", "hd_head_img"}

const Domain = "mp.weixin.qq.com"

// TimeOut 文件下载超时时间
const TimeOut = 60 * time.Second

// CollectHomeUrl 首页合计 URL
const CollectHomeUrl = "https://mp.weixin.qq.com/mp/homepage"
const CollectAlbumUrl = "https://mp.weixin.qq.com/mp/appmsgalbum"

// AudioPrefix 音频资源前缀
const AudioPrefix = "https://res.wx.qq.com/voice/getvoice?mediaid="

// ExcludeFolder 定义排除目录
var ExcludeFolder = []string{".DS_Store", "css", "task"}
