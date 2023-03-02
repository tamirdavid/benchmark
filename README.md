# benchmark
A [JFrog CLI plugin](https://www.jfrog.com/confluence/display/CLI/JFrog+CLI#JFrogCLI-JFrogCLIPlugins) to benchmark Artifactory upload or download speeds

## About this plugin
This plugin is for measuring upload and download speeds to and from an Artifactory instance.

It should be used to benchmarks your instance uploads and downloads speeds.<br>
The plugin provides metrics on upload or download speeds, as well as the number of successful and failed operations. The results are written to a `csv` formatted file.

This plugin can be used to compare your Artifactory instance behaviour before and after a change like a version upgrade or a configuration change.



## CLI Configuration

### Configuring Artifactory credentials using jf c command
1. Add a new server using,

    ```
    jf c add
    ```
2. Add the Artifactory URL and authentication details.

### Configure Artifactory credentials without using jf c
1. Documentation will soon explain how to utilize the url, password, and username flags to send credentials to the benchmark plugin.

## Installation with JFrog CLI
Installing the latest version:

`$ jf plugin install benchmark`

Installing a specific version:

`$ jf plugin install benchmark@version`

Uninstalling a plugin

`$ jf plugin uninstall benchmark`

## Building from source
To build the **benchmark** binary
```shell
go build -o benchmark
```
To build the **benchmark** binary for multiple operating systems and architectures (Mac, Linux and Windows)
```shell
./build-binary.sh
```

## Usage
### Commands
* up
    - Flags:
        - size [Optional] -  Determine the size of the files (in MB) that will be generated for testing the upload process. **[Default: 50]**
        - iterations [Optional] - How many files will be created for testing the upload process. **[Default: 30]**
        - repo_name [Optional] - Repository the tests will be executed on. **[Default: benchmark-up-tests]** <br> <br>
        - url [Optional] - If using custom server (not already configured one) **[No default value]**
        - username [Optional] - **[No default value]**
        - password [Optional] - **[No default value]**
        - append [Optional] - Append the csv results to existing file **[No default value]**
        - same_file [Optional] - benchmark will upload the same file instead of generating and uploading multiple files
    - Example:
    ```
  $ jf benchmark up
  $ jf benchmark up --size 50 --iterations 5
  $ jf benchmark up --size 50 --iterations 5 --repo_name mytestrepo
  $ jf benchmark up --size 50 --iterations 5 --repo_name mytestrepo --append benchmark-upload-2023-02-21T11:30:29.csv
  $ jf benchmark up --url <myserverurl> --username <username> --password <password>
  $ jf benchmark up --url <myserverurl> --username <username> --password <password> --iterations 2 --size 150
  ```
* dl
    - Flags:
        - size [Optional] - Determine the size of the files (in MB) that will be generated for testing the download process. **[Default: 50]**
        - iterations [Optional] - How many files will be created for testing the download process. **[Default: 30]**
        - repo_name [Optional] - Repository the tests will be executed on. **[Default: benchmark-up-tests]** <br> <br>
        - url [Optional] - If using custom server (not already configured one) **[No default value]**
        - username [Optional] - **[No default value]**
        - password [Optional] - **[No default value]**
        - append [Optional] - Append the csv results to existing file **[No default value]**
        - same_file [Optional] - benchmark will download the same file instead of generating and uploading multiple files
    - Example:
    ```
  $ jf benchmark dl  
  $ jf benchmark dl --size 50 --iterations 5
  $ jf benchmark dl --size 50 --iterations 5 --repo_name mytestrepo
  $ jf benchmark dl --size 50 --iterations 5 --append benchmark-download-2023-02-21T11:30:29.csv
  $ jf benchmark dl --url <myserverurl> --username <username> --password <password>
  $ jf benchmark dl --url <myserverurl> --username <username> --password <password> --iterations 15 --size 73
  ```

### Output file Example
* Both the 'dl' and 'up' commands produce CSV files that contain the filename, size, and the elapsed time for uploading/downloading:
```
file,size (MB),time taken (sec),speed (MB/sec)
/tmp/testfiles/File1.txt,50,14.664069103s,3.41
/tmp/testfiles/File2.txt,50,15.302840585s,3.27
/tmp/testfiles/File3.txt,50,17.859408288s,2.80
/tmp/testfiles/File4.txt,50,14.003771159s,3.57
/tmp/testfiles/File5.txt,50,14.498499084s,3.45
/tmp/testfiles/File6.txt,50,14.844652081s,3.37
/tmp/testfiles/File7.txt,50,14.286648744s,3.50
/tmp/testfiles/File8.txt,50,14.35213191s,3.48
/tmp/testfiles/File9.txt,50,14.667853445s,3.41
/tmp/testfiles/File10.txt,50,14.682988659s,3.41
```


## Release Notes
The release notes are available [here](RELEASE.md).