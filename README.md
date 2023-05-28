# Document Analysis and Image Comparison Application

This application, written in GoLang, takes input fields, a selfie, and a picture of a document and uses ML services to analyse and compare the text in the document image to the input values and  the image in the document to the selfie.

## Features

- Extracts text from a document image using Amazon Textract.
- Compares the extracted text with the input values provided by the user.
- Uses Amazon Rekognition to compare the image in the document with the selfie.
- Provides results indicating the match or mismatch between the document and the input values.

## Prerequisites

Before running the application, make sure you have the following:

- An AWS account with appropriate permissions to access Amazon Textract and Amazon Rekognition services.
- Go programming language installed on your development machine.
- AWS SDK for Go installed.
- Docker installed

## Getting Started

Follow these steps to set up and run the application:

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/document-analysis-app.git

2. Install the project dependencies:

    ```bash
    cd doc-verify
    go mod download

3. Configure AWS credentials:

    Set up your AWS credentials by either configuring the AWS CLI or by setting the

    `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables.

4. Configure application settings:

    Create your env file and setup the enviroment variables using `sample.env` as a guide.

5. Build and run the application:

    ```bash
    make build
    make start-db
    make run

## Usage

## Architecture

- Backend: Written in GoLang, the backend server receives requests from the API endpoints,
  
  interacts with AWS services (Amazon Textract and Amazon Rekognition)
  
  to perform document analysis and image comparison, and returns the results.

## Contributing

Contributions to the Document Analysis and Image Comparison Application are welcome! If you have any ideas, improvements, or bug fixes, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. Feel free to modify and use it according to your needs.

## Tools

- golang
- gin
- aws-sdk-go
- gorm
- mongodb
- go-mongo-driver
- godotenv
- golang-jwt
