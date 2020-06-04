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

服务容错：超时，熔断，隔离，限流，降级

gin 集成警告邮件
````

````
时间轮: 分层时间轮
添加任务
    1.定时执行任务   
    2.循环执行任务
删除任务
    1.立即删除
    2.下次任务结束后删除

````
````
//1.restful api
//2.限流
//3.接口防刷
//4.权限验证
//5.结果返回
//6.反爬虫
//7.熔断 //go hystrix
//8.ip黑/白名单
//9.灰度发布
//10.读写库
//11.敏感词
````

````
目前只关注功能测试,性能测试后期测试
````

````
本地缓存系统:
使用：github.com/allegro/bigcache

````

````
common:程序用到的所有中间件或其他公用程序
config:配置文件
core:核心服务
    chat:聊天服务
    connect:长链接维护服务
    logmanager:日志收集服务
    queue:redis 队列服务
    user:用户服务
gateway:网关
utils:工具代码
````