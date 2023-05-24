`Powered by Chat-GPT`

# iggo - Image Gallery Generator

Iggo is a command-line tool written in Go, which creates an HTML image gallery from images in a specified directory. The gallery is paginated with a set number of images per page.

## Installation

This project requires Go 1.16 or later. You can download it from the [official Go website](https://golang.org/dl/).

To install iggo, use the `go install` command with the path to the repository:

```bash
go install github.com/sat8bit/iggo@latest
```

This will download the source code and create an executable named `iggo` in your `$GOPATH/bin` directory.

## Usage

```
iggo -p [NUMBER] [TEMPLATE_PATH] [INPUT_DIRECTORY] [OUTPUT_DIRECTORY]
```

- `-p [NUMBER]` (optional): Specify the number of images per page. The default is 100.
- `[TEMPLATE_PATH]`: Path to the HTML template file.
- `[INPUT_DIRECTORY]`: Path to the directory containing the images.
- `[OUTPUT_DIRECTORY]`: Path to the directory where the gallery should be created.

If the number of arguments is not enough or incorrect, the program will display an error message with the correct usage.

## Example

To create a gallery with 50 images per page, using a template at `./template.html`, sourcing images from `./input`, and writing the gallery to `./output`, you would use the following command:

```bash
iggo -p 50 ./template.html ./input ./output
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the terms of the MIT license.
