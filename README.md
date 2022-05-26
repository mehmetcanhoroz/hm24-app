# HM24 Case App

This repo was developed for hm24 case. Built-in Golang.

## Installation

You only need to download go modules.

```bash
go mod tidy
```

## Running Service

This application is a basic api. In main directory, you should execute following command.
```bash
go run cmd/main.go 
```

## Using Service

Postman collection is prepared for you to see all endpoints and ready to use.
2 Files are ready in the repo, `HM24-App.postman_collection.json` and `HM24-App - Local.postman_environment.json`
You should import them into postman.
Please follow this url to learn how to import:

https://learning.postman.com/docs/getting-started/importing-and-exporting-data/#importing-data-into-postman

After that, please select environment on the right top corner.
Finally, left-side menu, you can see collections and use HM24 case collections to test app.


## Endpoints

| Case Task                             | URL                             |
|---------------------------------------|---------------------------------|
| Test                                  | Test                            |
| HTML Version                          | /analyse/html-version?url={URL} |
| Page Title                            | /analyse/title?url={URL}        |
| Headings count by level               | /analyse/hx?url={URL}           |
| Amount of internal and external links | /analyse/links?url={URL}        |
| Amount of inaccessible links          | /analyse/links?url={URL}        |
| If a page contains a login form       | /analyse/login-form?url={URL}   |

## Contributing

Pull requests are welcome. Anything want to discuss, please open an issue then we can have a chat what you would like to
change.
As I mentioned, It is for case. So, I already know there some part could be improved or refactored. However, Just
completed to provide expected case.

## License

[MIT](https://choosealicense.com/licenses/mit/)