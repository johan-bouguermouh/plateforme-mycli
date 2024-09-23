package bucket

import (
	utils "bucketool/utils"
	color "bucketool/utils/colorPrint"
	"errors"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/urfave/cli"
)

// Description how to use the command
var createBucketUsageText string = color.GreyP("Command Exemple : bucketool bucket create <name> -p <port> -h <host> -k <keyname> -s <Secret> \n") +
	`This command will create a new bucket.
	The first argument is the name of the bucket.
	This argument is required, and must be unique. Argument must be lowercase, and contain only lowercase letters, numbers, hyphens (-), and periods (.).
	th flag -alias can be used to specify the alias to use, if you have specified the current alias, you can omit this flag.
	Usage of flag -alias is optional. If you use it, this flags must be placed after "bucket" and before " create <name>", like this : bucket -alias <alias> create <name> \n`+
	color.GreyP("Example : bucket create mybucket") + "\n" + color.GreyP("Example : bucket -alias myalias create mybucket") + "\n"
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
		Category: "Bucket",
		Usage:   "create a bucket",
		Description: createBucketDesc,
		UsageText: createBucketUsageText,
		Flags:   CreateBucketFlags,
		Aliases: []string{"c"},
		Before: BeforeUseAlias,
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

	_, err := Connexion.S3Service.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		println(color.RedP("Failed to create bucket"))
		//if 409 error code is returned, it means the bucket already exists
		if strings.Contains(err.Error(), "409") {
			println(color.RedP("Bucket "+ bucketName +" already exists"))
		}
		return err
	}

	err = Connexion.S3Service.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
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
