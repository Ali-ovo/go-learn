package main

import (
	"log"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var client *oss.Client // 全局变量用来存储OSS客户端实例

func main() {
	// yourBucketName填写存储空间名称。
	bucketName := ""
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	endpoint := ""

	// 检查环境变量是否已经设置。
	if endpoint == "" || bucketName == "" {
		log.Fatal("Please set yourEndpoint and bucketName.")
	}

	// 从环境变量中获取访问凭证。
	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		handleError(err)
	}

	// 创建OSSClient实例。
	client, err = oss.New(endpoint, "", "", oss.SetCredentialsProvider(&provider))
	if err != nil {
		handleError(err)
	}

	// 输出客户端信息。
	log.Printf("Client: %#v\n", client)

	// 示例操作：创建存储空间。
	// if err := createBucket(bucketName); err != nil {
	// 	handleError(err)
	// }

	// 示例操作：上传文件。
	objectName := "cat.jpeg"
	localFileName := "/Users/ctw/Pictures/my/cat.jpeg"
	if err := uploadFile(bucketName, objectName, localFileName); err != nil {
		handleError(err)
	}

	// // 示例操作：下载文件。
	// downloadedFileName := "/Users/ctw/Downloads"
	// if err := downloadFile(bucketName, objectName, downloadedFileName); err != nil {
	// 	handleError(err)
	// }

	// // 示例操作：列举对象。
	// if err := listObjects(bucketName); err != nil {
	// 	handleError(err)
	// }

	// // 示例操作：删除对象。
	// if err := deleteObject(bucketName, objectName); err != nil {
	// 	handleError(err)
	// }
}

// handleError 用于处理不可恢复的错误，并记录错误信息后终止程序。
func handleError(err error) {
	log.Fatalf("Error: %v", err)
}

// createBucket 用于创建一个新的OSS存储空间。
// 参数：
//
//	bucketName - 存储空间名称。
//	endpoint - Bucket对应的Endpoint。
//
// 如果成功，记录成功日志；否则，返回错误。
func createBucket(bucketName string) error {
	// 创建存储空间。
	err := client.CreateBucket(bucketName)
	if err != nil {
		return err
	}

	// 存储空间创建成功后，记录日志。
	log.Printf("Bucket created successfully: %s", bucketName)
	return nil
}

// uploadFile 用于将本地文件上传到OSS存储桶。
// 参数：
//
//	bucketName - 存储空间名称。
//	objectName - Object完整路径，完整路径中不包含Bucket名称。
//	localFileName - 本地文件的完整路径。
//	endpoint - Bucket对应的Endpoint。
//
// 如果成功，记录成功日志；否则，返回错误。
func uploadFile(bucketName, objectName, localFileName string) error {
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// 上传文件。
	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		return err
	}

	// 文件上传成功后，记录日志。
	log.Printf("File uploaded successfully to %s/%s", bucketName, objectName)
	return nil
}

// downloadFile 用于从OSS存储桶下载一个文件到本地路径。
// 参数：
//
//	bucketName - 存储空间名称。
//	objectName - Object完整路径，完整路径中不能包含Bucket名称。
//	downloadedFileName - 本地文件的完整路径。
//	endpoint - Bucket对应的Endpoint。
//
// 如果成功，记录成功日志；否则，返回错误。
func downloadFile(bucketName, objectName, downloadedFileName string) error {
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// 下载文件。
	err = bucket.GetObjectToFile(objectName, downloadedFileName)
	if err != nil {
		return err
	}

	// 文件下载成功后，记录日志。
	log.Printf("File downloaded successfully to %s", downloadedFileName)
	return nil
}

// listObjects 用于列举OSS存储空间中的所有对象。
// 参数：
//
//	bucketName - 存储空间名称。
//	endpoint - Bucket对应的Endpoint。
//
// 如果成功，打印所有对象；否则，返回错误。
func listObjects(bucketName string) error {
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// 列举文件。
	marker := ""
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker))
		if err != nil {
			return err
		}

		// 打印列举文件，默认情况下一次返回100条记录。
		for _, object := range lsRes.Objects {
			log.Printf("Object: %s", object.Key)
		}

		if !lsRes.IsTruncated {
			break
		}
		marker = lsRes.NextMarker
	}

	return nil
}

// deleteObject 用于删除OSS存储空间中的一个对象。
// 参数：
//
//	bucketName - 存储空间名称。
//	objectName - 要删除的对象名称。
//	endpoint - Bucket对应的Endpoint。
//
// 如果成功，记录成功日志；否则，返回错误。
func deleteObject(bucketName, objectName string) error {
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// 删除文件。
	err = bucket.DeleteObject(objectName)
	if err != nil {
		return err
	}

	// 文件删除成功后，记录日志。
	log.Printf("Object deleted successfully: %s/%s", bucketName, objectName)
	return nil
}
