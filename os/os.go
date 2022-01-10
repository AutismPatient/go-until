package os2

import (
	"fmt"
	"github.com/AutismPatient/go-until/driver"
	string2 "github.com/AutismPatient/go-until/string"
	"github.com/StackExchange/wmi"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"runtime"
)

type Storage struct {
	Name       string
	FileSystem string
	Total      uint64
	Free       uint64
}

type storageInfo struct {
	Name       string
	Size       uint64
	FreeSpace  uint64
	FileSystem string
}

/*
	获取磁盘信息 仅WIN生效
*/
func GetStorageInfo() []Storage {
	var (
		storageInfo   []storageInfo
		localStorages []Storage
	)
	err := wmi.Query("Select * from Win32_LogicalDisk", &storageInfo)
	if err != nil {
		panic(err)
	}
	for _, storage := range storageInfo {
		info := Storage{
			Name:       storage.Name,
			FileSystem: storage.FileSystem,
			Total:      storage.Size,
			Free:       storage.FreeSpace,
		}
		localStorages = append(localStorages, info)
	}
	return localStorages
}

type SystemInfo struct {
	HostName        string          `json:"host_name"`
	PageSize        int             `json:"page_size"`
	NumCPU          int             `json:"num_cpu"`
	NumCgoCall      int64           `json:"num_cgo_call"`
	Version         string          `json:"version"`
	OS              string          `json:"os"`
	GOARCH          string          `json:"goarch"`
	InterfaceAdders []net.Addr      `json:"interface_adders"`
	Interfaces      []net.Interface `json:"interfaces"`
}

/*
	返回系统信息
*/
func GetSystemInfo() (info SystemInfo) {
	var (
		err error
	)
	info.HostName, err = os.Hostname()
	if err != nil {
		panic(err)
	}
	info.NumCPU = runtime.NumCPU()
	info.NumCgoCall = runtime.NumCgoCall()
	info.Version = runtime.Version()
	info.OS = runtime.GOOS
	info.GOARCH = runtime.GOARCH
	info.Interfaces, err = net.Interfaces()
	if err != nil {
		panic(err)
	}
	info.InterfaceAdders, err = net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	return info
}

/*
	上传批量文件（支持断点续传） todo 2020年12月15日21:48:39
*/
func UploadFile(req *http.Request, key string) (info map[string]string, err error) {
	multipartForm := req.MultipartForm
	files := multipartForm.File[key]
	for i := 0; i < len(files); i++ {
		fileName := files[i].Filename
		file, err := files[i].Open()
		if err != nil {
			info[fileName] = "文件打开错误：" + err.Error()
			continue
		}
		err = doUpload(file, fileName)
	}
	fmt.Println(multipartForm)
	return info, err
}

/*
	临时文件结构体
*/
type TempFileInfo struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Position int    `json:"position"` // 上传偏移量
}

func doUpload(file multipart.File, name string) (err error) {
	fileInfo := TempFileInfo{
		Id:       string2.Helper.CreateFileHash(file),
		Name:     name,
		Position: 0,
	}
	rely, err := driver.RedisClient.Do("get", fileInfo)
	fmt.Println(rely)
	defer file.Close()
	return err
}
