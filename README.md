# Order Pack Calculator

This repository contains a simple Order Pack Calculator that helps users determine the optimal way to pack items into predefined order packs. The calculator takes into account the amount of the items and the available pack sizes to minimize wasted space.

The available pack sizes are:
- Pack of 250
- Pack of 500
- Pack of 1000
- Pack of 2000
- Pack of 5000

## Features

- Input the total number of items to be packed.
- Calculates the optimal combination of pack sizes to minimize leftover items.
- Provides a clear output of how many packs of each size are needed.

## Rules

- Only whole packs can't be sent, no packs can be broken open
- The calculator yields the least amount of items to fullfil the order
- The calculator yields the least amount of packs to fullfil the order

## Tools

- Go programming language
- Fiber web framework for building the API
- Docker for containerization
- PostgreSQL for database management
- Goose for database migrations
- Swagger for API documentation

## How to run

1. Clone the repository:
   ```bash
   git clone github.com/luk3skyw4lker/gymshark-challenge-backend
   cd gymshark-challenge-backend
   ```

2. Start the Docker containers:
   ```bash
   make up_containers
   ``` 

3. Use the API endpoints to create orders and manage pack sizes. Refer to the Swagger documentation for detailed API usage.