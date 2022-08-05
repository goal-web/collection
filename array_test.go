package collection

import (
	"errors"
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	collect := New([]int{1})
	assert.NotNil(t, collect)
}

func TestArray(t *testing.T) {
	anyCollection := New([]any{
		1, 2, 3, true, "字符串", "true",
	})
	fmt.Println(anyCollection.ToFloat64Array())
	assert.True(t, anyCollection.Len() == 6)

	// 第二个参数是数据索引
	anyCollection.Map(func(data any, index int) any {
		return data
	})

	// 返回一个值会生成一个新的 collection
	fmt.Println(anyCollection.Map(func(data any, index int) any {
		if i, ok := data.(int); ok && i > 0 {
			return 1
		}
		return 0
	}).ToIntArray())
}

type User struct {
	id    int
	Name  string `json:"name"`
	Money float64
}

func TestToJson(t *testing.T) {
	users := New([]User{
		{id: 1, Name: "qbhy", Money: 1},
		{id: 2, Name: "goal", Money: 2},
		{id: 3, Name: "collection", Money: 0},
	})
	firstUser, _ := users.First()
	arr := New([]any{
		"1", 2, 3.0, true, firstUser, users,
	})

	fmt.Println(users.ToJson())
	fmt.Println(arr.ToJson())
}

func TestStructArray(t *testing.T) {

	users := New([]User{
		{id: 1, Name: "qbhy"},
		{id: 2, Name: "goal"},
	})

	// 使用 fields 接收的时候，未导出字段默认是 nil
	users.Each(func(user User, int2 int) {
		fmt.Printf("user: id:%v Name:%s \n", user.id, user.Name)
	})

	// 使用 map 修改数据后在用 where 筛选
	assert.True(t, users.Map(func(user User, int2 int) User {
		if user.id == 1 {
			user.Money = 100
		}
		return user
	}).Filter(func(user User, int2 int) bool {
		return user.Money == 100
	}).Len() == 1)
}

func TestFilterArray(t *testing.T) {

	users := New([]User{
		{id: 1, Name: "qbhy", Money: 10000000},
		{id: 2, Name: "goal", Money: 10},
	})

	fmt.Println("第一个数据", users.ToAnyArray()[0])

	richUsers := users.Filter(func(user User, int2 int) bool {
		return user.Money > 100
	})

	assert.True(t, richUsers.Len() == 1)
	fmt.Println(richUsers.ToAnyArray())

	poorUsers := users.Skip(func(user User, int2 int) bool {
		return user.Money > 100
	})

	assert.True(t, poorUsers.Len() == 1)
	fmt.Println(poorUsers.ToAnyArray())

	qbhyUsers := users.Filter(func(user User, int2 int) bool {
		return user.Name == "qbhy"
	})

	assert.True(t, qbhyUsers.Len() == 1)
	fmt.Println(qbhyUsers.ToAnyArray())
	fmt.Println(users.ToAnyArray())

	assert.True(t, users.Filter(func(user User, int2 int) bool {
		fmt.Println("user.money", user.Money)
		return user.Money == 10
	}).Len() == 1)
	assert.True(t, users.Filter(func(user User, int2 int) bool {
		return user.Money <= 50
	}).Len() == 1)

}

// TestAggregateArray 聚合函数测试
func TestAggregateArray(t *testing.T) {
	//
	//numbers := New([]float64{
	//	10000000000000000,
	//	10000000000000000,
	//	0.645624123,
	//}).(*Collection[float64])
	//
	//sum, _ := decimal.NewFromString("20000000000000000.645624123")
	//avg, _ := decimal.NewFromString("6666666666666666.8818747076666667")
	//max, _ := decimal.NewFromString("10000000000000000")
	//min, _ := decimal.NewFromString("0.645624123")
	//
	//fmt.Println("numbers.SafeSum()", numbers.SafeSum())
	//assert.True(t, numbers.SafeSum().Equal(sum))
	//assert.True(t, numbers.SafeAvg().Equal(avg))
	//assert.True(t, numbers.SafeMax().Equal(max))
	//assert.True(t, numbers.SafeMin().Equal(min))
	//users := New([]User{
	//	{id: 1, Name: "qbhy", Money: 10000000000000000},
	//	{id: 2, Name: "goal", Money: 10000000000000000},
	//	{id: 3, Name: "collection", Money: 0.645624123},
	//}).(*Collection[User])
	//
	//// SafeSum、SafeAvg、SafeMax、SafeMin 等方法需要 *collection.Collection 类型
	//fmt.Println("Sum", users.SafeSum())
	//fmt.Println("Avg", users.SafeAvg())
	//fmt.Println("Max", users.SafeMax())
	//fmt.Println("Min", users.SafeMin())
	//sum, _ := decimal.NewFromString("20000000000000000.645624123")
	//avg, _ := decimal.NewFromString("6666666666666666.8818747076666667")
	//max, _ := decimal.NewFromString("10000000000000000")
	//min, _ := decimal.NewFromString("0.645624123")

	newUsers := New([]User{
		{id: 1, Name: "qbhy", Money: 1},
		{id: 2, Name: "goal", Money: 2},
		{id: 3, Name: "collection", Money: 0},
	})
	//
	//assert.True(t, newUsers.Sum() == 3)
	//assert.True(t, newUsers.Avg() == 1)
	//assert.True(t, newUsers.Max() == 2)
	//assert.True(t, newUsers.Min() == 0)
	assert.True(t, newUsers.Len() == 3)
}

