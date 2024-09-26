package cmd

import (
	global "bucketool/cmd/global"
	utils "bucketool/utils"
	color "bucketool/utils/colorPrint"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/urfave/cli"
)

// Description hiw to use the command
var deleteObjectUsageText string = `This command delete a file from a bucket.
	The flag -b is the destination bucket.
	The flag -n is the name of the file in the bucket.
	The flag -alias can be used to specify the alias to use, if you have specified the current alias, you can omit this flag.
	Usage of flag -alias is optional. If you use it, this flags must be before command name, like this : -alias <alias> delete -d <bucketName> -n <ObjectName>\n`+
	"Example : delete -b mybucket -n myfile.txt" + "\n" + "Example : -alias myalias delete -b mybucket -n myfile.txt" + "\n"

// Description of the flags
var deleteObjectDesc string = color.ColorPrint("Black",
	" Flags Descriptions ", &color.Options{
		Underline: true,
		Bold:      true,
		Background: "White",
	}) + "\n" +
	color.BlueP(" -b, --bucket ") +  color.GreyP("(Required)")+" : Destination bucket\n" +
	color.BlueP(" -n, --name ") +  color.GreyP("(Required)")+" : Name of the file in the bucket\n"

// define the flags
var DeleteObjectFlags = []cli.Flag{
	cli.StringFlag{
		Name: "bucket, b",
		Usage: "Destination bucket",
		Required: true,
	},
	cli.StringFlag{
		Name: "name, n",
		Usage: "Name of the file in the bucket",
		Required: true,
	},
}

var DeleteObjectCMD = cli.Command{
		Name:    "delete, del",
		Category: "Object",
		Aliases: []string{"del"},
		Usage:   "Delete a file from a bucket",
		UsageText: deleteObjectUsageText,
		Description: deleteObjectDesc,
		OnUsageError: utils.OnUsageError,
		Flags: DeleteObjectFlags,
		Action: deleteObjectCmd,
		CustomHelpTemplate: utils.HelpTemplate,
}

func deleteObjectCmd(c *cli.Context) error {

	// get the destination bucket
	destination := c.String("bucket")

	// get the name of the file in the bucket
	name := c.String("name")


	//On verfie que le bucket existe bien
	_, err := global.Connexion.S3Client.HeadBucket(context.TODO(),&s3.HeadBucketInput{
		Bucket: aws.String(destination),
	})
	if err != nil {
		println(color.RedP("The bucket " + destination + " does not exist"))
		return err
	}

	//On verifie que l'objet existe bien
	_, err = global.Connexion.S3Client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(destination),
		Key:    aws.String(name),
	})
	if err != nil {
		println(color.RedP("The object " + name + " does not exist in the bucket " + destination))
		return err
	}

	// create the input for the request
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(destination),
		Key:    aws.String(name),
	}

	// send the request
	_, errDelete := global.Connexion.S3Client.DeleteObject(context.TODO(), input)
	if errDelete != nil {
		println(color.RedP("Error while deleting the object " + name + " from the bucket " + destination))
		return errDelete
	}

	println(color.GreenP("The object " + name + " has been deleted from the bucket " + destination))

	return nil
}