package imageOpt

import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	"image/jpeg"
	"mime/multipart"
	"path/filepath"
	"strings"
)

func FileHeaderToImage(header *multipart.FileHeader) (image.Image, string, error) {
	file, err := header.Open()
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	fileExt := strings.ToLower(filepath.Ext(header.Filename))
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, "", err
	}

	return img, fileExt, nil
}

func Cropping(img image.Image, croppedX, croppedY, width, height int) ([]byte, error) {
	// 检查参数是否有效
	if croppedX < 0 || croppedY < 0 || width <= 0 || height <= 0 {
		return nil, errors.New("invalid parameters")
	}

	// 检查裁剪区域是否超出原图的范围
	if croppedX+width > img.Bounds().Dx() || croppedY+height > img.Bounds().Dy() {
		return nil, errors.New("crop area is out of image bounds")
	}

	// 创建裁剪区域
	cropRect := image.Rect(0, 0, width, height)

	// 创建一个新的图像，大小和裁剪区域相同
	croppedImg := image.NewRGBA(cropRect)

	// 将原图像的内容复制到新图像中，使用 cropRect 作为目标区域
	draw.Draw(croppedImg, cropRect, img, image.Point{croppedX, croppedY}, draw.Src)

	// 创建一个字节数组缓冲区来存储图像数据
	var buffer bytes.Buffer

	// 将正方形图像编码为JPEG格式并将结果写入缓冲区
	err := jpeg.Encode(&buffer, croppedImg, nil)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
