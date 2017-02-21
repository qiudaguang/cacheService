# cacheService
api接口缓存服务

config.ini

		[server]
			port = 9090  #服务监听(访问)端口
		[redis]  #redis配置
			host = 192.168.1.188
			port = 6379
			index = 4
		[api]  #真实提供服务地址
			host = qiudaguang.api.zhigoubao.com
			port = 80
		[cache_service]
			expiretime = 3600    #默认缓存有效期
			updatecycle = 10000  #每访问xxx次执行更新

apiConfig.json  #具体接口配置

		[{
			"UrlPath": "/model/getList",   #接口地址
			"ParamsArr": [                 #接口参数
				"model_style_id",
				"category_ename",
				"attribute_ename_arr"
			],
			"ExpireTime": 3600,     #接口有效期(若为0，则取config.ini中配置)
			"CheckCount": 2000      #每访问xxx次执行更新(若为0，则取config.ini中配置)
		},
		{
			"UrlPath": "/site/t1",
			"ParamsArr": [],
			"ExpireTime": 10,
			"CheckCount": 50000
		}]
