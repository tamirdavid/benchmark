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
* Bouth 'dl' and 'up' commands output is csv files containg the filename, size, time it took to upload/download.
- for Example:
```
file,size,time_taken
File1.txt,50MB,14.76517576s
File2.txt,50MB,14.835556461s
File3.txt,50MB,14.561855731s
File4.txt,50MB,15.041486972s
```


## Release Notes
The release notes are available [here](RELEASE.md).