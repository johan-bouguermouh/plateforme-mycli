package cmd

import (
	global "bucketool/cmd/global"
	env "bucketool/environment"
	utils "bucketool/utils"
	color "bucketool/utils/colorPrint"
	"context"
	"io"
	"net/http"
	"strings"

	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/urfave/cli"
)

// Description hiw to use the command
var downloadObjectUsageText string = `This command download a file from a bucket and insert it in a path.
	The first argument is the path where the file will be copied.
	The flag -b, --bucket is the bucket where the file is.
	The flag -n is the name of the file in the bucket.
	The flag -rename is the name of the file in the path (optionel), you don't need to specify extention, it will be added automatically. If you don't specify it, the name of the file will be the same as the file copied, but the could be changed if the MIME type is different.
	The flag -alias can be used to specify the alias to use, if you have specified the current alias, you can omit this flag.
	Usage of flag -alias is optional. If you use it, this flags must be placed  before command, like this : -alias <alias> donwload <path> -b <bucketName> -n <ObjectName>\n`+
	"Example : dl /path/to/file -d mybucket -n myfile.txt" + "\n" + "Example : -alias myalias dl /path/to/file -d mybucket -n myfile.txt" + "\n"

// Description of the flags
var downloadObjectDesc string = color.ColorPrint("Black",
	" Flags Descriptions ", &color.Options{
		Underline: true,
		Bold:      true,
		Background: "White",
	}) + "\n" +
	color.BlueP(" -b, --bucket ") +  color.GreyP("(Required)")+" : Bucket where the file is\n" +
	color.BlueP(" -n, --name ") +  color.GreyP("(optionel)")+" : Name of the file in the bucket, if you don't specify it, the name of the file will be the same as the file copied\n"+
	color.BlueP(" -rename, --rn ") +  color.GreyP("(optionel)")+" : Name of the file in the path, you don't need to specify extention, it will be added automatically. If you don't specify it, the name of the file will be the same as the file copied, but the could be changed if the MIME type is different.\n"

// define the flags
var DownloadObjectFlags = []cli.Flag{
	cli.StringFlag{
		Name: "bucket, b",
		Usage: "Bucket where the file is",
		Required: true,
	},
	cli.StringFlag{
		Name: "name, n",
		Usage: "Name of the file in the bucket",
		Required: true,
	},
	cli.StringFlag{
		Name: "rename, rn",
		Usage: "Name of the file in the path",
		Required: false,
	},
}

var DownloadObjectCMD = cli.Command{
		Name:    "download, dl",
		Category: "Object",
		Aliases: []string{"dl"},
		Usage:   "Download a file from a bucket and insert it in a path",
		UsageText: downloadObjectUsageText,
		Description: downloadObjectDesc,
		OnUsageError: utils.OnUsageError,
		Flags: DownloadObjectFlags,
		Action: downloadObjectCmd,
		CustomHelpTemplate: utils.HelpTemplate,
}

// downloadObjectCmd is the function that will be called when the command is executed
func downloadObjectCmd(c *cli.Context) error {
	path := c.Args().First()
	bucketName := c.String("bucket")
	fileName := c.String("name")
	rename := c.String("rename")
	if(rename == ""){
		rename = fileName
	}
	if(fileName == ""){
		return errors.New(color.RedP("You must specify the name of the file in the bucket"))
	}
	if(bucketName == ""){
		return errors.New(color.RedP("You must specify the bucket"))
	}
	if(path == ""){
		return errors.New(color.RedP("You must specify the path"))
	}

	if errPath := verifyPath(path); errPath != nil {
		return errPath
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.New(color.RedP("Path " + path + " not found"))
	}

	if !global.BucketExists(bucketName){
		return errors.New(color.RedP("Bucket " + bucketName + " not found"))
	}

	err := DownloadFile(bucketName, fileName, path, rename)
	if err != nil {
		return errors.New(color.RedP("Error while downloading file"))
	}

	println(color.GreenP("File " + fileName + " downloaded from " + bucketName + " and copied to " + path))
	return nil
}


// DownloadFile gets an object from a bucket and stores it in a local file.
func DownloadFile(bucketName string, objectKey string, fileName string, rename string) error {
	// Obtenir l'objet
	result, err := global.Connexion.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		println(color.RedP("Couldn't get object " + objectKey + ". Here's why:"))
		return err
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		println(color.RedP("Couldn't read object body from " + objectKey + ". Here's why:"))
		return err
	}

	rename = setContentType(body, rename)
	
	file, err := os.Create(fileName+"/"+rename)
	if err != nil {
		println(color.RedP("Couldn't create file "+ fileName+". Here's why:"),err)
		return err
	}
	defer file.Close()

	_, err = file.Write(body)
	return err
}

func setContentType(body []byte, rename string) string {
    // Détecter le type MIME du contenu
    contentType := http.DetectContentType(body)
    println("Detected content type: ", contentType)

    // Séparer le type MIME pour obtenir l'extension
    categorie, extention, canCutct := strings.Cut(contentType, "/")
    if canCutct {
		if(contentType == "text/plain; charset=utf-8"){
			extention = "txt"
		}
        if env.IsDebugMode {
            println(color.GreyP("Read categorie of MIME type : " + categorie))
            println(color.GreyP("Read extention of MIME type : " + extention))
        }
        if rename != "" {
            nameTosaved, oldExt, canCut := strings.Cut(rename, ".")
            if oldExt != "" && oldExt != extention {
                if env.IsDebugMode {
                    println(color.GreyP("Removed old extention : " + oldExt))
                    println(color.GreyP("Added new extention : " + extention))
                }
                return nameTosaved + "." + extention
            } else if !canCut {
                if env.IsDebugMode {
                    println(color.GreyP("No extention found in the name : " + rename))
                    println(color.GreyP("Added new extention : " + extention))
                }
                return rename + "." + extention
            } else {
				if env.IsDebugMode {
					println(color.GreyP("Extention is the same as the MIME type, no need to change"))
				}
				return rename
			}
		} else {
			if env.IsDebugMode {
				println(color.RedP("No name specified, added extention : " + extention))
			}
			return "file." + extention
		}
    } else {
        if env.IsDebugMode {
            println(color.RedP("No MIME type detected"))
        }
    }
    return rename
}
