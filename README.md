# sam-app-go

This is a modified sample app from aws `sam init` using Go.

## Requirements

* AWS CLI already configured with Administrator permission
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
* MySQL Database (for local development)
* AWS RDS (for AWS deployment)

## Setup process

### Install MySQL

Make sure MySQL is installed in your local. The table this CRUD application is going to use is called `students`. The database name is "crud":

``` SQL
create database crud;

use crud;

CREATE TABLE `students` (
  `id` int NOT NULL AUTO_INCREMENT,
  `first_name` varchar(256) NOT NULL,
  `last_name` varchar(256) NOT NULL,
  `email` varchar(256) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  UNIQUE KEY `email_UNIQUE` (`email`)
) ENGINE=InnoDB;

insert students (first_name, last_name, email) values ('Jim', 'Rome', 'jim.rome@gmail.com')
```

### Prepare the DB connection string

To connect to local, the mysql connection string for Go would be:

```
<username>:<password>@tcp(<your machine name>.local:3306)/crud
```

### Generate template.yaml

`sam` requires a `template.yaml` either to run it locally or in AWS Lambda. Please note that this project does not include the required `template.yaml`. It has to be generated through the `template.goyaml` by running:

```
cd crud-cli
go install
cd ..
~/go/bin/crud-cli template --set ConnectionString='root:Welcome1@tcp(<your hostname>.local:3306)/crud' -f ./template.goyaml > template.yaml
```

It will generate the `template.yaml` into the project root directory. The connection string above uses `root` as the username, `Welcome1` as the password, and `<your hostname>` as the hostname. Please also note that `sam` is running in a Docker container locally. `localhost` or `127.0.0.1` will not work as it will connect to the docker container locally.

### sam build

In this example we use the built-in `sam build` to automatically download all the dependencies and package our build target.   
Read more about [SAM Build here](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-build.html) 

 
```shell
sam build
```

### Local development

**Invoking function locally through local API Gateway**

```bash
sam local start-api -p 3001
```

If `sam `started successfully, you should be able to see REST endpoints as:

```
http://127.0.0.1:3001/hello [GET]
http://127.0.0.1:3001/crud/student [DELETE, GET, HEAD, OPTIONS, PATCH, POST, PUT]
http://127.0.0.1:3001/crud/students [GET]
```

```
% curl http://127.0.0.1:3001/hello 
Hello, 160.34.93.102
```

### REST endpoints

List students

```
% curl http://127.0.0.1:3001/crud/students
{
  "Count": 1,
  "Students": [
    {
      "Id": "1",
      "FirstName": "Jim",
      "LastName": "Rome",
      "Email": "jim.rome@gmail.com"
    }
  ]
}

```

Please see the [frontend repo](https://github.com/ssdSSF/React-CRUD-Operation-V2) for the full REST API invokement.

### Production deployment

Make sure you have aa [AWS RDS DB](https://aws.amazon.com/rds/). The connection string will be very similar to:

```
admin:Welcome1@tcp(sam-app.cxxxxxxx3.us-west-2.rds.amazonaws.com:3306)/crud
```

Generate the `template.yaml` for Lambda deployment that will connect to AWS RDS:

```
~/go/bin/crud-cli template --set ConnectionString='admin:Welcome1@tcp(sam-app.cxxxxxxx3.us-west-2.rds.amazonaws.com:3306)/crud' -f ./template.goyaml > template.yaml
```

Build again:
```
sam build
```

Deploy:
```
sam deploy --guided
```