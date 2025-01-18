# CodexGigantus

CodexGigantus is a command-line tool written in Go that processes files in a specified directory based on given criteria. It's designed to integrate seamlessly with Language Learning Models (LLMs) for extracting smaller code components.

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
    git clone https://github.com/baditaflorin/codexgigantus.git
    ```
2. Navigate to the project directory:
    ```sh
    cd codexgigantus
    ```
3. Build the project:
    ```sh
    ./build.sh
    ```

## Usage

### Basic Command
```sh
./codexgigantus -dir /path/to/dir --ignore-dir logs,temp --ignore-ext log,tmp --include-ext txt,md
```

### How to test it on this repo
```shell
 ./CodexGigantus -dir . --ignore-file CodexGigantus,.DS_Store,qodana.yaml --ignore-ext txt --ignore-dir .git,.idea --save --output-file chatgpt_code.txt
```
### Flags Explanation
- `--dir` or `-dir`: Comma-separated list of directories to search (default: current directory).
- `--ignore-file` or `-ignore-file`: Comma-separated list of files to ignore.
- `--ignore-dir` or `-ignore-dir`: Comma-separated list of directories to ignore.
- `--ignore-ext` or `-ignore-ext`: Comma-separated list of file extensions to ignore.
- `--include-ext` or `-include-ext`: Comma-separated list of file extensions to include.
- `--ignore-suffix` or `-ignore-suffix`: Comma-separated list of file suffixes to ignore.
- `--recursive` or `-recursive`: Recursively search directories (default: true).
- `--debug` or `-debug`: Enable debug output.
- `--save`: Save the output to a file.
- `--output-file`: Specify the output file name (default: output.txt).
- `--show-size`: Show the size of the result in bytes.
- `--show-funcs`: Show only functions and their parameters.

### Internal Use Examples

#### Frontend
```sh
codexgigantus -dir social-network-frontend -ignore-file package-lock.json -ignore-dir node_modules,__previewjs__ -ignore-ext svg,png,ico,md -output-file frontend.txt -save
```

#### Backend
```sh
codexgigantus -dir social-network-backend -ignore-file package-lock.json,auth_test.go -ignore-dir tests -ignore-ext sum,mod -output-file backend.txt -save
```

#### Debugging
```sh
codexgigantus -debug -dir . -ignore-file package-lock.json,codexgigantus,frontend.txt -ignore-dir cmd,pkg,.idea,.git,node_modules,__previewjs__ -ignore-ext svg,png,ico,md -output-file frontend.txt -save
```

## Development

## Notes
Configuration Parsing: The ParseFlags function in config.go handles all command-line arguments.
File Processing: The ProcessFiles function in file_processor.go handles directory traversal and file filtering.
Functional Style: The code uses functional programming principles for better modularity and testability.
Debug Information: Use the -debug flag to enable detailed debug output.
Utility Functions: Common utility functions are consolidated in utils.go.

## Testing
The code is organized for easy unit testing.
Each function handles a single responsibility.
Functional parameters can be mocked during tests.

### Testing

To ensure the code is easy to test, functional parameters are used for gathering and processing files, allowing easy mocking during tests. Each function handles a single responsibility, making the codebase modular and maintainable.

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss any changes.

## License

This project is licensed under the MIT License.
