package table

type Reqeust struct {
	XNM string `form:"xnm"`
	XQM string `form:"xqm"`
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
