package os2

import (
	"fmt"
	"testing"
)

/*
	[{C: NTFS 107376275456 33885372416} {D: NTFS 404283002880 271718969344}]
*/
func TestGetStorageInfo(t *testing.T) {
	fmt.Println(GetStorageInfo())
}
