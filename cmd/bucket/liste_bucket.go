package bucket

import (
	utils "bucketool/utils"
	color "bucketool/utils/colorPrint"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/urfave/cli"
)

// Description how to use the command
var listBucketUsageText string = color.GreyP("Command Exemple : bucketool bucket list \n") +
	`This command will list all the bucket.
	The flag -d or --detail can be used to list all Buckets registred with details url.
	Usage of flag -alias is optional. If you use it, this flags must be placed after "bucket" and before " list", like this : bucket -alias <alias> list \n`+
	color.GreyP("Example : bucket list") + "\n" + color.GreyP("Example : bucket -alias myalias ls") + "\n"


// Description of the flags
var listBucketDesc string = color.ColorPrint("Black",
	" Flags Descriptions ", &color.Options{
		Underline: true,
		Bold:      true,
		Background: "White",
	}) + "\n" +
	color.BlueP(" -d, --detail ") +  color.GreyP("(Optionnal)")+" : List all Buckets registred with details url\n"



var ListBucketFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "details, d",
		Usage: "Show details of the bucket",
	},
}

// ListBucketCMD create the command to list a bucket
var ListBucketCMD = cli.Command{
		Name:    "list, ls",
		Aliases: []string{"ls"},
		Category: "Bucket",
		Usage:   "list all the bucket",
		Description: listBucketDesc,
		UsageText: listBucketUsageText,
		Flags:   ListBucketFlags,
		Before: BeforeUseAlias,
		Action:  listBucketCmd,
		OnUsageError: utils.OnUsageError,
		CustomHelpTemplate : utils.HelpTemplate,
}

func listBucketCmd(c *cli.Context) error {
	showDetails := c.Bool("details")

	result, err := Connexion.S3Service.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return err
	}

	for _, bucket := range result.Buckets {
		println(color.BlueP(*bucket.Name))
		if showDetails {
            printBucketDetails(*bucket.Name)
        }
	}

	return nil
}

func printBucketDetails(bucketName string) {
    // Get bucket location
	indentation := "  "
    location, err := Connexion.S3Service.GetBucketLocation(&s3.GetBucketLocationInput{
        Bucket: aws.String(bucketName),
    })
    if err != nil {
        fmt.Println(color.RedP("Failed to get bucket location:"), err)
        return
    }
	locationStr := aws.StringValue(location.LocationConstraint)
	if locationStr == "" {
		locationStr = "us-east-1"
	}
    fmt.Println(indentation,color.GreenP("Location:"), color.GreyP(locationStr))

    // Get bucket ACL
    acl, err := Connexion.S3Service.GetBucketAcl(&s3.GetBucketAclInput{
        Bucket: aws.String(bucketName),
    })
    if err != nil {
        fmt.Println(color.RedP("Failed to get bucket ACL:"), err)
        return
    }
    fmt.Println(indentation, color.GreenP("ACL:"))
    for _, grant := range acl.Grants {
        fmt.Printf("    Grantee: %s\n    Permission: %s\n", color.GreyP(aws.StringValue(grant.Grantee.DisplayName)), color.GreyP(aws.StringValue(grant.Permission)))
    }

    // Get bucket logging
    logging, err := Connexion.S3Service.GetBucketLogging(&s3.GetBucketLoggingInput{
        Bucket: aws.String(bucketName),
    })
    if err != nil {
        fmt.Println(color.RedP("Failed to get bucket logging:"), err)
        return
    }
    if logging.LoggingEnabled != nil {
        fmt.Println(indentation,color.GreenP("Logging: Enabled"))
        fmt.Printf("    TargetBucket: %s\n    TargetPrefix: %s\n", aws.StringValue(logging.LoggingEnabled.TargetBucket), aws.StringValue(logging.LoggingEnabled.TargetPrefix))
    } else {
        fmt.Println(indentation,color.GreenP("Logging:"), color.GreyP("Disabled"))
    }

    // Get bucket versioning
    versioning, err := Connexion.S3Service.GetBucketVersioning(&s3.GetBucketVersioningInput{
        Bucket: aws.String(bucketName),
    })
    if err != nil {
        fmt.Println(color.RedP("Failed to get bucket versioning:"), err)
        return
    }

	strVersionningStatus := aws.StringValue(versioning.Status)
	if strVersionningStatus == "" {
		strVersionningStatus = "Disabled"
	}
    fmt.Println(indentation, color.GreenP("Versioning:"), color.GreyP(strVersionningStatus))
}