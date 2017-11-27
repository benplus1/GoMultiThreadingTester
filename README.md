# GoMultiThreadingTester

This is a program that stress tests URL webservices & prints out statistics for bad responses and errors given a text file of links to call.


## Setup

Make sure the latest version of Go is installed before cloning and downloading the project.

Go: https://golang.org/project/


### Running the code

To run the main project, go to the "test" folder.

Make sure to create a textfile with URLs of services you wish to call and test. A file within the test folder called test.txt is included.

In that folder, run in the command line:
```
$ go run MultiThreadingTester.go
```
You will then be asked to input the file path of the textfile. For test.txt, it would be:

```
../GoMultiThreadingAppTester/test/test.txt
```
The program will then run until a specified number of seconds (inputted by the user).

## Credit

Built by Benjamin Yang at Rutgers University, kept under the MIT License.
