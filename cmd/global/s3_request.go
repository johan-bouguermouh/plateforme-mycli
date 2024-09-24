package global

import (
	env "bucketool/environment"
	color "bucketool/utils/colorPrint"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func BucketExists(bucketName string) bool {
	if env.IsDebugMode {
		println(color.GreyP("Checking if bucket exists : "), color.GreenP(bucketName))
	}
	_, err := Connexion.S3Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if(err != nil && env.IsDebugMode){
		println(color.RedP("Bucket " + bucketName + " not found"))
		println(color.RedP(err.Error()))
	}
	return err == nil
}