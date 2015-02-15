package zico_test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ugorji/go/codec"
	"hash/crc32"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	//	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Query struct {
	_type string
	unit  string
	key   string
	from  time.Time
	to    time.Time
}

type MarcoInfo map[int]int64

var MarcoInfoMap map[int]string = map[int]string{
	1:  "Rcnt",
	2:  "Rtime",
	3:  "Rsize",
	4:  "Cc",
	5:  "Inj",
	6:  "Xss",
	7:  "Lg_rcnt",
	8:  "Lg_rtime",
	9:  "Lg_rsize",
	10: "Eghits",
	11: "Hits",
	12: "S200",
	13: "S206",
	14: "S301",
	15: "S302",
	16: "S304",
	17: "S400",
	18: "S403",
	19: "S404",
	20: "S411",
	21: "S499",
	22: "S500",
	23: "S502",
	24: "S503",
	25: "S504",
}

type MarcoDetailInfo map[string]map[string]MarcoInfo

type MarcoRecord struct {
	Id    int64           `codec:"id"`
	Now   int64           `codec:"now"`
	Key   string          `codec:"key"`
	Stats MarcoInfo       `codec:"stats"`
	Nwpv  MarcoDetailInfo `codec:"nwpv"`
	Nwct  MarcoDetailInfo `codec:"nwct"`
}

type MarcoRecords []MarcoRecord

type ZicogoMsg struct {
	Checksum uint32   `codec:"checksum"`
	Records  [][]byte `codec:"records"`
}

var ct = []string{
	//	"北京", "上海", "深圳", "杭州",
	"北京",
}

var nw = []string{
	//	"移动", "联通", "电信", "unknown",
	"移动",
}

var pv = []string{
	"北京", //"上海", "浙江", "江苏", "辽宁", "台湾",
}

func GenMarcoInfo(args ...string) (msg MarcoInfo) {
	msg = MarcoInfo{}
	msg[1] = rand.Int63n((int64)(12345678910)) + 1
	for i := 2; i <= 2; i++ {
		if rand.Intn(2) == 1 {
			msg[i] = rand.Int63n((int64)(12345678910))
		}
	}
	return
}

func MergeMarcoInfo(msg1 MarcoInfo, msg2 MarcoInfo) (msg MarcoInfo) {
	msg = MarcoInfo{}
	for i := 1; i <= 2; i++ {
		msg[i] = msg1[i] + msg2[i]
	}
	return
}

func DisplayMarcoInfo(msg MarcoInfo, offset int) {
	leading := strings.Repeat(" ", 4*offset)
	var keys []int

	for k, _ := range msg {
		if msg[k] != 0 {
			keys = append(keys, k)
		}
	}
	sort.Ints(keys)

	for _, k := range keys {
		fmt.Printf("%s%s = %d\n", leading, MarcoInfoMap[k], msg[k])
	}
}

func GenMarcoDetailInfo(name []string, args ...string) (msg MarcoDetailInfo) {
	msg = MarcoDetailInfo{}
	for _, _nw := range nw {
		if rand.Intn(2) == 1 {
			msg[_nw] = map[string]MarcoInfo{}
			for _, w := range name {
				if rand.Intn(2) == 1 {
					msg[_nw][w] = GenMarcoInfo()
				}
			}
		}
	}
	return
}

func nwmapMerge(map1, map2 map[string]MarcoInfo) (res map[string]MarcoInfo) {
	res = map[string]MarcoInfo{}
	for k := range map1 {
		res[k] = map1[k]
	}
	for k := range map2 {
		if _, ok := map2[k]; ok {
			res[k] = MergeMarcoInfo(res[k], map2[k])
		} else {
			res[k] = map2[k]
		}
	}
	return
}

func MergeMarcoDetailInfo(msg1, msg2 MarcoDetailInfo) (msg MarcoDetailInfo) {
	msg = MarcoDetailInfo{}
	for k := range msg1 {
		msg[k] = msg1[k]
	}
	for k := range msg2 {
		if _, ok := msg2[k]; ok {
			msg[k] = nwmapMerge(msg[k], msg2[k])
		} else {
			msg[k] = msg2[k]
		}
	}
	return
}

func DisplayMarcoDetailInfo(msg MarcoDetailInfo, offset int) {
	leading := strings.Repeat(" ", 4*offset)
	for _nw, v := range msg {
		fmt.Printf("%s%s:{\n", leading, _nw)
		for w, mi := range v {
			offset++
			leading := strings.Repeat(" ", 4*offset)
			fmt.Printf("%s%s:{\n", leading, w)
			DisplayMarcoInfo(mi, offset+1)
			fmt.Printf("%s},\n", leading)
			offset--
		}
		leading := strings.Repeat(" ", 4*offset)
		fmt.Printf("%s},\n", leading)
	}
}

func GenMarcoRecord(now int64, domain string, args ...string) (record MarcoRecord) {
	record = MarcoRecord{}
	record.Now = now
	record.Key = domain
	record.Stats = GenMarcoInfo()
	record.Nwpv = GenMarcoDetailInfo(pv)
	record.Nwct = GenMarcoDetailInfo(ct)
	return
}

func MergeMarcoRecord(msg1, msg2 MarcoRecord, now int64) (record MarcoRecord) {
	record = MarcoRecord{}
	record.Now = now
	record.Key = msg1.Key
	record.Stats = MergeMarcoInfo(msg1.Stats, msg2.Stats)
	record.Nwpv = MergeMarcoDetailInfo(msg1.Nwpv, msg2.Nwpv)
	record.Nwct = MergeMarcoDetailInfo(msg1.Nwct, msg2.Nwct)
	return
}

