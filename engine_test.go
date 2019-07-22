/*
 * Copyright (c) 2019 uplus.io
 */

package udb

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/hashicorp/memberlist"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"
	"unsafe"
	"uplus.io/udb/btree"
	"uplus.io/udb/hash"
	"uplus.io/udb/logger"
	"uplus.io/udb/mmap"
	"uplus.io/udb/store"
)

func TestEngine_Storage(t *testing.T) {
	//config := config.EngineConfig{StorePath: "data/engine-test", Namespace: "test-db", partitionSize: 3}
	//engine := NewEngine(config)
	//storage := engine.Storage
	//storage().Open()
	//defer storage().Close()
	//
	////test create
	//schema := store.NewSchema("users")
	//schema.Create("id", store.BASE_TYPE_STRING, 32)
	//schema.Create("username", store.BASE_TYPE_STRING, 16)
	//schema.Create("age", store.BASE_TYPE_INT, 0)
	//schema.Create("mobile", store.BASE_TYPE_STRING, 11)
	////3422 0119 8006 0124 53
	//schema.Create("certificate", store.BASE_TYPE_STRING, 18)
	//schema.Create("password", store.BASE_TYPE_STRING, 32)
	//updateSchema := storage().UpdateSchema(schema)
	//
	////test delete
	////schema.Remove("age")
	////storage().UpdateSchema(schema)
	//
	//storage().Put("users", "id", "1", "1", 1)
	//storage().Put("users", "username", "1", "sunding", 1)
	//storage().Put("users", "age", "1", 38, 1)
	//storage().Put("users", "mobile", "1", "13601800602", 1)
	//storage().Put("users", "certificate", "1", "3422 0119 8006 0124 53", 1)
	//storage().Put("users", "password", "1", "admin", 1)
	//
	//fmt.Println(updateSchema)
}

func TestEngine_Hash(t *testing.T) {
	for i := 0; i < 1000; i++ {
		hash := hash.UInt32Of(fmt.Sprintf("%d", i))
		fmt.Printf("hash %d %d\n", i, hash)
	}
}

var testPath = filepath.Join(".", "test-data")
var testData = []byte("0123456789ABCDEF")

func init() {
	f := openFile(os.O_RDWR | os.O_CREATE | os.O_TRUNC)
	f.Write(testData)
	f.Close()
}

func openFile(flags int) *os.File {
	f, err := os.OpenFile(testPath, flags, 0644)
	if err != nil {
		panic(err.Error())
	}
	return f
}

func TestMmap(t *testing.T) {
	f := openFile(os.O_RDWR)
	defer f.Close()
	mmap, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		t.Errorf("error mapping: %s", err)
	}

	defer mmap.Unmap()
	mmap[0] = 'a'
	mmap.Flush()
}

func TestMmapRegion(t *testing.T) {
	const pageSize = 4096
	//const pageSize = 65536

	// Create a 2-page sized file
	bigFilePath := filepath.Join(".", "test-data-region")
	fileobj, err := os.OpenFile(bigFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err.Error())
	}

	//init data
	bigData := make([]byte, 2*pageSize, 2*pageSize)
	fileobj.Write(bigData)
	fileobj.Close()

	// Map the first page by itself
	fileobj, err = os.OpenFile(bigFilePath, os.O_RDWR, 0)
	if err != nil {
		panic(err.Error())
	}
	m, err := mmap.MapRegion(fileobj, pageSize, mmap.RDWR, 0, 0)
	if err != nil {
		t.Errorf("error mapping file: %s", err)
	}
	m[0] = 'a'
	m[1] = 'b'
	m[2] = 'c'
	m.Flush()
	//m.Unmap()
	fileobj.Close()

	// Map the second page by itself
	fileobj, err = os.OpenFile(bigFilePath, os.O_RDWR, 0)
	if err != nil {
		panic(err.Error())
	}
	pagesize := os.Getpagesize()
	fmt.Printf("os pageSize is %d\n", pagesize)
	m, err = mmap.MapRegion(fileobj, pageSize, mmap.RDWR, 0, pageSize)
	if err != nil {
		t.Errorf("error mapping file: %s", err)
	}
	m[0] = 'z'
	m[1] = 'y'
	m[2] = 'x'
	m.Flush()
	fileobj.Close()

}

type TestKey struct {
	Key  string
	Hash uint32
}

func NewTestKey(k string) TestKey {
	return TestKey{Key: k, Hash: hash.UInt32Of(k)}
}

func (t TestKey) Less(than btree.Item) bool {
	l, _ := strconv.Atoi(t.Key)
	r, _ := strconv.Atoi(than.(TestKey).Key)
	return l < r
	//return strings.Compare(t.Key, than.(TestKey).Key) < 0
}

//func (t TestKey) Less(than llrb.Item) bool {
//	l, _ := strconv.Atoi(t.Key)
//	r, _ := strconv.Atoi(than.(TestKey).Key)
//	return l < r
//	//return strings.Compare(t.Key, than.(TestKey).Key) < 0
//}

