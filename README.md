[![Build Status][travis-image]][travis-url] 
# Read/write uboot environment

Small go package/app to read/write uboot env files that contain crc32 + 1 byte
padding. Unlike fw_{set,print}env it does not needs a
/etc/fw_env.config config file.

Example of the API:
```
package main

import (
	"fmt"
	"github.com/mvo5/uboot-go/uenv"
	"os"
)

func main() {
	env, _ := uenv.Open(os.Args[1])
	fmt.Print(env)
	env.Set("foo", "bar")
	fmt.Print(env)
}
```

Example of the cmdline app for existing files:
```
$ uboot-go uboot.env print
initrd_addr=0x88080000
uenvcmd=load mmc ${bootpart} ${loadaddr} snappy-system.txt; env import -t $loadaddr $filesize; run snappy_boot
bootpart=0:1

$ uboot-go uboot.env set key value
$ uboot-go uboot.env print
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

Example of the cmdline app for creating new env files:
```
$ uboot-go uboot.env create 4096
$ uboot-go uboot.env set foo bar
$ uboot-go uboot.env print
foo=bar
```

[travis-image]: https://travis-ci.org/mvo5/uboot-go.svg?branch=master
[travis-url]: https://travis-ci.org/mvo5/uboot-go
