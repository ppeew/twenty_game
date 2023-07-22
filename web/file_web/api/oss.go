package api

//type OSS struct {
//	bucket *oss.Bucket
//}

//func NewOSS() (*OSS, error) {
//	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
//	endpoint := "oss-cn-guangzhou.aliyuncs.com"
//	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
//	accessKeyId := ""
//	accessKeySecret := ""
//	// yourBucketName填写Bucket名称。
//	bucketName := "twelve-files"
//	// 创建OSSClient实例。
//	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
//	if err != nil {
//		return nil, err
//	}
//	// 创建存储空间。
//	err = client.CreateBucket(bucketName)
//	if err != nil {
//		return nil, err
//	}
//	bucket, err := client.Bucket(bucketName)
//	if err != nil {
//		return nil, err
//	}
//	return &OSS{bucket: bucket}, nil
//}
//
//func (oss *OSS) UploadFile(objectKey, filePath string) {
//	//创文件
//	err := oss.bucket.PutObjectFromFile(objectKey, filePath)
//	if err != nil {
//		panic(err)
//	}
//}
//
//func (oss *OSS) DownloadFile(objectKey, descPath string) {
//	err := oss.bucket.GetObjectToFile(objectKey, "./first.jpg")
//	if err != nil {
//		panic(err)
//	}
//}
