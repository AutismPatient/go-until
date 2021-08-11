package oss

// oss 服务方法

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go-until/date"
	string2 "go-until/string"
	"io/ioutil"
	"mime/multipart"
	"time"
)

type ClientStruct struct {
	*oss.Client
}

// 上传选择类型
const (
	SIMPLE = iota
	TADD
	BREAKPOINT
	BURST
)

var (
	Client *ClientStruct

	ConfigValue map[string]string
)

func init() {
	Client = clientInit()
}
func printVersion() {
	fmt.Println("OSS Go SDK Version: ", oss.Version)
}

// 初始化OSS客户端
func clientInit() *ClientStruct {
	client, err := oss.New(ConfigValue["Endpoint"], ConfigValue["Access-Key-ID"], ConfigValue["Access-Key-Secret"],
		oss.Timeout(12, 120),
		oss.EnableCRC(false),
		oss.EnableMD5(true), oss.RedirectEnabled(true))
	if err != nil {
		panic(err)
	}
	defer printVersion()
	return &ClientStruct{client}
}

// 获取空间实例
func (client *ClientStruct) getBucket(name string) *oss.Bucket {
	bucket, err := client.Bucket(name)
	if err != nil {
		panic(err)
	}
	return bucket
}

// CreateBucketOfVersion 创建空间，根据版本
func (client *ClientStruct) CreateBucketOfVersion(name string, isVersion bool) bool {
	err := client.CreateBucket(name)
	if err != nil {
		panic(err)
	}
	if isVersion {
		// 设置存储空间版本控制状态为Enabled或Suspended，此处以设置存储空间版本状态为Enabled为例。
		config := oss.VersioningConfig{Status: "Enabled"}
		err = client.SetBucketVersioning(name, config)
		if err != nil {
			panic(err)
		}
	}
	return true
}

// GetBucketInfo 获取空间版本状态信息
func (client *ClientStruct) GetBucketInfo(name string) oss.GetBucketVersioningResult {
	versionInfo, err := client.GetBucketVersioning(name)
	if err != nil {
		panic(err)
	}
	return versionInfo
}

// Upload 上传
func (client *ClientStruct) Upload(uploadType int, suffix string, data multipart.File) (bool, string) {
	switch uploadType {
	case 0:
		return client.uploadSimple(suffix, data)
		//case 2:
		//	client.uploadToAdd(bucket, data)
		//case 3:
		//	client.uploadBreakPoint(bucket, data)
		//case 4:
		//	client.uploadBurst(bucket, data)
	}
	return false, ""
}

// 生成ObjectKey | 时间戳 + 随机大小写字母
func createObjectKey() string {
	var (
		Unix = time.Now().Unix()
		Str  = string2.Helper.GetRandomString(6)
	)
	return fmt.Sprintf("%d%s", Unix, Str)
}

// 简单上传
func (client *ClientStruct) uploadSimple(suffix string, file multipart.File) (bool, string) {
	bucket := client.getBucket(ConfigValue["Bucket"])
	key := fmt.Sprintf("%s/%s.%s", date.GetSortTimeText(time.Now().Unix()), createObjectKey(), suffix)
	// 上传文件流
	defer file.Close()
	err := bucket.PutObject(key, file)
	if err != nil {
		panic(err)
	}
	fmt.Println("image upload ok")
	return true, key
}

// GetObjectSoftLink 获取Object资源
func (client *ClientStruct) GetObjectSoftLink(key string) (data []byte) {
	bucket := client.getBucket(ConfigValue["Bucket"])
	if client.IsExist(key) {
		body, err := bucket.GetObject(key)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(body)
		defer body.Close()
		if err != nil {
			panic(err)
		}
	}
	return
}

// DeleteObject 删除资源
func (client *ClientStruct) DeleteObject(key string) error {
	bucket := client.getBucket(ConfigValue["Bucket"])
	return bucket.DeleteObject(key)
}

