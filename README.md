# LLM-FileProcessor

LLM-CodeInject is a command-line tool written in Go that processes files in a specified directory based on given criteria. It's designed to integrate seamlessly with Language Learning Models (LLMs) for extracting smaller code components.

## Features
- Specify root directory
- Ignore specific directories
- Exclude files with specific extensions
- Include only files with specific extensions
- Process and display contents of text files

## Installation

### Prerequisites
- Go 1.16 or later

### Steps
1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/llm-fileprocessor.git
    ```
2. Navigate to the project directory:
    ```sh
    cd llm-codeinject
    ```
3. Build the project:
    ```sh
    ./build.sh
    ```

## Usage

```sh
./llm-codeinject --directory /path/to/dir --ignore-dir logs,temp --ignore-ext log,tmp --include-ext txt,md


Internal Use Examples
#frontend 
llm-codeinject -dir social-network-frontend -ignore-file package-lock.json -ignore-dir node_modules,__previewjs__ -ignore-ext svg,png,ico,md -output-file frontend.txt -save  #-show-size

#backend
llm-codeinject -dir social-network-backend -ignore-file package-lock.json,auth_test.go -ignore-dir tests -ignore-ext sum,mod -output-file backend.txt -save  #-show-size

