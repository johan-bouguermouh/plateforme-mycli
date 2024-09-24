package bucket

import (
	global "bucketool/cmd/global"
	utils "bucketool/utils"
	color "bucketool/utils/colorPrint"
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/urfave/cli"
)

// Description how to use the command
var deleteBucketUsageText string = color.GreyP("Command Exemple : bucketool bucket delete <name> \n") +
	`This command will delete a bucket.
	The first argument is the name of the bucket.
	This argument is required, and must be existing.const
	Usage of flag -alias is optional. If you use it, this flags must be placed after "bucket" and before " delete <name>", like this : bucket -alias <alias> delete <name> \n`+
	color.GreyP("Example : bucket delete mybucket") + "\n" + color.GreyP("Example : bucket -alias myalias delete mybucket") + "\n"

// Description of the flags
var deleteBucketDesc string = color.ColorPrint("Black",
	" Flags Descriptions ", &color.Options{
		Underline: true,
		Bold:      true,
		Background: "White",
	}) + "\n" + color.GreyP("No flags available") + "\n"

var DeleteBucketFlags = []cli.Flag{}

// DeleteBucketCMD create the command to delete a bucket
var DeleteBucketCMD = cli.Command{
	Name:    "delete",
	Aliases: []string{"d"},
	Category: "Bucket",
	Usage:   "delete a bucket",
	Description: deleteBucketDesc,
	UsageText: deleteBucketUsageText,
	Flags:   DeleteBucketFlags,
	Action:  deleteBucketCmd,
	OnUsageError: utils.OnUsageError,
	CustomHelpTemplate : utils.HelpTemplate,
}

func deleteBucketCmd(c *cli.Context) error {
	bucketName := c.Args().First()

	validationsProblems := valideNameBucket(bucketName)
	if validationsProblems != "" {
		return errors.New(color.RedP(validationsProblems))
	}

	//On verifi qu'un bucket existe bien avec ce nom
	_, err := global.Connexion.S3Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return errors.New(color.RedP("Bucket " + bucketName + " not found"))
	}

	_, err = global.Connexion.S3Client.DeleteBucket(context.TODO(),&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return err
	}

	fmt.Println(color.GreenP("Bucket " + bucketName + " deleted"))
	return nil
}
