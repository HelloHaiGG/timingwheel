# HelloMyWorld

#### es + tablestore 实现系统日志管理

####gateway 网关
* app 网关
    
* wms 网关

#### integrate 微服务整合，结果返回给网关

````
LRU算法
redis 分布式锁
gorm
雪花算法
etcd
验证码
调用链工具集成 jaeger
ant携程池
grpc通信
时间轮

gin 集成警告邮件
````

````
目前只关注功能测试,性能测试后期测试
````

````
缓存系统:
L1:go-cache 热点数据 基本不会改变的数据
L2:redis 相对热点数据

策略:
    系统提供三个方法
    get(key string) string,error
    put(key string) error
    reset(key string) error

    系统维护一个sync.Map[key][flag] //考虑到分布式,可以用redis代替
    flag:表示 key 对应的 value 值是否在进行 reset 操作
         flag = 1 
            get操作:  
            sleep 0.01s 判断,并判断三次 (0.01s,0.02,0.04)
            三次都失败后返回 error
            reset操作: sleep 0.01s 再次判断,判断两次 (0.01s,0.03s)
            两次都失败后返回 error 这时候可能存在异常情况
            put操作:
            直接返回 error
         flag = 0
            get操作:
            直接在L1,L2中获取数据
            reset操作:
            将 key 对应的flag设置为1,并更改L1,L2 
            put 将flag设置为1,并写入L1,L2
````
 