func TestBTree(t *testing.T) {
	var stats runtime.MemStats

	for i := 0; i < 10; i++ {
		runtime.GC()
	}
	fmt.Println("-------- BEFORE ----------")
	runtime.ReadMemStats(&stats)
	fmt.Printf("%+v\n", stats)

	var size = 10000
	tree := btree.New(3)
	for i := 1; i <= size; i++ {
		tree.ReplaceOrInsert(NewTestKey(fmt.Sprintf("%d", i)))
	}
	tree.Delete(NewTestKey("1"))
	tree.ReplaceOrInsert(NewTestKey("1"))
	min := tree.Min()
	fmt.Printf("min %s %d\n", min.(TestKey).Key, min.(TestKey).Hash)
	fmt.Println(tree.Len())

	//tree.Ascend(loop)

	fmt.Println("-------- AFTER ----------")
	runtime.ReadMemStats(&stats)
	fmt.Printf("%+v\n", stats)
	for i := 0; i < 10; i++ {
		runtime.GC()
	}
	fmt.Println("-------- AFTER GC ----------")
	runtime.ReadMemStats(&stats)
	fmt.Printf("%+v\n", stats)

	tree.PrintTree(os.Stdout)
}
func loop(i btree.Item) bool {
	fmt.Printf("key:%s val:%d\n", i.(TestKey).Key, i.(TestKey).Hash)
	return true
}

func lessInt(a, b interface{}) bool {
	return a.(int) < b.(int)
}

//func TestRBTree(t *testing.T) {
//	var stats runtime.MemStats
//
//	for i := 0; i < 10; i++ {
//		runtime.GC()
//	}
//	fmt.Println("-------- BEFORE ----------")
//	runtime.ReadMemStats(&stats)
//	fmt.Printf("%+v\n", stats)
//
//	tree := llrb.New()
//	tree.ReplaceOrInsert(NewTestKey("1"))
//	tree.ReplaceOrInsert(NewTestKey("2"))
//	tree.ReplaceOrInsert(NewTestKey("3"))
//	tree.ReplaceOrInsert(NewTestKey("4"))
//	tree.DeleteMin()
//	tree.Delete(NewTestKey("4"))
//	tree.AscendGreaterOrEqual(NewTestKey("1"),rbtreeLoop)
//
//	fmt.Println("-------- AFTER ----------")
//	runtime.ReadMemStats(&stats)
//	fmt.Printf("%+v\n", stats)
//	for i := 0; i < 10; i++ {
//		runtime.GC()
//	}
//	fmt.Println("-------- AFTER GC ----------")
//	runtime.ReadMemStats(&stats)
//	fmt.Printf("%+v\n", stats)
//}

//func rbtreeLoop(i llrb.Item) bool {
//	fmt.Printf("key:%s val:%d\n", i.(TestKey).Key, i.(TestKey).recoverMeta)
//	return true
//}

func TestBytesPoint(t *testing.T) {
	bytes := []byte{'a', 'b', 'c', 'd', 'e'}
	//len := len(bytes)

	var datPoint *[]byte
	fmt.Printf("%T %v\n", bytes, bytes)
	fmt.Printf("%T %v\n", &bytes, &bytes)
	p := unsafe.Pointer(&bytes)
	datPoint = (*[]byte)(unsafe.Pointer(&bytes))
	fmt.Printf("%T %v\n", p, p)
	fmt.Printf("%T %v\n", datPoint, datPoint)
}

//https://github.com/emirpasic/gods
func TestGodsTrees(t *testing.T) {
	store.NewDataMap(store.DataMapConfig{})
}

//https://github.com/syndtr/goleveldb
func TestGoLevelDb(t *testing.T) {
	var db *leveldb.DB
	db, err := leveldb.OpenFile("data/level-db", nil)
	if err != nil {
		t.Errorf("open leveldb fail")
	}
	defer db.Close()

	//var dat []byte;
	//for i := 0; i < 100000000; i++ {
	//	dat = []byte(fmt.Sprintf("%d", i+1))
	//	db.Put(dat, dat, nil)
	//}

	value, _ := db.Get([]byte("1"), nil)
	fmt.Println(string(value))
}

//https://github.com/dgraph-io/badger
//github.com/dgraph-io/badger v1.6.0
func TestBadger(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions("data/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(txn *badger.Txn) error {
		var dat []byte;
		for i := 0; i < 1000000; i++ {
			dat = []byte(fmt.Sprintf("%d", i+1))
			txn.Set(dat, dat)
		}
		return nil
	})
}

