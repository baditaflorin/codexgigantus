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

### Basic Command
```sh
./llm-codeinject --directory /path/to/dir --ignore-dir logs,temp --ignore-ext log,tmp --include-ext txt,md
```

### Flags Explanation
- `--dir` or `-dir`: Comma-separated list of directories to search (default: current directory).
- `--ignore-file` or `-ignore-file`: Comma-separated list of files to ignore.
- `--ignore-dir` or `-ignore-dir`: Comma-separated list of directories to ignore.
- `--ignore-ext` or `-ignore-ext`: Comma-separated list of file extensions to ignore.
- `--include-ext` or `-include-ext`: Comma-separated list of file extensions to include.
- `--recursive` or `-recursive`: Recursively search directories (default: true).
- `--debug` or `-debug`: Enable debug output.
- `--save` or `-save`: Save the output to a file.
- `--output-file` or `-output-file`: Specify the output file name (default: output.txt).
- `--show-size` or `-show-size`: Show the size of the result in bytes.
- `--show-funcs` or `-show-funcs`: Show only functions and their parameters.

### Internal Use Examples

#### Frontend
```sh
llm-codeinject -dir social-network-frontend -ignore-file package-lock.json -ignore-dir node_modules,__previewjs__ -ignore-ext svg,png,ico,md -output-file frontend.txt -save
```

#### Backend
```sh
llm-codeinject -dir social-network-backend -ignore-file package-lock.json,auth_test.go -ignore-dir tests -ignore-ext sum,mod -output-file backend.txt -save
```

#### Debugging
```sh
llm-codeinject -debug -dir . -ignore-file package-lock.json,llm-codeinject,frontend.txt -ignore-dir cmd,pkg,.idea,.git,node_modules,__previewjs__ -ignore-ext svg,png,ico,md -output-file frontend.txt -save
```

## Development

The code has been refactored to improve modularity and testability using functional programming principles. Key changes include:

- **Configuration Parsing**: Decoupled from the main function for better readability and testability.
- **File Processing**: Extracted into a separate function to simplify the main function.
- **Functional Style**: Refactored file handling using higher-order functions.
- **Debug Information**: Improved handling using a functional approach.
- **Utility Functions**: Consolidated into a single module for improved organization and reusability.

### Testing

To ensure the code is easy to test, functional parameters are used for gathering and processing files, allowing easy mocking during tests. Each function handles a single responsibility, making the codebase modular and maintainable.

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss any changes.

## License

This project is licensed under the MIT License.