// TestSortArray 测试排序功能
func TestSortArray(t *testing.T) {
	users := New([]User{
		{id: 1, Name: "qbhy", Money: 12},
		{id: 2, Name: "goal", Money: 1},
		{id: 2, Name: "goal", Money: 15},
		{id: 2, Name: "goal99", Money: 99},
		{id: 3, Name: "collection", Money: -5},
		{id: 3, Name: "移动", Money: 10086},
	})

	fmt.Println(users.ToAnyArray())

	// 暂不支持转成 contracts.Fields
	usersOrderByMoneyDesc := users.Sort(func(user User, next User) bool {
		return user.Money > next.Money
	})
	fmt.Println(usersOrderByMoneyDesc.ToAnyArray())
	assert.True(t, usersOrderByMoneyDesc.ToAnyArray()[0].(User).Money == 10086)

	usersOrderByMoneyAsc := users.Sort(func(user User, next User) bool {
		return user.Money < next.Money
	})
	fmt.Println(usersOrderByMoneyAsc.ToAnyArray())
	assert.True(t, usersOrderByMoneyAsc.ToAnyArray()[0].(User).Money == -5)

	numbers := New([]any{
		8, 0, 1, 2, 0.6, 4, 5, 6, -0.2, 7, 9, 3, "10086",
	})

	sortedNumbers := numbers.Sort(func(i, j any) bool {
		return cast.ToInt(i) > cast.ToInt(j)
	}).ToFloat64Array()

	fmt.Println(sortedNumbers)
	assert.True(t, sortedNumbers[0] == 10086)
}

// TestCombine 测试组合集合功能
func TestCombine(t *testing.T) {
	users := New([]User{
		{id: 1, Name: "qbhy", Money: 12},
	})

	users = users.Push(User{id: 2, Name: "goal", Money: 1000})
	//users = users.Prepend(User{id: 2, Name: "goal", Money: 1000}) // 插入到开头

	assert.True(t, users.Len() == 2)
	fmt.Println(users.ToAnyArray())

	others := New([]User{
		{id: 3, Name: "马云", Money: 100000000},
	})

	all := others.Merge(users).Sort(func(pre User, next User) bool {
		return pre.Money > next.Money
	})

	assert.True(t, all.Len() == 3)
	fmt.Println(all.ToAnyArray())

	firstUser, existsFirst := all.First()
	assert.True(t, firstUser.Name == "马云" && existsFirst) // 最有钱还是马云

	//normalUsers := all.Where("money", ">", 100)
	normalUsers := all.Filter(func(user User, index int) bool {
		return user.Money > 100
	})
	lastUser, existsLast := normalUsers.Last()
	assert.True(t, normalUsers.Len() == 2)                // 两个普通人
	assert.True(t, lastUser.Name == "goal" && existsLast) // 筛选不影响排序，跟马云比还差了点
	assert.False(t, normalUsers.IsEmpty())                // 有普通人
	assert.True(t, normalUsers.Filter(func(user User, index int) bool {
		return user.Money < 0
	}).IsEmpty()) // 普通人都没有负债

	randomUsers := all.Random(2)
	// 随机获取两个数据
	assert.True(t, randomUsers.Len() == 2)
	fmt.Println(randomUsers.ToAnyArray())

	pullUser, existsPull := all.Pull()
	assert.True(t, existsPull && pullUser.Name == "qbhy") // 从末尾取走一个
	assert.True(t, all.Len() == 2)                        // 判断取走后的长度
	shiftUser, existsUser := all.Shift()
	assert.True(t, existsUser && shiftUser.Name == "马云") // 从开头取走一个
	assert.True(t, all.Len() == 1)                       // 判断取走后的长度

}

func TestChunk(t *testing.T) {
	collect := New([]int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
	})

	err := collect.Chunk(5, func(collection contracts.Collection[int], page int) error {
		fmt.Printf("页码：%d，数量：%d %v\n", page, collection.Len(), collection.ToAnyArray())
		switch page {
		case 4:
			assert.True(t, collection.Len() == 4)
		default:
			assert.True(t, collection.Len() == 5)
		}
		return nil
	})
	assert.Nil(t, err)

	err = New([]User{
		{id: 1, Name: "qbhy", Money: 12},
		{id: 2, Name: "goal", Money: 1},
		{id: 2, Name: "goal", Money: 15},
		{id: 2, Name: "goal99", Money: 99},
		{id: 3, Name: "collection", Money: -5},
		{id: 3, Name: "移动", Money: 10086},
	}).Chunk(3, func(collection contracts.Collection[User], page int) error {
		assert.True(t, page == 1)
		firstUser, _ := collection.First()
		lastUser, _ := collection.Last()
		assert.True(t, firstUser.Name == "qbhy")
		assert.True(t, lastUser.Name == "goal")
		return errors.New("第一页退出")
	})

	assert.Error(t, err)

}
