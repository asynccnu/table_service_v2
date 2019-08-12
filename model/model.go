package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type TableRowItem struct {
	Kcmc	string `json:"kcmc" bson:"kcmc"`		// 课程名称
	Zcd		string `json:"zcd" bson:"zcd"`			// 周次
	Jcor	string `json:"jcor" bson:"jcor"`		// 节次
	Cdmc	string `json:"cdmc"bson:"cdmc"`			// 上课地点
	Xm		string `json:"xm" bson:"xm"`			// 教师名
	Xqj		string `json:"xqj" bson:"xqj"`			// 星期几
	KchID	string `json:"kch_id" bson:"kch_id"`	// 课程号ID
	JxbID	string `json:"jxb_id" bson:"jxb_id"`	// 教学班ID
}

type TableItem struct {
	Course		string `json:"course" bson:"course"`
	Teacher  	string `json:"teacher" bson:"teacher"`
	Place 		string `json:"place" bson:"place"`		// 上课地点
	Start 		string `json:"start" bson:"start"`		// 课程开始时间(start=3表示上午第三节课开始上)
	During 		string `json:"during" bson:"during"`	// 课程持续时间(during=2表示持续2节课)
	Day 		string `json:"day" bson:"day"`			// 上课星期,如 "1","2"..."7"
	Weeks		[]int32 `json:"weeks" bson:"weeks"`		// 上课周次如[1, 2, 3 ... 19]
	Remind 		bool   `json:"remind" bson:"remind"`	// 是否提醒
}

// 教务课表和素质课的mongo存储结构
type TableModel struct {
	Sid 	string			`bson:"sid"`
	Table 	[]*TableItem	`json:"table "bson:"table"`
}

// 自定义课程的mongo存储结构
type UserColModel struct {
	Id 		primitive.ObjectID	`bson:"_id"`
	Sid 	string			`bson:"sid"`
	Table 	*TableItem		`bson:"table"`
}