//https://github.com/go-ego/riot
//github.com/go-ego/riot v0.0.0-20190307162011-3d971d90bc83
//func TestRiot(t *testing.T) {
//	searcher := riot.Engine{}
//	// 初始化
//	searcher.Init(types.EngineOpts{
//		UseStore:    true,
//		StoreFolder: "data/riot",
//		StoreEngine: "ldb",
//		//StoreEngine: "bg",
//		StoreShards: 3,
//		PinYin:      true,
//		Using:       3,
//		GseDict:     "zh",
//		// GseDict: "your gopath"+"/src/github.com/go-ego/riot/data/dict/dictionary.txt",
//	})
//	defer searcher.Close()
//
//	text := "《复仇者联盟3：无限战争》是全片使用IMAX摄影机拍摄"
//	text1 := "在IMAX影院放映时"
//	text2 := "全片以上下扩展至IMAX 1.9：1的宽高比来呈现"
//
//	// 将文档加入索引，docId 从1开始
//	searcher.Index("1", types.DocData{Content: text})
//	searcher.Index("2", types.DocData{Content: text1}, false)
//	searcher.Index("3", types.DocData{Content: text2}, true)
//
//	// 等待索引刷新完毕
//	searcher.Flush()
//	// engine.FlushIndex()
//
//	// 搜索输出格式见 types.SearchResp 结构体
//	log.Print(searcher.Search(types.SearchReq{Text: "fc"}))
//}

//github.com/boltdb/bolt v1.3.1
func TestBoltDb(t *testing.T) {

}

func launchCluster(port int) {
	logger := logger.NewLogger(logger.LoggerLevelDebug, "./", fmt.Sprintf("log-%d.log", port))

	logger.Debugf("cluster [%d] launch", port)
	config := memberlist.DefaultLocalConfig()
	config.Name = strconv.Itoa(port)
	config.SecretKey = []byte("abcdef0123456789")
	config.BindPort = port
	config.AdvertisePort = port
	config.EnableCompression = true

	//config.Transport
	//config.Delegate = &MessageDelegate{}
	list, err := memberlist.Create(config)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	// Join an existing cluster by specifying at least one known member.
	n, err := list.Join([]string{"192.168.1.106:1107", "192.168.1.106:1108"})
	if err != nil {
		panic("Failed to join cluster: " + err.Error())
	}

	fmt.Println(n)
	logger.Debugf("cluster [%d] launch completed", port)
	for {
		// Ask for members of the cluster
		members := list.Members()
		fmt.Printf("[%d]has members %d\n", port, len(members))
		for _, member := range members {
			fmt.Printf("[%d]Member: %s %s\n", port, member.Name, member.Addr)
			err := list.SendToTCP(member, []byte("hi"))
			if err != nil {
				fmt.Println(err)
			}
		}
		time.Sleep(time.Second * 3)
	}
}

func TestGossip(t *testing.T) {
	exit := make(chan bool)
	go func() {
		launchCluster(1107)
	}()
	go func() {
		launchCluster(1108)
	}()
	go func() {
		launchCluster(1109)
	}()
	go func() {
		launchCluster(1110)
	}()
	go func() {
		launchCluster(1111)
	}()

	<-exit
}

type MessageDelegate struct {
}

// NodeMeta is used to retrieve meta-data about the current node
// when broadcasting an alive message. It's length is limited to
// the given byte size. This metadata is available in the Node structure.
func (p *MessageDelegate) NodeMeta(limit int) []byte {
	fmt.Printf("NodeMeta [%d]\n", limit)
	return nil
}

// NotifyMsg is called when a user-data message is received.
// Care should be taken that this method does not block, since doing
// so would block the entire UDP packet receive loop. Additionally, the byte
// slice may be modified after the call returns, so it should be copied if needed
func (p *MessageDelegate) NotifyMsg(dat []byte) {
	fmt.Printf("received [%s]\n", string(dat))
}

// GetBroadcasts is called when user data messages can be broadcast.
// It can return a list of buffers to send. Each buffer should assume an
// overhead as provided with a limit on the total byte size allowed.
// The total byte size of the resulting data to send must not exceed
// the limit. Care should be taken that this method does not block,
// since doing so would block the entire UDP packet receive loop.
func (p *MessageDelegate) GetBroadcasts(overhead, limit int) [][]byte {
	fmt.Printf("GetBroadcasts [%d/%d]\n", overhead, limit)
	return nil
}

// LocalState is used for a TCP push/pull. This is sent to
// the remote side in addition to the membership information. Any
// data can be sent here. See MergeRemoteState as well. The `join`
// boolean indicates this is for a join instead of a push/pull.
func (p *MessageDelegate) LocalState(join bool) []byte {
	fmt.Printf("LocalState [%v]\n", join)
	return nil
}

// MergeRemoteState is invoked after a TCP push/pull. This is the
// state received from the remote side and is the result of the
// remote side's LocalState call. The 'join'
// boolean indicates this is for a join instead of a push/pull.
func (p *MessageDelegate) MergeRemoteState(buf []byte, join bool) {
	fmt.Printf("MergeRemoteState [%s|%v]\n", string(buf), join)
}

