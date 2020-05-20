package models

import (
	"fmt"
	_ "fmt"
	"github.com/gao111/canal-adapter-go/libs"
	"github.com/gao111/canal-adapter-go/sysinit"
	_ "github.com/jinzhu/gorm"
	protocol "github.com/gao111/canal-adapter-go/protocol"
	"strings"
)

type Sync struct {
	DbName     string
	TableName  string
}

//func (User) TableName() string {
//  return "user"
//}

func NewSync(dbName string, tableName string) *Sync {
	return &Sync{
		DbName:    dbName,
		TableName: tableName,
	}
}
//func (u *User) GetUserById() {
//  u.ID = 0
//  //sysinit.Db.Where("user_id = ?", u.ID).Select("user_id, name , nickname , mobile").First(&u).Error
//  sysinit.Db.Table("user").Where("user_id = ?", 11202).Select("user_id, name , nickname , mobile").First(&u)
//  fmt.Print("haha \n")
//  fmt.Println(u)
//}

func (s *Sync)InsertSync(columns []*protocol.Column) {
	if len(columns) == 0 {
		return
	}

	setValue := ""
	for _, col := range columns {
		colName := col.GetName()
		colValue := col.GetValue()

		// 开头和结尾都是#号的, 为阿里云数据库特殊字段
		if strings.HasPrefix(colName, "#") && strings.HasSuffix(colName,"#") {
			continue
		}
		// 其中需要mysql加引号的有
		if setValue != "" {
			setValue += " , "
		}
		if isMysqlString(col.GetSqlType()){
			setValue += fmt.Sprintf(" %s = '%s'" , colName , colValue)
		} else {
			setValue += fmt.Sprintf(" %s = %s" , colName , colValue)
		}
	}
	if setValue == "" {
		return
	}
	//sysinit.Db.Exec()
	sql := fmt.Sprintf("insert into %s.%s set %s;" , s.DbName , s.TableName , setValue)
	fmt.Println(sql)
	sysinit.Db.Exec(sql)
}


func (s *Sync)DeleteSync(columns []*protocol.Column) {
	if len(columns) == 0 {
		return
	}

	where := ""
	for _, col := range columns {
		colName := col.GetName()
		colValue := col.GetValue()
		//fmt.Println(fmt.Sprintf("%s : %s  update= %t", col.GetName(), col.GetValue(), col.GetUpdated()))

		// 开头和结尾都是#号的, 为阿里云数据库特殊字段
		if strings.HasPrefix(colName, "#") && strings.HasSuffix(colName,"#") {
			continue
		}
		fmt.Println(fmt.Sprintf("colName : %s ,类型为: %s" , colName , libs.Typeof(colValue)))
		// col.GetSqlType返回的是java的sqlType,对应值如下 https://www.cnblogs.com/vickylinj/p/9787250.html
		// 其中需要mysql加引号的有
		if isMysqlString(col.GetSqlType()){
			where += fmt.Sprintf(" and %s = '%s'" , colName , colValue)
		} else {
			where += fmt.Sprintf(" and %s = %s" , colName , colValue)
		}
	}
	if where == "" {
		return
	}
	//sysinit.Db.Exec()
	sql := fmt.Sprintf("delete from %s.%s where 1=1 %s;" , s.DbName , s.TableName , where)
	// 增加limit1后,一次只会删除一条数据,保证数据安全,但无法保证数据真实一致
	//sql := fmt.Sprintf("delete from %s.%s where 1=1 %s limit 1")
	fmt.Println(sql)
	sysinit.Db.Exec(sql)
}


func (s *Sync)UpdateSync(columns []*protocol.Column) {
	if len(columns) == 0 {
		return
	}

	where := ""
	setValue := ""
	for _, col := range columns {
		colName := col.GetName()
		colValue := col.GetValue()
		isUpdated := col.GetUpdated()
		//fmt.Println(fmt.Sprintf("%s : %s  update= %t", col.GetName(), col.GetValue(), col.GetUpdated()))

		// 开头和结尾都是#号的, 为阿里云数据库特殊字段
		if strings.HasPrefix(colName, "#") && strings.HasSuffix(colName,"#") {
			continue
		}

		// 被修改的字段
		if isUpdated == true {
			if setValue != "" {
				setValue += " , "
			}
			if isMysqlString(col.GetSqlType()){
				setValue += fmt.Sprintf(" %s = '%s'" , colName , colValue)
			} else {
				setValue += fmt.Sprintf(" %s = %s" , colName , colValue)
			}
		} else {
			// 其中需要mysql加引号的有
			if isMysqlString(col.GetSqlType()){
				where += fmt.Sprintf(" and %s = '%s'" , colName , colValue)
			} else {
				where += fmt.Sprintf(" and %s = %s" , colName , colValue)
			}
		}
	}
	if where == "" || setValue == "" {
		return
	}
	//sysinit.Db.Exec()
	sql := fmt.Sprintf("update %s.%s set %s where 1=1 %s;" , s.DbName , s.TableName , setValue , where)
	fmt.Println(sql)
	sysinit.Db.Exec(sql)
}

/**
  判断获取的java type是不是mysql中的字符类型, 字符类型需要加上引号
*/
func isMysqlString(javaType int32) bool {
	// col.GetSqlType返回的是java的sqlType,对应值如下 https://www.cnblogs.com/vickylinj/p/9787250.html

	sqlTypeString := []int32{1,12,-1,2004,2005,-15,-9,-16,2011}
	exists , _ := libs.InArray(javaType , sqlTypeString)
	return exists
}