func DisplayMarcoRecord(record MarcoRecord, offset int) {
	leading := strings.Repeat(" ", 4*offset)
	fmt.Printf("%snow: %d,\n", leading, record.Now)
	fmt.Printf("%skey: %s,\n", leading, record.Key)

	fmt.Printf("%sstats: {\n", leading)
	DisplayMarcoInfo(record.Stats, offset+1)
	fmt.Printf("%s},\n", leading)

	fmt.Printf("%snwpv: {\n", leading)
	DisplayMarcoDetailInfo(record.Nwpv, offset+1)
	fmt.Printf("%s},\n", leading)

	fmt.Printf("%snwct: {\n", leading)
	DisplayMarcoDetailInfo(record.Nwct, offset+1)
	fmt.Printf("%s},\n", leading)
}

func GenMarcoRecords(size int, now int64, domain string, args ...string) (records []MarcoRecord) {
	records = make([]MarcoRecord, size)
	for i := 0; i < size; i++ {
		records[i] = GenMarcoRecord(now, domain)
	}
	return
}

func EncMarcoRecords(records []MarcoRecord, secret string) ([]byte, error) {
	// records[] to byte[]
	msg := ZicogoMsg{
		Records: [][]byte{},
	}
	var record_b = new([]byte)
	for _, r := range records {
		enc := codec.NewEncoderBytes(record_b, &codec.MsgpackHandle{})
		err := enc.Encode(&r)
		if err != nil {
			return nil, err
		}
		msg.Records = append(msg.Records, *record_b)
	}

	// encode records
	buffer := new(bytes.Buffer)
	buffer.WriteString(secret)
	buffer.Write(msg.Records[0])
	msg.Checksum = crc32.ChecksumIEEE(buffer.Bytes())

	var b = new([]byte)
	enc := codec.NewEncoderBytes(b, &codec.MsgpackHandle{})
	err := enc.Encode(msg)
	return *b, err

}

func ZicogoSend(zicod_addr string, records []byte) error {
	hostname, _ := os.Hostname()
	request, _ := http.NewRequest("POST", zicod_addr, bytes.NewBuffer(records))
	request.Header.Set("Content-Type", "application/octet-stream")
	request.Header.Set("X-Hostname", hostname)

	tr := &http.Transport{ResponseHeaderTimeout: 30 * time.Second}
	client := &http.Client{Transport: tr, Timeout: 30 * time.Second}

	resp, err := client.Do(request)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp != nil && resp.StatusCode != 200 {
		return errors.New("httpcode" + strconv.Itoa(resp.StatusCode))
	}

	return nil
}

func ZicodQuery(zicod_addr string, query Query) (interface{}, error) {
	var u int64
	switch query.unit {
	case "min":
		u = 60
	case "hour":
		u = 60 * 60
	case "day":
		u = 24 * 60 * 60
	}

	form := url.Values{
		"type": {query._type},
		"key":  {query.key},
		"unit": {query.unit},
		"from": {fmt.Sprint(query.from.Unix() / u)},
		"to":   {fmt.Sprint(query.to.Unix() / u)},
	}

	resp, err := http.PostForm(zicod_addr+"/query", form)
	if err != nil {
		// handle error
		fmt.Printf("failed query data, %s\n", err.Error)
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Printf("failed query data, %d\n", resp.StatusCode)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("response body: %s\n", body)
	return body, nil
}

//func main() {
//	rand.Seed(time.Now().UTC().UnixNano())
//	//	from := time.Unix(1354773785, 0)
//	//	to := time.Unix(1354773885, 0)
//	//	query := Query{"dm", "min", "test.zico2.com", from, to}
//	//	ZicodQuery("http://10.0.1.233:9095", query)
//	//	GenMarcoInfo()
//	record1 := GenMarcoDetailInfo(ct)
//	record2 := GenMarcoDetailInfo(ct)
//	record3 := MergeMarcoDetailInfo(record1, record2)
//	DisplayMarcoDetailInfo(record1, 1)
//	DisplayMarcoDetailInfo(record2, 2)
//	DisplayMarcoDetailInfo(record3, 3)
//	fmt.Println(reflect.ValueOf(record1).Type())
//	fmt.Println(reflect.ValueOf(record2).Type())
//	fmt.Println(reflect.ValueOf(record3).Type())
//	//  GenMarcoRecord(12345, "test.zico.com")
//	now := time.Now().Unix()
//	fmt.Println(now / 60)
//	records := GenMarcoRecords(2, now, "test.zico.com")
//	for i := 0; i < 2; i++ {
//		DisplayMarcoRecord(records[i], 2)
//	}
//	record4 := MergeMarcoRecord(records[0], records[1], 0)
//	DisplayMarcoRecord(records[0], 3)
//	DisplayMarcoRecord(records[1], 3)
//	DisplayMarcoRecord(record4, 3)
//
//	record_b, _ := EncMarcoRecords(records, "jaychou")
//	_ = ZicogoSend("http://10.0.1.233:9095", record_b)
//
//	time.Sleep(10 * time.Second)
//
//	from := time.Unix(now, 0)
//	to := time.Unix(now+60, 0)
//	query := Query{"dm", "min", "test.zico.com", from, to}
//	ZicodQuery("http://10.0.1.233:9095", query)
//
//	//	for i := 0; i < 1; i++ {
//	//		fmt.Println(*b)
//	//	}
//	//	s := reflect.TypeOf((*MarcoInfo)(nil)).Elem()
//	//	typeOfT := s.Type()
//	//	for i := 0; i < s.NumField(); i++ {
//	//		fmt.Println(typeOfT.Field(i).Name)
//	//	}
//}
