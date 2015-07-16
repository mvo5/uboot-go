# Read/write uboot environment

Small go app to read/write uboot env files that contain crc + 1 byte
padding. Its more flexible than fw_{set,print}env that needs a
/etc/fw_env.config config file.


Example:
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