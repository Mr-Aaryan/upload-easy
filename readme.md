# Upload Easy

A simple CLI tool to upload files to **Google Drive**, **Mega** and **Cloudinary**. Upload files directly from your terminal with ease, just like using an npm package.

## Features

- Upload files to Google Drive.
- Upload files to Cloudinary.
- Select the desired upload service based on your preference.
- Simple and intuitive CLI interface.

## Installation

1. Install Go if you don't have it already. [Download Go](https://go.dev/dl/).
2. Clone this repository:
   ```bash
   git clone https://github.com/Mr-Aaryan/upload-easy.git
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
go run main.go --file "./upload/file.png" -g
```

### Options

- `--file` (required): Path to the file to be uploaded.
- `-g` or `-c` or `-m` (required): Choose between where to upload. -`-g` for Google -`-m` for Mega -`-c` for Cloudinary

Example:

```bash
go run main.go --file "./upload/file.png" -g
```

### Configuration

Before using the tool, set up your environment variables:

#### Google Drive

1. Obtain credentials for Google Drive API by following [this guide](https://developers.google.com/drive/api/v3/quickstart/go).
2. Save the credentials JSON file as `credentials.json` in the project directory as `./googleutils/credentials.json`
3. `./googleutils/token.json` File is automatically created after successful authentication


#### Cloudinary

1. Log in to your Cloudinary account and obtain your API key and secret.
2. Set up the `.env` file with the following variables:
   ```env
   CLOUDINARY_URL=cloudinary://<API_KEY>:<API_SECRET>@<CLOUD_NAME>
   ```

#### Mega
1. Create an mega account if you don't already have it.
2. Add the following in your `.env` file
   ```env
   MEGA_EMAIL=<mega_email>
   MEGA_PASSWORD=<mega_password>
   ```

### Examples

#### Upload to Google Drive

```bash
go run main.go --file "./upload/file.png" -g
```

#### Upload to Cloudinary

```bash
go run main.go --file "./upload/file.png" -c
```

#### Upload to Mega

```bash
go run main.go --file "./upload/file.png" -m
```

### Prerequisites

- Go 1.20+ installed.
- Environment variables configured in `.env`.

## Roadmap

- [ ] Add support for more cloud storage services.
- [ ] Implement parallel uploads.
- [ ] Provide support for folder uploads.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

