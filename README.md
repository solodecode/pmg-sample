# pmg-sample

An example of simple service of gRPC, postgreSQL with basic authentication and mTLS security
![Diagramm](/assets/diag.png)

## Usage
Clone repo into your directory:
```bash
$ git clone https://github.com/solodecode/pmg-sample
```

Change the constant dbUrl in server.go and use sql-query create.sql in sql directory. Build and run server.go and client.go
### OR
You can use docker to run pre-built images:
1. Pull images:
```bash
$ docker pull ghcr.io/solodecode/pmg-sample:latest
$ docker pull ghcr.io/solodecode/pmg-sample/client:latest
$ docker pull postgres:latest
```
2. Create network for images:
```bash
$ docker network create pg-pmg-sample
```
3. Run PostgreSQL:
```bash
docker run --network pg-pmg-sample --name pg -p 5432:5432 -e POSTGRES_PASSWORD=test -e POSTGRES_DB=test -d postgres:latest
```
4. Open PostgreSQL terminal in docker and connect to test database:
```bash
$ psql -U postgres
$ \c test
```
5. And now you can create the table by sql-query in pmg-sample/sql:
```bash
$ CREATE TABLE test (
$ id BIGSERIAL NOT NULL,
$ title VARCHAR(100),
$ description TEXT,
$ price FLOAT,
$ stock BOOL
$ );
```
6. Run server and client images:
```bash
$ docker run --network pg-pmg-sample --name server -p 5333:5333 ghcr.io/solodecode/pmg-sample
$ docker run --network pg-pmg-sample --name client ghcr.io/solodecode/pmg-sample/client
```
Now you will receive logs of adding a product to the table every 5 seconds.