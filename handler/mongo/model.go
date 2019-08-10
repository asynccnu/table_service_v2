package mongo

type TableItem struct {
	Course		string `json:"course"`
	Teacher  	string `json:"teacher"`
	Place 		string `json:"place"`		// 上课地点
	//Start 		string `json:"start"`		// 课程开始时间(start=3表示上午第三节课开始上)
	//During 		string `json:"during"`		// 课程持续时间(during=2表示持续2节课)
	Day 		string `json:"day"`			// 上课星期,如 "1","2"..."7"
	Week 		string `json:"week"` 		// 上课周次
	During 		string `json:"during"`		// 上课节次
	//Weeks		[]int32 `json:"weeks"`		// 上课周次如[1, 2, 3 ... 19]
	//Remind 		bool   `json:"remind"`		// 是否提醒
}