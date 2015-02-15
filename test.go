package zico_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestA(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	//	from := time.Unix(1354773785, 0)
	//	to := time.Unix(1354773885, 0)
	//	query := Query{"dm", "min", "test.zico2.com", from, to}
	//	ZicodQuery("http://10.0.1.233:9095", query)
	//	GenMarcoInfo()
	record1 := GenMarcoDetailInfo(ct)
	record2 := GenMarcoDetailInfo(ct)
	record3 := MergeMarcoDetailInfo(record1, record2)
	DisplayMarcoDetailInfo(record1, 1)
	DisplayMarcoDetailInfo(record2, 2)
	DisplayMarcoDetailInfo(record3, 3)
	fmt.Println(reflect.ValueOf(record1).Type())
	fmt.Println(reflect.ValueOf(record2).Type())
	fmt.Println(reflect.ValueOf(record3).Type())
	//  GenMarcoRecord(12345, "test.zico.com")
	now := time.Now().Unix()
	fmt.Println(now / 60)
	records := GenMarcoRecords(2, now, "test.zico.com")
	for i := 0; i < 2; i++ {
		DisplayMarcoRecord(records[i], 2)
	}
	record4 := MergeMarcoRecord(records[0], records[1], 0)
	DisplayMarcoRecord(records[0], 3)
	DisplayMarcoRecord(records[1], 3)
	DisplayMarcoRecord(record4, 3)

	record_b, _ := EncMarcoRecords(records, "jaychou")
	_ = ZicogoSend("http://10.0.1.233:9095", record_b)

	time.Sleep(10 * time.Second)

	from := time.Unix(now, 0)
	to := time.Unix(now+60, 0)
	query := Query{"dm", "min", "test.zico.com", from, to}
	ZicodQuery("http://10.0.1.233:9095", query)

	//	for i := 0; i < 1; i++ {
	//		fmt.Println(*b)
	//	}
	//	s := reflect.TypeOf((*MarcoInfo)(nil)).Elem()
	//	typeOfT := s.Type()
	//	for i := 0; i < s.NumField(); i++ {
	//		fmt.Println(typeOfT.Field(i).Name)
	//	}
}
