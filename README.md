# **git-repo-filesize-scanner**

A Command Line Tool that enable users scan through a public repository for large files based off a defined threshold (MB) and return a summary of those files along with the size of the files.

---

## **Set up**

### **Building git-repo-filesize-scanner Application**
- Clone github to your local system
- To build the application run the command ` go build -o grfscan . ` to build the binary
- Run the command ` ./grfscan --help ` to view all commands available
- To scan a repository using inline json ` ./grfscan scan --json '{"clone_url":"https://github.com/example.git","size":5, "token": "ghp_randomtoken"}' `
- N.B the token is optional as it only accesses public repositories
- To view all the scan commands available use ` ./grfscan scan --help `

  **Sample Response:**  
  ```json
    { 
    "total_num_of_files": 4,
    "files": [
        {
        "name": "server/tests/large_lzw_frame.gif",
        "size": "7.03 MB"
        },
        {
        "name": "server/tests/orientation_test_9.jpeg",
        "size": "6.83 MB"
        },
        {
        "name": "webapp/channels/src/images/crt-in-product.gif",
        "size": "7.58 MB"
        },
        {
        "name": "webapp/channels/src/images/emoji-sheets/apple-sheet.png",
        "size": "16.76 MB"
        }
    ]
    }
  ```

### **Technologies Used**
- Golang, Cobra for the command line tooling, github.com/go-git/go-git/v6 for making git client calls, uber zap logger for logging.
