# benchmark
A [JFrog CLI plugin](https://www.jfrog.com/confluence/display/CLI/JFrog+CLI#JFrogCLI-JFrogCLIPlugins) to benchmark Artifactory upload or download speeds

## About this plugin
This plugin is for measuring upload and download speeds to and from an Artifactory instance.

It should be used to benchmarks your instance uploads and downloads speeds.<br>
The plugin provides metrics on upload or download speeds, as well as the number of successful and failed operations. The results are written to a `csv` formatted file.

This plugin can be used to compare your Artifactory instance behaviour before and after a change like a version upgrade or a configuration change.

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
        - size: Determine the size of the files (in MB) that will be generated for testing the upload process. **[Default: 50]** **[Mandatory: False]**
        - iterations: How many files will be created for testing the upload process. **[Default: 30]** **[Mandatory: False]**
        - repo_name: Repository the tests will be executed on. **[Default: benchmark-up-tests]** **[Mandatory: False]**
    - Example:
    ```
  $ jf benchmark up
  $ jf benchmark up --size 50 --iterations 5
  $ jf benchmark up --size 50 --iterations 5 --repo_name mytestrepo
  ```
* dl
    - Flags:
        - size: Determine the size of the files (in MB) that will be generated for testing the download process. **[Default: 50]** **[Mandatory: False]**
        - iterations: How many files will be created for testing the download process. **[Default: 30]** **[Mandatory: False]**
        - repo_name: Repository the tests will be executed on. **[Default: benchmark-up-tests]** **[Mandatory: False]**
    - Example:
    ```
  $ jf benchmark dl  
  $ jf benchmark dl --size 50 --iterations 5
  $ jf benchmark dl --size 50 --iterations 5 --repo_name mytestrepo
  ```

### Output file Example
* Both the 'dl' and 'up' commands produce CSV files that contain the filename, size, and the elapsed time for uploading/downloading:
```
file,size,time_taken,Speed
/tmp/testfiles/File1.txt,50MB,14.664069103s,3.41 MB/s
/tmp/testfiles/File2.txt,50MB,15.302840585s,3.27 MB/s
/tmp/testfiles/File3.txt,50MB,17.859408288s,2.80 MB/s
/tmp/testfiles/File4.txt,50MB,14.003771159s,3.57 MB/s
/tmp/testfiles/File5.txt,50MB,14.498499084s,3.45 MB/s
/tmp/testfiles/File6.txt,50MB,14.844652081s,3.37 MB/s
/tmp/testfiles/File7.txt,50MB,14.286648744s,3.50 MB/s
/tmp/testfiles/File8.txt,50MB,14.35213191s,3.48 MB/s
/tmp/testfiles/File9.txt,50MB,14.667853445s,3.41 MB/s
/tmp/testfiles/File10.txt,50MB,14.682988659s,3.41 MB/s
```


## Release Notes
The release notes are available [here](RELEASE.md).