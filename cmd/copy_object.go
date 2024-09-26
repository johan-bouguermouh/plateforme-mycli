package cmd

import (
	global "bucketool/cmd/global"
	utils "bucketool/utils"
	color "bucketool/utils/colorPrint"
	"context"
	"log"
	"strings"

	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/urfave/cli"
)

// Description hiw to use the command
var registryUsageText string = `This command copy a file from a path and insert it in a bucket destination.
	The first argument is the path of the file to copy.
	The flag -d is the destination bucket.
	The flag -n is the name of the file in the bucket.
	The flag -alias can be used to specify the alias to use, if you have specified the current alias, you can omit this flag.
	Usage of flag -alias is optional. If you use it, this flags must be placed before command, like this : -alias <alias> copy <path> \n`+
	"Example : copy /path/to/file -d mybucket -n myfile.txt" + "\n" + "Example : -alias myalias copy /path/to/file -d mybucket -n myfile.txt" + "\n"

// Description of the flags
var createBucketDesc string = color.ColorPrint("Black",
	" Flags Descriptions ", &color.Options{
		Underline: true,
		Bold:      true,
		Background: "White",
	}) + "\n" + 
	color.BlueP(" -d, --destination ") +  color.GreyP("(Required)")+" : Destination bucket\n" +
	color.BlueP(" -n, --name ") +  color.GreyP("(optionel)")+" : Name of the file in the bucket, if you don't specify it, the name of the file will be the same as the file copied\n"
	
// define the flags
var CopyObjectFlags = []cli.Flag{
	cli.StringFlag{
		Name: "destination, d",
		Usage: "Destination bucket",
	},
	cli.StringFlag{
		Name: "name, n",
		Usage: "Name of the file in the bucket",
	},
}

var CopyObjectCMD = cli.Command{
		Name:    "copy, cp",
		Category: "Object",
		Aliases: []string{"cp", "copy"},
		Usage:   "Copy a file from a path and insert it in a bucket destination",
		UsageText: registryUsageText,
		Description: createBucketDesc,
		OnUsageError: utils.OnUsageError,
		Flags: CopyObjectFlags,
		Action: copyObjectCmd,
		CustomHelpTemplate: utils.HelpTemplate,
}

func copyObjectCmd(c *cli.Context) error {
	path := c.Args().First()
	destination := c.String("destination")
	fileName := c.String("name")
	if(fileName == ""){
		fileName = extractFileName(path)
	}

	if errPath := verifyPath(path); errPath != nil {
		return errPath
	}
	_, err := global.Connexion.S3Client.HeadBucket(context.TODO(),&s3.HeadBucketInput{
		Bucket: aws.String(destination),
	})
	if err != nil {
		return errors.New(color.RedP("Bucket " + destination + " not found"))
	}

	err = UploadFile(destination, fileName, path)
	if err != nil {
		return errors.New(color.RedP("Error while uploading file"))
	}
	
	println(color.GreenP("File " + path + " copied to " + destination + " with the name " + fileName))
	return nil
}

func verifyPath(path string) error {
	if path == "" {
		return errors.New(color.RedP("You must specify a path"))
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.New(color.RedP("The path specified does not exist"))
	}
	return nil
}

func extractFileName(path string) string {
	isBackSlash := strings.Contains(path, "\\")
	splitedChar := "/"
	if isBackSlash {
		splitedChar = "\\"
	}
	pathSplitted := strings.Split(path, splitedChar)
	return pathSplitted[len(pathSplitted)-1]
}

// UploadFile reads from a file and puts the data into an object in a bucket.
func UploadFile(bucketName string, objectKey string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Couldn't open file %v to upload. Here's why: %v\n", fileName, err)
	} else {
		defer file.Close()
		_, err = global.Connexion.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			Body:   file,
		})
		if err != nil {
			log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
				fileName, bucketName, objectKey, err)
		}
	}
	return err
}

