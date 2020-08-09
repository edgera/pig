# PI Gadget
A repo for performing windows 10 input via usb from raspberry pi 4



## Resources this repo is built on
https://threadsec.wordpress.com/raspberry-pi-zero-usb-composite-gadget/
https://blog.gbaman.info/?p=791
https://blog.gbaman.info/?p=699
https://www.collabora.com/news-and-blog/blog/2019/06/24/using-dummy-hcd/
https://github.com/darrylburke/RaspberryPiZero_HID_MultiTool/tree/master/gadget/hid
http://irq5.io/2016/12/22/raspberry-pi-zero-as-multiple-usb-gadgets/
https://www.isticktoit.net/?p=1383

# Raspberry pi setup
Some combination of the following seemed to make stuff work.

## dwc2
dwc2 is an upstream driver which can do the OTG host/gadget flip dictated by OTG_SENSE.
In host mode performance will pale cf dwc_otg, hence it's only recommended for gadget mode.
```
echo "dtoverlay=dwc2" | sudo tee -a /boot/config.txt
echo "dwc2" | sudo tee -a /etc/modules
```

## libcomposite
```
sudo echo "libcomposite" | sudo tee -a /etc/modules

```

## libusbgx (optional?)
```
git clone https://github.com/libusbgx/libusbgx.git
cd libusbgx
autoreconf -i
./configure  --prefix=/usr
make
make install # as root
```

## gt (optional?)
```
git clone https://github.com/kopasiak/gt.git
cd gt/source
cmake -DCMAKE_INSTALL_PREFIX= .
make
make install # as root
```

## Checking kernel config (note)
```
sudo modprobe configs   # Will create file /proc/config.gz
zcat /proc/config.gz > .config
```


