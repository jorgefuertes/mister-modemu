~~~
 ______         __               ________ __ __     __
|   __ \.-----.|  |_.----.-----.|  |  |  |__|  |--.|__|
|      <|  -__||   _|   _|  _  ||  |  |  |  |    < |  |
|___|__||_____||____|__| |_____||________|__|__|__||__|
_______________________________________________________

      ESP8266 AT Modem Emulator for ZX-Next core
_______________________________________________________
~~~~

## What is

This is a modem emulator to be instaled into a
[Mister FPGA board](https://misterfpga.org).

It allows ZX-Next, or any other core that can use a serial port and
expect an ESP8266 WiFi module, to connect to internet, as in the
real machine with the WiFi module.

The modem emulator is a work in progress but has the basics and can
be used with [NXTel](https://github.com/Threetwosevensixseven/NXtel)
smoothly to connect with Next's Videotex sites.

## Requisites

- A Mister complete system, you can get one from [ManuFerHi](https://manuferhi.com/) for example.
- Internet connection through mister's ethernet or USB WiFi link.
- RetroWiki's [ZXNext_Mister](https://github.com/benitoss/ZXNext_Mister) core running.
- SD Card or VHD boot file with the [basic system distro](https://www.specnext.com/latestdistro).

## Recommended

- Latest [NXTel release](https://github.com/Threetwosevensixseven/NXtel/releases).

## Installation

SSH as `root` user to your Mister an type or paste:

~~~bash
curl -kLs http://bit.ly/modemu-install | bash
~~~

That's all.

## Usage

- Navigate to Mister's system menu, pressing `ESC` at the core's menu.
- Select `Scripts`, and `Yes` to continue.
- Select the `RetroWiki_Modemu_Start` script. This script knows how to install, run and update this modem emulator.
- Now you can load ZXNext's core and browse to `demos`.
- Execute `NXTel` and connect to any service.
