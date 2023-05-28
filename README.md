# ApiCat

English | [简体中文](https://github.com/apicat/apicat/blob/master/README-CN.md)

ApiCat is an API development tool based on AI technology, which aims to help developers develop APIs more quickly and efficiently through automation and intelligence. ApiCat supports the import and export of OpenAPI and Swagger data files, and can analyze and identify the API requirements entered by users, and automatically generate corresponding API documents and codes.

You can visit our [Online Demo](http://demo.apicat.net) to try it out.

ApiCat is still in its early stages, Star and Watch are welcome to follow the latest developments of the project.

## Features

### Demo

![AI-generate-schema](https://cdn.apicat.net/uploads/0c3518c1bfc421fc4f3f86c085f353d2.gif)

![AI-generate-api-by-schema](https://cdn.apicat.net/uploads/bbcae83511d797d22077d05d17c262cc.gif)

![AI-generate-api](https://cdn.apicat.net/uploads/cf617b56fa186960c228c79487cf6c5e.gif)

### Overview

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
cd frontend 
pnpm install
pnpm build

# Update collation dependencies
go mod tidy

# Compile project
go build

# Modify the configuration file
# You can copy the configuration file content of config/setting.default.yaml for configuration modification
cp ./config/setting.default.yaml ./
vim ./setting.default.yaml

# start service(default configuration)
./apicat
# start service(custom configuration)
./apicat -c setting.default.yaml
```

## Contact

The growth of ApiCat is inseparable from each of its users. If you have any content that you want to discuss with us, please contact us and join our WeChat discussion group through the QR code below.

![Wechat Group](https://cdn.apicat.net/uploads/01bfb23802cdfad49f0d560ee80fc5e3.png)

## License

[MIT](https://github.com/apicat/apicat/blob/main/LICENSE)