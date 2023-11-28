package wgorm

import (
	"fmt"
	"github.com/devfeel/mapper"
	"github.com/jinzhu/copier"
	"testing"
	"time"
)

type (
	User struct {
		Name string
		Age  int
		Id   string `mapper:"_id"`
		AA   string `json:"Score"`
		Time time.Time
	}

	Student struct {
		Name  string
		Age   int
		Id    string `mapper:"_id"`
		Score string
	}

	Teacher struct {
		Name  string
		Age   int
		Id    string `mapper:"_id"`
		Level string
	}
)

func BenchmarkCopy(b *testing.B) {
	b.Run("copier", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			user := &User{}
			userMap := &User{}
			teacher := &Teacher{}
			student := &Student{Name: "test", Age: 10, Id: "testId", Score: "100"}
			valMap := make(map[string]interface{})
			valMap["Name"] = "map"
			valMap["Age"] = 10
			valMap["_id"] = "x1asd"
			valMap["Score"] = 100
			valMap["Time"] = time.Now()

			//mapper.Mapper(student, user)
			//mapper.AutoMapper(student, teacher)
			//mapper.MapperMap(valMap, userMap)

			copier.Copy(user, student)
			copier.Copy(teacher, student)
			copier.Copy(userMap, valMap)

			//fmt.Println("student:", student)
			//fmt.Println("user:", user)
			//fmt.Println("teacher", teacher)
			//fmt.Println("userMap:", userMap)
		}
	})

	b.Run("mapper", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			user := &User{}
			userMap := &User{}
			teacher := &Teacher{}
			student := &Student{Name: "test", Age: 10, Id: "testId", Score: "100"}
			valMap := make(map[string]interface{})
			valMap["Name"] = "map"
			valMap["Age"] = 10
			valMap["_id"] = "x1asd"
			valMap["Score"] = 100
			valMap["Time"] = time.Now()

			mapper.Mapper(student, user)
			mapper.AutoMapper(student, teacher)
			mapper.MapperMap(valMap, userMap)

			//copier.Copy(user, student)
			//copier.Copy(teacher, student)
			//copier.Copy(userMap, valMap)

			//fmt.Println("student:", student)
			//fmt.Println("user:", user)
			//fmt.Println("teacher", teacher)
			//fmt.Println("userMap:", userMap)
		}
	})

}

func TestXxx(t *testing.T) {
	ss1 := map[string]any{
		"A": 1,
	}

	type s2 struct {
		A int
	}

	var ss2 s2

	//copier.Copy(&ss2, ss1)
	mapper.MapperMap(ss1, &ss2)

	fmt.Printf("%+v\n", ss2)
}