// 判断文件是否存在
func (client *ClientStruct) IsExist(key string) bool {
	bucket := client.getBucket(ConfigValue["Bucket"])
	isExist, err := bucket.IsObjectExist(key)
	if err != nil {
		panic(err)
	}
	return isExist
}

// 定义进度条监听器。
type OssProgressListener struct {
}

// 定义进度变更事件处理函数。
func (listener *OssProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferStartedEvent:
		fmt.Printf("Transfer Started, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferDataEvent:
		fmt.Printf("\rTransfer Data, ConsumedBytes: %d, TotalBytes %d, %d%%.",
			event.ConsumedBytes, event.TotalBytes, event.ConsumedBytes*100/event.TotalBytes)
	case oss.TransferCompletedEvent:
		fmt.Printf("\nTransfer Completed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferFailedEvent:
		fmt.Printf("\nTransfer Failed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	default:
	}
}

//// 追加上传 TODO 2020年9月8日22:33:12
//func (client *OssClientStruct) uploadToAdd(name string, file *os.File) {
//	var nextPos int64 = 0
//	var bucket = client.getBucket(name)
//	// 第一次追加的位置是0，返回值为下一次追加的位置。后续追加的位置是追加前文件的长度。
//	nextPos, err := bucket.AppendObject(bucket.BucketName, bytes.NewReader(data), nextPos)
//	if err != nil {
//		fmt.Println("Error:", err)
//		panic(err)
//	}
//
//	// 第二次追加。
//	nextPos, err = bucket.AppendObject(bucket.BucketName, bytes.NewReader(data), nextPos)
//	if err != nil {
//		fmt.Println("Error:", err)
//		panic(err)
//	}
//
//	// ....
//}
//
//// 断点续传上传 TODO
//func (client *OssClientStruct) uploadBreakPoint(name string, file *os.File) {
//	// 设置分片大小为100 KB，指定分片上传并发数为3，并开启断点续传上传。
//	// 其中<yourObjectName>与objectKey是同一概念，表示断点续传上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
//	// "LocalFile"为filePath，100*1024为partSize。
//	var bucket = client.getBucket(name)
//	err := bucket.UploadFile(bucket.BucketName, "LocalFile", 100*1024, oss.Routines(3), oss.Checkpoint(true, ""))
//	if err != nil {
//		fmt.Println("Error:", err)
//		panic(err)
//	}
//}
//
//// 分片上传 TODO
//func (client *OssClientStruct) uploadBurst(name string, file *os.File) {
//	bucket := client.getBucket(name)
//	objectName := "<yourObjectName>"
//	locaFilename := "<yourLocalFilename>"
//	chunks, err := oss.SplitFileByPartNum(locaFilename, 3)
//	fd, err := os.Open(locaFilename)
//	defer fd.Close()
//
//	// 指定存储类型为标准存储。
//	storageType := oss.ObjectStorageClass(oss.StorageStandard)
//
//	// 步骤1：初始化一个分片上传事件，并指定存储类型为标准存储。
//	imus, err := bucket.InitiateMultipartUpload(objectName, storageType)
//	// 步骤2：上传分片。
//	var parts []oss.UploadPart
//	for _, chunk := range chunks {
//		fd.Seek(chunk.Offset, os.SEEK_SET)
//		// 调用UploadPart方法上传每个分片。
//		part, err := bucket.UploadPart(imus, fd, chunk.Size, chunk.Number)
//		if err != nil {
//			fmt.Println("Error:", err)
//			os.Exit(-1)
//		}
//		parts = append(parts, part)
//	}
//
//	// 指定Object的读写权限为公共读，默认为继承Bucket的读写权限。
//	objectAcl := oss.ObjectACL(oss.ACLPublicRead)
//
//	// 步骤3：完成分片上传，指定文件读写权限为公共读。
//	cmdr, err := bucket.CompleteMultipartUpload(imus, parts, objectAcl)
//	if err != nil {
//		fmt.Println("Error:", err)
//		os.Exit(-1)
//	}
//	fmt.Println(cmdr)
//}
