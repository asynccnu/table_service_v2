package model

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	MongoDb		= "table"
	XkCol 		= "xk"		// 教务课表
	SzkcCol		= "szkc"	// 素质课程课表
	UserCol		= "users"	// 自定义课表
)

// 删除自定义课程
func DeleteTable(sid string, id int) error {
	collection := DB.Self.Database(MongoDb).Collection(UserCol)

	filter := bson.M{
		"sid": sid,
		"id": id,
	}

	if _, err := collection.DeleteOne(context.TODO(), filter); err != nil {
		return err
	}

	return nil
}

// 查看monggodb中是否有教务课表的记录
func HaveTable(sid string) (bool, error) {
	collection := DB.Self.Database(MongoDb).Collection(XkCol)

	count, err := collection.CountDocuments(context.TODO(), bson.M{"sid": sid})

	if err != nil {
		return false, err
	} else if count == 0 {
		return false, nil
	}

	return true, nil
}

// 添加自定义课程
func AddSelfTable(sid string, table *TableItem) (int64, error) {
	collection := DB.Self.Database(MongoDb).Collection(UserCol)

	// 以时间戳作为id
	// 一个缺陷，多个同时请求可能会产生相同的时间戳
	// 待解决
	id := time.Now().Unix()
	document := UserColModel{
		Id: id,
		Sid: sid,
		Table: table,
	}

	if _, err := collection.InsertOne(context.TODO(), document); err != nil {
		return 0, err
	}

	return id, nil
}

// 添加教务课表
func AddXKTable(sid string, tableList []*TableItem) error {
	collection := DB.Self.Database(MongoDb).Collection(XkCol)

	_, err := collection.InsertOne(context.TODO(), TableModel{Sid: sid, Table: tableList})

	if err != nil {
		return err
	}

	return nil
}

// 获取课表
func GetTable(sid string) ([]*TableItem, error) {
	// 获取教务课表
	tableList, err := GetXkTable(sid)
	if err != nil {
		return nil, err
	}

	// 获取自定义课表
	tableSelf, err := GetSelfTable(sid)
	if err != nil {
		return nil, err
	}

	if len(tableSelf) != 0 {
		tableList = append(tableList, tableSelf...)
	}

	// 素质课
	tableSzkc, err := GetSzkcTable(sid)
	if err != nil {
		return nil, err
	}

	if len(tableSzkc) != 0 {
		tableList = append(tableList, tableSzkc...)
	}

	return tableList, nil
}

// 获取素质课表
func GetSzkcTable(sid string) ([]*TableItem, error) {

}

// 从数据库中获取教务课表
func GetXkTable(sid string) ([]*TableItem, error) {
	var tableModel = new(TableModel)

	collection := DB.Self.Database(MongoDb).Collection(XkCol)

	err := collection.FindOne(context.TODO(), bson.M{"sid": sid}).Decode(tableModel)

	if err != nil {
		return nil, err
	}

	return tableModel.Table, nil
}

// 从数据库中获取自定义课表
func GetSelfTable(sid string) ([]*TableItem, error) {
	var result = make([]*TableItem, 0)

	collection := DB.Self.Database(MongoDb).Collection(UserCol)

	cur, err := collection.Find(context.TODO(), bson.M{"sid": sid})

	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var elem UserColModel
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		result = append(result, elem.Table)
	}

	return result, nil
}

// 加工得到的原始课表数据
func Process(rowTable *TableRowItem) (TableItem, error) {
	var weeks []int32
	var weekStart, weekEnd int32

	weeksString := rowTable.Zcd
	doubleWeek := strings.Contains(weeksString, "双")
	singleWeek := strings.Contains(weeksString, "单")

	_, err := fmt.Sscanf(weeksString, "%d-%d", &weekStart, &weekEnd)
	if err != nil {
		return TableItem{}, err
	}

	for i := weekStart; i <= weekEnd; i++ {
		if doubleWeek && i % 2 != 0  || singleWeek && i % 2 == 0{
			continue
		}
		weeks = append(weeks, i)
	}

	var classStart, classEnd int

	_, err = fmt.Sscanf(rowTable.Jcor, "%d-%d", &classStart, &classEnd)
	if err != nil {
		return TableItem{}, err
	}

	return TableItem{
		Course:		rowTable.Kcmc,
		Teacher:	rowTable.Xm,
		Place:		rowTable.Cdmc,
		Start:		strconv.Itoa(classStart),
		During:		strconv.Itoa(classEnd - classStart + 1),
		Day: 		rowTable.Xqj,
		Weeks:		weeks,
		Remind: 	false,
	}, nil
}
