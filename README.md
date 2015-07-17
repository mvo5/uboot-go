# Read/write uboot environment

Small go package/app to read/write uboot env files that contain crc32 + 1 byte
padding. Unlike fw_{set,print}env it does not needs a
/etc/fw_env.config config file.

Example of the go app:
```
$ ./uboot-env uboot.env print
initrd_addr=0x88080000
uenvcmd=load mmc ${bootpart} ${loadaddr} snappy-system.txt; env import -t $loadaddr $filesize; run snappy_boot
bootpart=0:1

$ ./uboot-env uboot.env set key value
$ ./uboot-env uboot.env print
initrd_addr=0x88080000
uenvcmd=load mmc ${bootpart} ${loadaddr} snappy-system.txt; env import -t $loadaddr $filesize; run snappy_boot
bootpart=0:1
key=value

# echo "$(pwd)/uboot.env 0x000 0x20000" > /etc/fw_env.config
$ fw_printenv
initrd_addr=0x88080000
uenvcmd=load mmc ${bootpart} ${loadaddr} snappy-system.txt; env import -t $loadaddr $filesize; run snappy_boot
bootpart=0:1
key=value
```

Example of the API:
```
package main

import (
	"fmt"
	"github.com/mvo5/uboot-env/uboot"
	"os"
)

func main() {
	env, _ := uboot.NewEnv(os.Args[1])
	fmt.Print(env)
	env.Set("foo", "bar")
	fmt.Print(env)
}
```