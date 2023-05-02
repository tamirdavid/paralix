# Paralix CLI

Paralix CLI is a command line interface tool written in Golang. The purpose of this CLI is to provide a faster and better way to execute bash commands in parallel. 

## Command

Currently, the CLI has one command called `command`. This command allows you to execute a bash command with placeholders in parallel. 

### Flags

Here are the flags available for the `command` command: 


- `--execute`, `-e`: A string flag that takes the command to execute with placeholders. The placeholders are denoted by `<KEY>` and will be replaced with values provided either by the `-p` flag or an input file.
<br>Example: `--execute 'echo <WHAT_SHOULD_ECHO>'`.

- `--placeholder`, `-p`: A string flag that takes placeholders in the format of `KEY={VALUE1,VALUE2,VALUE3}`. Multiple placeholders can be separated by a comma. These values will replace the placeholders in the command provided by the `-e` flag. Example: `-p 'WHAT_SHOULD_ECHO={HELLO,WORLD}'`.

- `--inputfile`, `-f`: A string flag that takes a file path that contains the inputs to run. Each input should be on a new line. These inputs will replace the placeholders in the command provided by the `-e` flag.<br> Example: `-f 'WHAT_SHOULD_ECHO'`.

- `--output`, `-o`: A string flag that takes a file path to write the output of the command. <br>
The output will be written in the following format: <br>
PlaceholderA<br>
Output of command with PlaceholderA<br><br>
PlaceholderB<br>
Output of command with PlaceholderB<br><br>
PlaceholerN
Output of command with PlaceholderN<br><br>

### Examples

Here are some examples of how to use the Paralix CLI:


#### Example 1:

Execute a command with placeholders and values provided by an input file.<br>
`$ paralix command --execute 'aws s3 ls s3://<BUCKET_NAME> --inputfile BUCKET_NAME -o upgrade_output`
<br>_Assuming the `BUCKET_NAME` file contains the following_:
myfiles1
myfiles2

This will execute the command `aws s3 ls s3://<BUCKET_NAME>` twice, once with `myfiles1` and once with `myfiles2`. <br>
- _The output file will looks like_:<br>
**bucket1**<br>
file1<br>
file2<br>
file3<br>
**bucket2**<br>
file4<br>
file5<br>
file6<br>

#### Example 2:

Execute a command with placeholders and write the output to a file.<br>
`$paralix command --execute 'curl https://some.api/endpoint/to/send/with/<NAME>' --inputfile 'NAME' --output 'curl_output.txt'`<br>
- _Assuming the NAME file contains:_ <br>
nameA<br>
nameB<br>

- _The output file will looks like:_<br>
**NameA**<br>
curl output for NameA<br>
**NameB**<br>
curl output for NameB<br>



## Installation
To use Paralix CLI, you need to have Go installed on your system. If you don't have Go installed, you can download it from the official Go website.

Once you have Go installed, you can install Paralix CLI using the following command:<br>
`go get github.com/tamirdavid/paralix`


## Contributing
If you would like to contribute to the Paralix CLI, please open a pull request on the GitHub repository.