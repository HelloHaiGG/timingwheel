package safety

//签名策略
//前端参数传递必传参数有 TS,Sign,Version
//通过TS对2取余,判断参数排序方式 0：正序 1：倒序
//将参数组合成字符串,生成sign
//判断前端穿过来的sign与生成的sign是否一样