package customType

type ApiConfType struct {

	UrlPath string		//url path
	ParamsArr []string	//参数数组
	ExpireTime int		//过期时间
	CheckCount int		//检查访问倍数
}

type RespCacheType struct {
	ApiUrl string		//原始url
	ParamsStr string	//请求参数
	Resp string			//返回结果
	Method string		//请求方式
}