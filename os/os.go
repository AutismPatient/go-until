package os2

import "github.com/StackExchange/wmi"

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
	获取磁盘信息 WIN
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
