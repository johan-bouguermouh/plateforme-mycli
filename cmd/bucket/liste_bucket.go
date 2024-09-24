package bucket

import (
	global "bucketool/cmd/global"
	utils "bucketool/utils"
	color "bucketool/utils/colorPrint"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

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
		Action:  listBucketCmd,
		OnUsageError: utils.OnUsageError,
		CustomHelpTemplate : utils.HelpTemplate,
}

func listBucketCmd(c *cli.Context) error {
	showDetails := c.Bool("details")

	result, err := global.Connexion.S3Client.ListBuckets(context.TODO(),&s3.ListBucketsInput{})
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
    indentation := "  "

    // Get bucket location
    location, err := global.Connexion.S3Client.GetBucketLocation(context.TODO(), &s3.GetBucketLocationInput{
        Bucket: aws.String(bucketName),
    })
    if err != nil {
        fmt.Println(color.RedP("Failed to get bucket location:"), err)
        return
    }
    locationStr := string(location.LocationConstraint)
    if locationStr == "" {
        locationStr = "us-east-1"
    }
    fmt.Println(indentation, color.GreenP("Location:"), color.GreyP(locationStr))

    // Get bucket ACL
    acl, err := global.Connexion.S3Client.GetBucketAcl(context.TODO(), &s3.GetBucketAclInput{
        Bucket: aws.String(bucketName),
    })
    if err != nil {
        fmt.Println(color.RedP("Failed to get bucket ACL:"), err)
        return
    }
    fmt.Println(indentation, color.GreenP("ACL:"))
    for _, grant := range acl.Grants {
        grantee := "Unknown"
        if grant.Grantee != nil && grant.Grantee.DisplayName != nil {
            grantee = *grant.Grantee.DisplayName
        }
        fmt.Printf("    Grantee: %s\n    Permission: %s\n", color.GreyP(grantee), color.GreyP(string(grant.Permission)))
    }

    // Get bucket logging
    logging, err := global.Connexion.S3Client.GetBucketLogging(context.TODO(), &s3.GetBucketLoggingInput{
        Bucket: aws.String(bucketName),
    })
    if err != nil {
        fmt.Println(color.RedP("Failed to get bucket logging:"), err)
        return
    }
    if logging.LoggingEnabled != nil {
        fmt.Println(indentation, color.GreenP("Logging: Enabled"))
        fmt.Printf("    TargetBucket: %s\n    TargetPrefix: %s\n", aws.ToString(logging.LoggingEnabled.TargetBucket), aws.ToString(logging.LoggingEnabled.TargetPrefix))
    } else {
        fmt.Println(indentation, color.GreenP("Logging:"), color.GreyP("Disabled"))
    }

    // Get bucket versioning
    versioning, err := global.Connexion.S3Client.GetBucketVersioning(context.TODO(), &s3.GetBucketVersioningInput{
        Bucket: aws.String(bucketName),
    })
    if err != nil {
        fmt.Println(color.RedP("Failed to get bucket versioning:"), err)
        return
    }

    strVersioningStatus := string(versioning.Status)
    if strVersioningStatus == "" {
        strVersioningStatus = "Disabled"
    }
    fmt.Println(indentation, color.GreenP("Versioning:"), color.GreyP(strVersioningStatus))
}