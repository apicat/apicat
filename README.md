# ApiCat

English | [简体中文](https://github.com/apicat/apicat/blob/master/README-CN.md)

ApiCat is an API development tool based on AI technology, which aims to help developers develop APIs more quickly and efficiently through automation and intelligence. ApiCat supports the import and export of OpenAPI and Swagger data files, and can analyze and identify the API requirements entered by users, and automatically generate corresponding API documents and codes.

You can visit our [Online Demo](http://demo.apicat.net) to try it out.

ApiCat is still in its early stages, Star and Watch are welcome to follow the latest developments of the project.

## Features

- Support OpenAPI and Swagger data file import and export, which is convenient for developers to describe and manage API specifications.
- Through AI technology, the requirements and structure of the API can be automatically identified, and corresponding API documents and codes can be generated to improve development efficiency and quality.

## Installation and deployment

### Get Code

```
git clone https://github.com/apicat/apicat.git
```

### Compile and start the service

```
# Enter project
cd apicat

# Compile the front-end code
cd fronted
pnpm install
pnpm build

# Update collation dependencies
go mod tidy

# Compile project
go build

# Modify the configuration file
vim ./setting.example.yaml

# start service
./apicat
```

## License

[MIT](https://github.com/apicat/apicat/blob/main/LICENSE)