###验证

本项目采用[Basic Auth](https://swagger.io/docs/specification/authentication/basic-authentication/).

在使用本项目任何API时都需要携带Authorization Header.

例如使用curl GET 方法访问API如下:

```bash
curl -X GET \
  http://127.0.0.1:8082/api/table/v2/ \
  -H 'Authorization: Basic MjAwdqqwsdqqqdDpdqwwqdqwTqdwqjEy'
```



### 数据格式

本项目返回数据(无论正确与否)采用如下数据格式

```
{
    "code": 0,
    "message": "OK",
    "data": {}
}
```



## APIS

### 1. 获取课程表

| Method | URL            | Header        |
| ------ | -------------- | ------------- |
| GET    | /api/table/v2/ | Authorization |



**Response**

```json
{
	"code":0,
	"message": "OK",
	"data": {		
		"table": [   
			{
				"id": "0",			// 课程id，可能形如"1"这种普通数字也可能为如 "5e36c9f5d3c47ef0e7bbe160"
				"course": "Java语言程序设计",		
				"teacher": "张勇",
				"place": "N113",		// 上课地点
				"start": "3",			// 课程开始时间(start=3表示上午第三节课开始上)
				"during": "2",			// 课程持续时间(during=2表示持续2节课)
				"day": "1",       		// 上课星期,如 "1","2"..."7"
				"source": "xk",   		// 课程来源：xk(教务处), szkc(素质课), user(自定义)
				"weeks": [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17],   // 哪些周上课
				"remind": false,
				"color": 0
			},
			{
				"id": "5e36c9f5d3c47ef0e7bbe160",
				"course": "用户添加课程",
				"teacher": "覃凤珍",
				"place": "网球场",
				"start": "5",
				"during": "2",
				"day": "2",
				"source": "xk",
				"weeks": [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17],
				"remind": false,
				"color": 1  // 颜色取0,1,2,3
			}
		]
	}
}
```





### 2. 添加课程

| Method | URL            | Header        |
| ------ | -------------- | ------------- |
| POST   | /api/table/v2/ | Authorization |



**POST DATA**

```json
{
    "course": "大学体育4",   
    "teacher": "覃凤珍",		 
    "place": "网球场",			  
    "start": "5",
    "during": "2",
    "day": "1",
    "weeks": [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17]
}
```



**RESPONSE DATA**

```json
{
    "code": 0,
    "message": "OK",
    "data": {
        "Id": "5e36d4b8baff5c52ee5ee631"   // 课程id
    }
}
```



### 3.删除课程

| Method | URL                        | Header        |
| ------ | -------------------------- | ------------- |
| DELETE | /api/table/v2/?id={课程id} | Authorization |



**RESPONSE DATA**

```json
{
    "code": 0,
    "message": "OK",
    "data": null
}
```



### 可能返回的HTTP状态码

| Status Code | Description                                     |
| ----------- | ----------------------------------------------- |
| 200         | OK  操作成功                                    |
| 400         | Bad Request 客户端发送的数据错误                |
| 401         | Unauthorized 客户端没有携带Authorization Header |
| 500         | 服务端内部错误                                  |

