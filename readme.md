# Goal-web/collection

## 安装 - install
```bash
go get github.com/goal-web/collection
```

## 使用
```go
package tests

import (
	"fmt"
	"github.com/goal-web/collection"
	"github.com/goal-web/contracts"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	collect, err := collection.New(1)
	assert.Nil(t, collect)
	assert.Error(t, err, err)

	// 使用 MustNew 的时候，如果参数不是 array 或者 slice 的话，将会 panic
	collect, err = collection.New([]int{1})
	assert.NotNil(t, collect)
	assert.Nil(t, err)
}

func TestArray(t *testing.T) {
	intCollection := collection.MustNew([]interface{}{
		1, 2, 3, true, "字符串", "true",
	})
	fmt.Println(intCollection.ToFloat64Array())
	assert.True(t, intCollection.Length() == 6)

	// 第二个参数是数据索引
	intCollection.Map(func(data, index int) {
		fmt.Println(fmt.Sprintf("第 %d 个，值：%d", index, data))
	})

	// 第三个参数是所有数据集合
	intCollection.Map(func(data, index int, allData []interface{}) {
		if index == 0 {
			fmt.Println("allData", allData)
		}
		fmt.Println(fmt.Sprintf("第 %d 个，值：%d", index, data))
	})

	// 甚至可以直接转换成你想要的类型
	intCollection.Map(func(data string, index int) {
		fmt.Println(fmt.Sprintf("第 %d 个，值：%s", index, data))
	})
	intCollection.Map(func(data bool, index int) {
		fmt.Println(fmt.Sprintf("第 %d 个，值：%v", index, data))
	})

	// 不返回任何值表示只遍历
	intCollection.Map(func(data int) {
		fmt.Println("只遍历: ", data)
	})
	fmt.Println(intCollection.ToIntArray())

	// 返回一个值会生成一个新的 collection
	fmt.Println(intCollection.Map(func(data int) int {
		if data > 0 {
			return 1
		}
		return 0
	}).ToIntArray())
}

type User struct {
	id    int
	Name  string
	Money int
}

func TestStructArray(t *testing.T) {

	users := collection.MustNew([]User{
		{id: 1, Name: "qbhy"},
		{id: 2, Name: "goal"},
	})

	assert.Nil(t, users.Index(5))

	users.Map(func(user User) {
		fmt.Printf("user: id:%d Name:%s \n", user.id, user.Name)
	})
	// 使用 fields 接收的时候，未导出字段默认是 nil
	users.Map(func(user contracts.Fields) {
		fmt.Printf("user: id:%v Name:%s \n", user["id"], user["name"])
	})

	// 使用 map 修改数据后在用 where 筛选
	assert.True(t, users.Map(func(user User) User {
		if user.id == 1 {
			user.Money = 100
		}
		return user
	}).Where("money", 100).Length() == 1)
}

func TestFilterArray(t *testing.T) {

	users := collection.MustNew([]User{
		{id: 1, Name: "qbhy", Money: 10000000},
		{id: 2, Name: "goal", Money: 10},
	})

	fmt.Println("第一个数据", users.Index(0))

	richUsers := users.Filter(func(user User) bool {
		return user.Money > 100
	})

	assert.True(t, richUsers.Length() == 1)
	fmt.Println(richUsers.ToInterfaceArray())

	poorUsers := users.Skip(func(user User) bool {
		return user.Money > 100
	})

	assert.True(t, poorUsers.Length() == 1)
	fmt.Println(poorUsers.ToInterfaceArray())

	qbhyUsers := users.Where("name", "qbhy")

	assert.True(t, qbhyUsers.Length() == 1)
	fmt.Println(qbhyUsers.ToInterfaceArray())

	assert.True(t, users.WhereLte("money", 50).Length() == 1)
	assert.True(t, users.Where("money", "<=", 50).Length() == 1)

}
```

[goal-web/collection](https://github.com/goal-web/collection)  
qbhy0715@qq.com
