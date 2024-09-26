package cmd

import (
	"bucketool/cmd/global"
	utils "bucketool/utils"
	color "bucketool/utils/colorPrint"
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/urfave/cli"
)

// Description hiw to use the command
var listBuucketObjectsUsageText string = `This command list all the objects in a bucket.
	The flag -b, -bucket is the destination bucket.
	The flag -d, -details is optional, it shows details of the objects.
	The flag -alias can be used to specify the alias to use, if you have specified the current alias, you can omit this flag.
	Usage of flag -alias is optional. If you use it, this flags must be placed before name of command, like this : -alias <alias> list \n`+
	"Example : list_objects -b mybucket" + "\n" + "Example : -alias myalias list -b mybucket" + "\n"

// Description of the flags
var listBucketObjectsDesc string = color.ColorPrint("Black",
	" Flags Descriptions ", &color.Options{
		Underline: true,
		Bold:      true,
		Background: "White",
	}) + "\n" +
	color.BlueP(" -b, --bucket ") + color.GreyP("(Required)") + " : Destination bucket\n"+
	color.BlueP(" -d, --details ") + color.GreyP("(optionel)") + " : Show details of the objects\n"

// define the flags
var ListBucketObjectsFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "bucket, b",
		Usage: "Destination bucket",
	},
	cli.BoolFlag{
		Name: "details, d",
		Usage: "Show details of the objects",
	},
}

var ListBucketObjectsCMD = cli.Command{
	Name:    "list, ls",
	Category: "Object",
	Aliases: []string{"ls", "list"},
	Usage:   "List all the objects in a bucket",
	UsageText: listBuucketObjectsUsageText,
	Description: listBucketObjectsDesc,
	OnUsageError: utils.OnUsageError,
	Flags: ListBucketObjectsFlags,
	Action: listBucketObjectsCmd,
	CustomHelpTemplate: utils.HelpTemplate,
}

func listBucketObjectsCmd(c *cli.Context) error {
	bucketName := c.String("bucket")
	hasDetails := c.Bool("details")
	if bucketName == "" {
		return errors.New(color.RedP("You must specify a bucket name"))
	}

	// Verify if the bucket exists
	if !global.BucketExists(bucketName) {
		return errors.New(color.RedP("The bucket " + bucketName + " doesn't exist"))
	}

	objects, err := ListObjects(bucketName)
	if err != nil {
		return err
	}

	if len(objects) == 0 {
		println(color.GreyP("No objects found in bucket"+ bucketName +"\n"))
	} else {
		println(color.ColorPrint("Grey","Objects in bucket "+ color.GreenP(bucketName), &color.Options{
			Bold: true,
			}))
		for _, object := range objects {
			println("  - "+color.BlueP(*object.Key))
			if(hasDetails){
				println("    Last Modified: "+color.GreyP(object.LastModified.String()))
				println("    Size: " + color.GreyP(strconv.FormatInt(*object.Size, 10) + " bytes"))
				println("    Storage Class: "+color.GreyP(string(object.StorageClass)))
				println("    ETag: "+color.GreyP(*object.ETag))
			}
		}
	}

	return nil
}

// ListObjects lists the objects in a bucket.
func ListObjects(bucketName string) ([]s3Types.Object, error) {
	result, err := global.Connexion.S3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	var contents []s3Types.Object
	if err != nil {
		log.Printf("Couldn't list objects in bucket %v. Here's why: %v\n", bucketName, err)
	} else {
		contents = result.Contents
	}
	return contents, err
}
