# Upload Easy

A simple CLI tool to upload files to **Google Drive** and **Cloudinary**. Upload files directly from your terminal with ease, just like using an npm package.

## Features

- Upload files to Google Drive.
- Upload files to Cloudinary.
- Select the desired upload service based on your preference.
- Simple and intuitive CLI interface.

## Installation

1. Install Go if you don't have it already. [Download Go](https://go.dev/dl/).
2. Clone this repository:
   ```bash
   git clone https://github.com/your-username/upload-easy.git
   cd upload-easy
   ```
3. Build and install the tool:
   ```bash
   go install
   ```
   Ensure the Go bin directory is in your PATH.

## Usage

### Upload a File

To upload a file, use the following command:
```bash
upload-easy --file "./upload/file.png"
```

### Options

- `--file` (required): Path to the file to be uploaded.
- `-g` or `-c` or `-m` (required): Choose between where to upload.
   -`-g` for Google
   -`-m` for Mega
   -`-c` for Cloudinary

Example:
```bash
upload-easy --file "./upload/file.png" --service "drive"
```

### Configuration

Before using the tool, set up your environment variables:

#### Google Drive

1. Obtain credentials for Google Drive API by following [this guide](https://developers.google.com/drive/api/v3/quickstart/go).
2. Save the credentials JSON file as `credentials.json` in the project directory.
3. Set up a `.env` file in the project directory:
   ```env
   GOOGLE_DRIVE_CREDENTIALS=./credentials.json
   ```

#### Cloudinary

1. Log in to your Cloudinary account and obtain your API key and secret.
2. Set up the \`.env\` file with the following variables:
   ```env
   CLOUDINARY_URL=cloudinary://<API_KEY>:<API_SECRET>@<CLOUD_NAME>
   ```

### Examples

#### Upload to Google Drive

```bash
upload-easy --file "./upload/file.png" --service "drive"
```

#### Upload to Cloudinary

```bash
upload-easy --file "./upload/file.png" --service "cloudinary"
```

#### Let the Tool Prompt You

```bash
upload-easy --file "./upload/file.png"
```
The tool will ask you to select the service.

## Development

### Prerequisites

- Go 1.20+ installed.
- Environment variables configured in `.env`.

### Run Locally

1. Build the project:
   ```bash
   go build -o upload-easy
   ```
2. Execute:
   ```bash
   ./upload-easy --file "./upload/file.png"
   ```

### Testing

Run tests using:
```bash
go test ./...
```

## Roadmap

- [ ] Add support for more cloud storage services.
- [ ] Implement parallel uploads.
- [ ] Provide support for folder uploads.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.


## Things to add 
./googleutils/credentials.json
//we should download from google console after creating google clientId and all

./googleutils/token.json 
//google downloads automatically after authentication


.env
```
CLOUDINARY_URL=<cloudinary_url>
MEGA_EMAIL=<mega_email>
MEGA_PASSWORD=<mega_password>
```