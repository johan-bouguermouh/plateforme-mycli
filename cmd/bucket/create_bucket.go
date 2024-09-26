package bucket

import (
	global "bucketool/cmd/global"
	conn "bucketool/connexion"
	utils "bucketool/utils"
	color "bucketool/utils/colorPrint"
	"context"
	"errors"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/urfave/cli"
)

// Description how to use the command
var createBucketUsageText string = color.GreyP("Command Exemple : bucketool bucket create <name> -p <port> -h <host> -k <keyname> -s <Secret> \n") +
	`This command will create a new bucket.
	The first argument is the name of the bucket.
	This argument is required, and must be unique. Argument must be lowercase, and contain only lowercase letters, numbers, hyphens (-), and periods (.).
	The flag -alias can be used to specify the alias to use, if you have specified the current alias, you can omit this flag.
	Usage of flag -alias is optional. If you use it, this flags must be placed before "bucket", like this : -alias <alias> bucket create <name> \n`+
	color.GreyP("Example : bucket create mybucket") + "\n" + color.GreyP("Example : bucketool -alias myalias bucket create mybucket") + "\n"
// Description of the flags
var createBucketDesc string = color.ColorPrint("Black",
	" Flags Descriptions ", &color.Options{
		Underline: true,
		Bold:      true,
		Background: "White",
	}) + "\n" + color.GreyP("No flags available") + "\n"
	
var CreateBucketFlags = []cli.Flag{}

// CreateBucketCmd create the command to create a bucket
var CreateBucketCMD = cli.Command{
		Name:    "create",
		Aliases: []string{"c"},
		Category: "Bucket",
		Usage:   "create a bucket",
		Description: createBucketDesc,
		UsageText: createBucketUsageText,
		Flags:   CreateBucketFlags,
		Action:  createBucketCmd,
		OnUsageError: utils.OnUsageError,
		CustomHelpTemplate : utils.HelpTemplate,
}

func createBucketCmd(c *cli.Context) error {
	bucketName := c.Args().First()

	validationsProblems := valideNameBucket(bucketName)
	if validationsProblems != "" {
		return errors.New(color.RedP(validationsProblems))
	}
	   

	_, err := global.Connexion.S3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(conn.Region),
		},
	})
	if err != nil {
		println(color.RedP("Failed to create bucket"))
		//if 409 error code is returned, it means the bucket already exists
		if strings.Contains(err.Error(), "409") {
			println(color.RedP("Bucket "+ bucketName +" already exists"))
		}
		log.Printf("Couldn't create bucket %v in Region %v. Here's why: %v\n",
		bucketName, conn.Region, err)
		return err
	}

	err = waitForBucketExists(bucketName)
	if err != nil {
		println(color.RedP("Failed to wait for bucket creation :"))
		return err
	}

	println(color.GreenP("Bucket "+ bucketName +" created successfully"))


	return nil
}

var RegexBucketName = regexp.MustCompile(`^[a-z0-9.-]+$`)

func valideNameBucket (name string) string {

	// Check if the name is empty
	if name == "" {
		return "Bucket name is required"
	}

	//Check if the name is to short
	if len(name) < 3 {
		return "Bucket name must be at least 3 characters long"
	}

	// Check if the name is to long
	if len(name) > 63 {
		return "Bucket name must be at most 63 characters long"
	}

	// Check if bucket name has no uppercase letters
	if strings.ToLower(name) != name {
		return "Bucket name must be lowercase"
	}

	//Check if bucketName have valid characters
	if !RegexBucketName.MatchString(name) {
		return "Bucket name must contain only lowercase letters, numbers, hyphens (-), and periods (.)"
	}

	return ""
}

func waitForBucketExists(bucketName string) error {
    for {
        _, err := global.Connexion.S3Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
            Bucket: aws.String(bucketName),
        })
        if err == nil {
            return nil
        }

        // Check if the error is not a 404 (bucket not found)
        var notFound *types.NotFound
        if !errors.As(err, &notFound) {
            return err
        }

        // Wait for a short period before retrying
        time.Sleep(5 * time.Second)
    }
}
