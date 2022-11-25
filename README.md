# Signing service

Implement a signing service using a message driven / microservice oriented solution.

Given a database of 100,000 records and a collection of 100 private keys, create a process to concurrently sign batches of records, storing the signatures in the database until all records are signed.

## Documentation

[Documentation on Notion](https://hxt365.notion.site/Signing-service-b51923f1e6004a8c80b8f3dd26ef9b0b)

## How to run

- Run docker-compose for Coordinator to start Coordinator service and Database
- Seed DB data using script/prepare_data
- Config the number of workers and run docker-compose for Worker