#!/usr/bin/env bash
cd build
wget -q http://downloads.arduino.cc/arduino-1.8.8-linux64.tar.xz
tar xf arduino-1.8.8-linux64.tar.xz
sudo mv arduino-1.8.8 /usr/local/share/arduino
sudo apt -qq install -y cmake >> /dev/null
rm -rf arduino*
