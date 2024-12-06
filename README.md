# MarkDog

MarkDog is a lightweight, user-friendly markdown editor that provides real-time preview and easy export options. Designed for developers and writers, MarkDog simplifies the process of creating and managing markdown documents.

## Features

- **Real-Time Markdown to HTML Conversion**: Instantly see your markdown content rendered as HTML.
- **PDF Export**: Convert your markdown documents to PDF format with ease.
- **Web-Based Interface**: Access the editor through a clean and responsive web interface.
- **File Upload Support**: Upload and edit your markdown files directly in the application.
- **Efficient Web Server**: Powered by the Fiber framework for fast and reliable performance.

## Installation

1. Ensure you have Go installed on your machine.
2. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/markdog.git
   ```
3. Navigate to the project directory:
   ```bash
   cd markdog
   ```
4. Install dependencies:
   ```bash
   go mod tidy
   ```
5. Run the application:
   ```bash
   go run main.go
   ```
6. Access the application at `http://localhost:3050`.

## Usage

- **Convert Markdown**: Use the `/convert` API endpoint to convert markdown text to HTML.
- **Export PDF**: Use the `/export-pdf` API endpoint to convert markdown to a PDF file.
- **Upload Files**: Use the `/upload` API endpoint to upload and read markdown files.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to the creators of the Fiber framework and the `gomarkdown/markdown` and `go-wkhtmltopdf` packages for their excellent tools that made this project possible.

---

Enjoy using MarkDog for all your markdown editing needs!
