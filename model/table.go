package model

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	MongoDb		= "table"
	XkCol 		= "xk"		// 教务课表
	SzkcCol		= "szkc"	// 素质课程课表
	UserCol		= "users"	// 自定义课表
)

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
func AddSelfTable(sid string, table *TableItem) error {
	collection := DB.Self.Database(MongoDb).Collection(UserCol)

	tableList, err := GetSelfTable(sid)
	if err != nil {
		return err
	}

	tableList = append(tableList, table)

	if _, err := collection.ReplaceOne(context.TODO(), bson.M{"sid": sid}, tableList); err != nil {
		return err
	}

	return nil
}

// 添加教务课表
func AddXKTable(sid string, tableList []*TableItem) error {
	collection := DB.Self.Database(MongoDb).Collection(XkCol)

	_, err := collection.InsertOne(context.TODO(), MgoTable{Sid: sid, Table: tableList})

	if err != nil {
		return err
	}

	return nil
}

// 获取课表
func GetTable(sid string) ([]*TableItem, error) {
	/*
	collection := DB.Self.Database(MongoDb).Collection(MgoCollection)

	var result = make([]*TableItem, 0)

	cur, err := collection.Find(context.TODO(), bson.D{})

	if err != nil {
		return result, err
	}

	for cur.Next(context.TODO()) {
		var elem TableItem
		err := cur.Decode(&elem)
		if err != nil {
			return result, err
		}
		result = append(result, &elem)
	}

	return result, nil
	 */

	// 获取教务课表
	tableList, err := GetTableFromXK(sid)
	if err != nil {
		return nil, err
	}

	// 获取自定义课表
	tableOwn, err := GetSelfTable(sid)
	if err != nil {
		return nil, err
	}

	tableList = append(tableList, tableOwn...)

	return tableList, nil
}

// 从数据库中获取教务课表
func GetTableFromXK(sid string) ([]*TableItem, error) {
	var tableModel = new(MgoTable)

	collection := DB.Self.Database(MongoDb).Collection(XkCol)

	err := collection.FindOne(context.TODO(), bson.M{"sid": sid}).Decode(tableModel)

	if err != nil {
		return nil, err
	}

	return tableModel.Table, nil
}

// 从数据库中获取自定义课表
func GetSelfTable(sid string) ([]*TableItem, error) {
	var tableModel = new(MgoTable)

	collection := DB.Self.Database(MongoDb).Collection(UserCol)

	err := collection.FindOne(context.TODO(), bson.M{"sid": sid}).Decode(tableModel)

	if err != nil {
		return nil, err
	}

	return tableModel.Table, nil
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
