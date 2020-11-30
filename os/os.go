package os2

import (
	"fmt"
	"github.com/StackExchange/wmi"
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
	上传批量文件（支持断点续传）
*/
func UploadFile(req *http.Request) {
	multipartForm := req.MultipartForm
	fmt.Println(multipartForm)
}
