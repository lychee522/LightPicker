package utils

// @author 肖肖雨歇
// @description 工具类：计算文件哈希值，为秒传功能提供弹药

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime/multipart"
)

// CalcFileHash 计算上传文件的 MD5 特征值
func CalcFileHash(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	hash := md5.New()
	// 将文件流拷贝进哈希计算器，不占用过多内存
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
