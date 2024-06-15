package main

import (
	"image/png"
	"os"

	"github.com/disintegration/imaging"
)

func main() {
	// 打开一个图片文件
	inFile, err := os.Open("avatar.png")
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	// 解码图片
	img, err := png.Decode(inFile)
	if err != nil {
		panic(err)
	}

	// 创建一个新的图片，并将其尺寸设置为原图的一半

	// 创建输出文件
	outFile, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	err = imaging.Encode(outFile, img, imaging.PNG, imaging.PNGCompressionLevel(png.BestCompression))
	if err != nil {
		panic(err)
	}
}
