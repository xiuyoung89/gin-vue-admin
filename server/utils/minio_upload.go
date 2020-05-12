package utils

import (
	"fmt"
	"gin-vue-admin/global"
	"github.com/minio/minio-go/v6"
	"log"
	"mime/multipart"
	"time"
)

// 接收两个参数 一个文件流 一个 bucket 你的七牛云标准空间的名字
func MinioUpload(file *multipart.FileHeader) (err error, path string, key string) {

	// checkout bucket
	endpoint := global.GVA_CONFIG.Minio.EndPoint
	accessKeyID := global.GVA_CONFIG.Minio.AccessKey
	secretAccessKey := global.GVA_CONFIG.Minio.SecretKey
	bucketName := global.GVA_CONFIG.Minio.Bucket
	location := global.GVA_CONFIG.Minio.Location
	useSSL := global.GVA_CONFIG.Minio.SSL

	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Panicf("cant create minio bucket: %v", err)
	}
	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// 检查存储桶是否已经存在。
		exists, err := minioClient.BucketExists(bucketName)
		if err == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalf("create bucket error: %v",err)
		}
	}

	// upload file

	f, e := file.Open()
	if e != nil {
		fmt.Println(e)
		return e, "", ""
	}
	defer f.Close()

	dataLen := file.Size
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename) // 文件名格式 自己可以改 建议保证唯一性
	n, err := minioClient.PutObject(bucketName, fileKey, f, dataLen, minio.PutObjectOptions{ContentType:"image/png"})
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("minio put object(%s) success, size: %d", fileKey, n)

	//putPolicy := storage.PutPolicy{
	//	Scope: global.GVA_CONFIG.Minio.Bucket,
	//}
	//mac := qbox.NewMac(global.GVA_CONFIG.Qiniu.AccessKey, global.GVA_CONFIG.Qiniu.SecretKey)
	//upToken := putPolicy.UploadToken(mac)
	//cfg := storage.Config{}
	//// 空间对应的机房
	//cfg.Zone = &storage.ZoneHuadong
	//// 是否使用https域名
	//cfg.UseHTTPS = false
	//// 上传是否使用CDN上传加速
	//cfg.UseCdnDomains = false
	//formUploader := storage.NewFormUploader(&cfg)
	//ret := storage.PutRet{}
	//putExtra := storage.PutExtra{
	//	Params: map[string]string{
	//		"x:name": "github logo",
	//	},
	//}
	//f, e := file.Open()
	//if e != nil {
	//	fmt.Println(e)
	//	return e, "", ""
	//}
	//dataLen := file.Size
	//fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename) // 文件名格式 自己可以改 建议保证唯一性
	//err = formUploader.Put(context.Background(), &ret, upToken, fileKey, f, dataLen, &putExtra)
	//if err != nil {
	//	global.GVA_LOG.Error("upload file fail:", err)
	//	return err, "", ""
	//}
	//return err, global.GVA_CONFIG.Qiniu.ImgPath + "/" + ret.Key, ret.Key
	path = global.GVA_CONFIG.Minio.EndPoint + "/" + bucketName + "/" + fileKey
	return
}


func MinioDeleteFile(key string) (err error) {

	endpoint := global.GVA_CONFIG.Minio.EndPoint
	accessKeyID := global.GVA_CONFIG.Minio.AccessKey
	secretAccessKey := global.GVA_CONFIG.Minio.SecretKey
	bucketName := global.GVA_CONFIG.Minio.Bucket
	//location := global.GVA_CONFIG.Minio.Location
	useSSL := global.GVA_CONFIG.Minio.SSL

	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Printf("remove bucket: %s, object: %s, error: %v",bucketName, key,  err)
		return
	}

	err = minioClient.RemoveObject(bucketName, key)
	if err != nil {
		fmt.Println(err)
		return
	}
	return

	//mac := qbox.NewMac(global.GVA_CONFIG.Qiniu.AccessKey, global.GVA_CONFIG.Qiniu.SecretKey)
	//cfg := storage.Config{
	//	// 是否使用https域名进行资源管理
	//	UseHTTPS: false,
	//}
	//// 指定空间所在的区域，如果不指定将自动探测
	//// 如果没有特殊需求，默认不需要指定
	////cfg.Zone=&storage.ZoneHuabei
	//bucketManager := storage.NewBucketManager(mac, &cfg)
	//err := bucketManager.Delete(global.GVA_CONFIG.Qiniu.Bucket, key)
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}
	return nil
}
