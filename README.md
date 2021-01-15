# xkcdDL

A tiny utility to help you concurrently download memes(images) from [xkcd.com](https://xkcd.com/) by specifying the number of images to download. Default is the 10 most current images

## Install

go get -u -x github.com/vanderkilu/xkcdDL

## Usage

    xkcdDL [options]

## Options

- `-d` The directory(relative) to save the downloaded images. Default is currently directory the tool is running in.
- `-w` The number of workers to concurrently use to download the images. Default is 10.
- `-m` The number of images to download from xkcd.com. Default is the 10 most current images

## Examples

`xkcdDL` //uses default values

`xkcdDL -d /home/vndrkl/xkcd -w 20 -m 100
