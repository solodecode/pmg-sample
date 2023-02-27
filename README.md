# pmg-sample

An example of simple service of gRPC, postgreSQL with basic authentication and mTLS security
![Diagramm](/png/diag.png)

## Usage
Clone repo into your directory:
```bash
git clone https://github.com/solodecode/pmg-sample
```

Change the constant dbUrl in server.go and use sql-query create.sql in sql directory. Build and run server.go and client.go