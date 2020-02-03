# table service

## releases

* ##### release-v1.0 @ 2020.02.01  完善基本功能，部署上线



## 运维Tips

1. 学期切换需要手动清除原来的课表数据库,以避免新学期获取不到课表时返回上学期课表

2. 可设置的环境变量如下

```
CCNUBOX_DATA_SERVICE_URL="url.to.data.service:12345"
CCNUBOX_DB_URL="mongodb://username:password@127.0.0.1:27017/?authSource=admin"
CCNUBOX_TABLE_XN=2018		// 当自动获取的学年不准确时使用此直接设定学年 
CCNUBOX_TABLE_XQ=16		// 当自动获取的学期不准确时使用此直接设定学期  3:第一学期 12:第二学期 16:第三学期
```
