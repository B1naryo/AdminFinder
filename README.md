# Admin Finder in Go

This is a simple program written in Go to find directories on a website. It checks if the directories listed in a file are present on a provided website.

## How to Use

1. Make sure you have Go installed on your machine. You can download and install it from the [official website](https://golang.org/).

2. Clone this repository to your local machine: \
git clone https://github.com/your-user/admin-finder-go.git


3. Navigate to the project directory:
cd admin-finder-go

4. Compile the code:
go build AdminFinder.go


5. Run the program:
- To check a single URL:
  ```
  ./AdminFinder -u https://site.com -d directories.txt
  ```
- To check multiple URLs from a file:
  ```
  ./AdminFinder -f urlz.txt -d directories.txt
  ```

Replace `https://site.com` with the base URL you want to check and `directories.txt` with the name of the file containing the directories to be checked.

## Command-Line Options

- `-u`: Specifies a single base URL.
- `-f`: Specifies a file containing multiple base URLs.
- `-d`: Specifies the file containing the directories to be checked (default: `directories.txt`).

## Directories File

The `directories.txt` file should contain the directories you want to check, with each directory on a separate line.
