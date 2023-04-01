# pmg-sample

An example of simple service of gRPC, postgreSQL with basic authentication and mTLS security
![Diagramm](/assets/diag.png)

## Usage with pre-built images:
1. Clone repo into your directory:
```bash
$ git clone https://github.com/solodecode/pmg-sample
```
2. Pull docker images:
```bash
$ docker pull ghcr.io/solodecode/pmg-sample/cmd/server:latest
$ docker pull ghcr.io/solodecode/pmg-sample/cmd/client:latest
$ docker pull postgres:latest
```
3. Create network for images:
```bash
$ docker network create pmg-sample
```
3. Run PostgreSQL:
```bash
$ docker run --network pmg-sample --name pg -e POSTGRES_DB=test -e POSTGRES_PASSWORD=test -v $(pwd)/scripts/sql/create.sql:/docker-entrypoint-initdb.d/init.sql postgres
```
4. Run server and client images:
```bash
$ docker run -e DATABASE_URL=postgres://yourlogin:yourpwd@pg:5432/test --name server -p 5333:5333 --network pmg-sample ghcr.io/solodecode/pmg-sample/cmd/server
$ docker run --network pmg-sample ghcr.io/solodecode/pmg-sample/cmd/client
```
Now you will receive logs of adding a product to the table every 5 seconds.
