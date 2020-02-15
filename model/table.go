package model

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/lexkong/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	MongoDb = "ccnubox"
	XkCol   = "table_xk"    // 教务课表
	SzkcCol = "table_szkc"  // 素质课程课表
	UserCol = "table_users" // 自定义课表
)

// 删除自定义课程
func DeleteTable(sid string, id string, objId primitive.ObjectID) (int64, error) {
	collection := DB.Self.Database(MongoDb).Collection(UserCol)

	// objId, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	log.Error("primitive.ObjectIDFromHex error", err)
	// 	return 0, err
	// }

	filter := bson.M{
		"sid": sid,
		"_id": objId,
	}

	delRes, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Error("collection.DeleteOne error", err)
		return 0, err
	}

	return delRes.DeletedCount, nil
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
func AddSelfTable(sid string, table *TableItem) (string, error) {
	collection := DB.Self.Database(MongoDb).Collection(UserCol)

	id := primitive.NewObjectID()
	table.Id = id.Hex()
	table.Source = "user"

	document := UserColModel{
		Id:    id,
		Sid:   sid,
		Table: table,
	}

	if _, err := collection.InsertOne(context.TODO(), document); err != nil {
		log.Error("InsertOne error", err)
		return "", err
	}

	return id.Hex(), nil
}

// 添加教务课表
func AddXKTable(sid string, tableList []*TableItem) error {
	collection := DB.Self.Database(MongoDb).Collection(XkCol)
	var err error

	// 有记录则为替换，无记录就插入
	if haveDoc, _ := HaveTable(sid); haveDoc {
		_, err = collection.ReplaceOne(context.TODO(), bson.M{"sid": sid}, TableModel{Sid: sid, Table: tableList})
	} else {
		_, err = collection.InsertOne(context.TODO(), TableModel{Sid: sid, Table: tableList})
	}

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
		log.Error("GetSelfTable error", err)
		return nil, err
	}

	if len(tableSelf) != 0 {
		tableList = append(tableList, tableSelf...)
	}

	// 素质课
	/*
		tableSzkc, err := GetSzkcTable(sid)
		if err != nil {
			return nil, err
		}

		if len(tableSzkc) != 0 {
			tableList = append(tableList, tableSzkc...)
		}
	*/

	return tableList, nil
}

// 获取素质课表
func GetSzkcTable(sid string) ([]*TableItem, error) {
	return nil, nil
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

	// 1-10周，1周，两种情况分别处理
	multiWeek := strings.Contains(weeksString, "-")
	if multiWeek {
		_, err := fmt.Sscanf(weeksString, "%d-%d", &weekStart, &weekEnd)
		if err != nil {
			return TableItem{}, err
		}
	} else {
		_, err := fmt.Sscanf(weeksString, "%d", &weekStart)
		if err != nil {
			return TableItem{}, err
		}
		weekEnd = weekStart
	}

	for i := weekStart; i <= weekEnd; i++ {
		if doubleWeek && i%2 != 0 || singleWeek && i%2 == 0 {
			continue
		}
		weeks = append(weeks, i)
	}

	var classStart, classEnd int

	// 节次：1-10，1，两种情况分别处理
	multiDuring := strings.Contains(rowTable.Jcor, "-")
	if multiDuring {
		_, err := fmt.Sscanf(rowTable.Jcor, "%d-%d", &classStart, &classEnd)
		if err != nil {
			return TableItem{}, err
		}
	} else {
		_, err := fmt.Sscanf(rowTable.Jcor, "%d", &classStart)
		if err != nil {
			return TableItem{}, err
		}
		classEnd = classStart
	}

	return TableItem{
		Course:  rowTable.Kcmc,
		Teacher: rowTable.Xm,
		Place:   rowTable.Cdmc,
		Start:   strconv.Itoa(classStart),
		During:  strconv.Itoa(classEnd - classStart + 1),
		Day:     rowTable.Xqj,
		Weeks:   weeks,
		Remind:  false,
	}, nil
}
