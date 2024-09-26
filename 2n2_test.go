package main

import (
	color "bucketool/utils/colorPrint"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)
var Path string

func init() {
	println(
		color.ColorPrint("Black", "    INIT TEST N2N    ",&color.Options{
			Background: "Yellow",
			Bold: true,
		}))
	// On récupère le path du fichier de test
	path, _ := os.Getwd()
	Path = path
	println(path)
}

func TestCLI(t *testing.T) {
    tests := []struct {
        name    string
        command string
        args    []string
        want    string
		print string
    }{
		{ 
			name: "Test alias command Set",
			command: "bucketool",
			args: []string{"alias", "set","minio","-H", "http://localhost", "-p", "9000", "-k", "minioadmin", "-s", "minioadmin"},
			want: "\x1b[3m\x1b[90mAlias has minio been saved\x1b[0m\nRegistered Alias on Name:  \x1b[32mminio\x1b[0m \x1b[3m\x1b[90mhttp://localhost:9000\x1b[0m\n",
			print: color.GreenP("🗸 Testing Alias Set"),
		},
		{
			name: "Test alias command List",
			command: "bucketool",
			args: []string{"alias", "list", "-d"},
			want: " \x1b[32mminio\x1b[0m \x1b[90m(http://localhost:9000)\x1b[0m\n",
			print: color.GreenP("🗸 Testing Alias List"),
		},
		{
			name : "Test alias command Current switch",
			command: "bucketool",
			args: []string{"alias", "current", "-s","minio"},
			want: "\x1b[3m\x1b[90mSwitch Alias to \x1b[32mminio\x1b[0m\x1b[0m\n",
			print: color.GreenP("🗸 Testing Current switch Set"),
		},
		{
			name: "Create Bucket",
			command: "bucketool",
			args: []string{"bucket", "create", "mybuckettest"},
			want: "mybuckettest\n\x1b[32mBucket mybuckettest created successfully\x1b[0m\n",
			print: color.GreenP("🗸 Testing Bucket Create"),
		},
		{
			name: "List Bucket",
			command: "bucketool",
			args: []string{"bucket", "ls"},
			want: "\x1b[34mmybuckettest\x1b[0m\n",
			print: color.GreenP("🗸 Testing Bucket List"),
		},
		{
			name: "Copy Object to Bucket",
			command: "bucketool",
			args: []string{"cp", Path+"\\dataTest\\valid_file.txt", "-d", "mybuckettest","-n","filetest.txt"},
			want: "valid_file.txt copied to mybuckettest with the name filetest.txt\x1b[0m\n",
			print: color.GreenP("🗸 Testing Copy Object to Bucket"),
		},
		{
			name: "List Object in Bucket",
			command: "bucketool",
			args: []string{"ls", "-b", "mybuckettest", "-d"},
			want: "\x1b[34mfiletest.txt\x1b[0m\n",
			print: color.GreenP("🗸 Testing List Object in Bucket"),
		},
		{
			name: "Download Object from Bucket",
			command: "bucketool",
			args: []string{"dl", Path+"\\dataTest\\", "-b", "mybuckettest", "-n", "filetest.txt", "-rn", "filetestdownload.txt"},
			want: "Detected content type:  text/plain; charset=utf-8\n\x1b[32mFile filetest.txt downloaded from mybuckettest and copied to",
			print: color.GreenP("🗸 Testing Download Object from Bucket"),
		},

        // {
        //     name:    "Test copy command with original name",
        //     command: "bucketool",
        //     args:    []string{"cp", "/path/to/file", "-d", "mybucket"},
        //     want:    "File copied successfully",
        // },
        // {
        //     name:    "Test copy command with specific name",
        //     command: "bucketool",
        //     args:    []string{"copy", "/path/to/file", "-d", "mybucket", "-n", "myfile.txt"},
        //     want:    "File copied successfully",
        // },
        // {
        //     name:    "Test download command",
        //     command: "bucketool",
        //     args:    []string{"download", "/path/to/file", "-b", "mybucket", "-n", "myfile.txt"},
        //     want:    "File downloaded successfully",
        // },
        // {
        //     name:    "Test download command with rename",
        //     command: "bucketool",
        //     args:    []string{"download", "/path/to/file", "-b", "mybucket", "-n", "myfile.txt", "-rename", "newfile"},
        //     want:    "File downloaded successfully",
        // },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cmd := exec.Command(tt.command, tt.args...)
            output, err := cmd.CombinedOutput()
			if err != nil {
                // Log the error but continue to check the output
                t.Logf("Command failed with error: %v", err)
            }
            assert.NoError(t, err)
            assert.Contains(t, string(output), tt.want)
			t.Log(tt.print)
        })
    }
}