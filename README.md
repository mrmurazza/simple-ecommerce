# Quickstart Guide to Set up Dev Workspace

1. Git clone git@github.com:mrmurazza/simple-ecommerce.git
2. Make sure you have Go installed and set up Go workspace
3. Install the dependencies by running `go mod vendor`
4. Go to your project folder using `cd` and run your project using:

   `$ go run main.go`

5. For deployment, run this command to build binary file

   `$ go build main.go`

# Project Structure

```
.
├── config                          > contains config file that reads from env variable
├── domain                          > contains the main business logic grouped by its domain
│     ├── user                      > contains domain entity, interface & its implementation of domain user, consist of Login, Register, and Get Users
│     │  ├── impl                   > contains class implementation of service & repo
│     │  ├── entity.go              > file contains domain model and data object related to the domain 
│     │  └── type.go                > file contains interface declaration of service & repo
│     └── order                     > responsible domain order, consisting: order creation, showing products, & Get Orders
│        ├── impl                   
│        ├── entity.go              
│        └── type.go                
├── dto                             > contains data object used in endpoint requests & response
├── handler                         > contains handler class responsible to handle inputs from external to domain & parsing response from domain to external
├── pkg                             > contains package or util class used by the domain
│     ├── auth                      > contains JWT related code to generate Token & ParseToken, also logic to authenticate Admin or normal user via middleware
│     ├── database                  > contains code to init DB connection and running simple migrations
│     └── ratelimit                 > contains middleware logic to rate limit users API calls
├── README.md
└── main.go                         > Main Go file, contains service initialization, routing, and middleware assignment


```

# Project Component

Basically, this project is separated into several components:

1. Handler : responsible on preparing the inputs and serving the outputs and the one to call the business logic code.
2. Domain: contains main business logics consist of Type, Service, Repository, & Entity/Model

   2.1. Type : contains interface declaration of service & repository, one of dependency inversion implementation

   2.2. Service : contains service implementation responsible on handling the main core of the business logic. (
   ex: `domain/user/impl/service.impl.go`)

   2.3. Repository : responsible on handling DB queries. (ex: `domain/user/impl/repo.impl.go`)

   2.4. Entity : contains data object (struct/entity) highly related with core domain or event its domain model itself (
   ex: `domain/user/entity.go`)
3. DTO (Data Transfer Object): contains data object related to API requests & responses
4. PKG (Package): contains utils or specific package used by the core domain

# DB Schema
Based on the requirements, this application use 4 tables:
1. Users 
    To store user with its credentials & roles

    ```
    CREATE TABLE users (
        `id` integer NOT NULL primary key AUTOINCREMENT, 
        email varchar(50) not null default ``, 
        name varchar(50) not null default ``, 
        password varchar(50) not null default ``, 
        role varchar(50) not null, 
        created_at datetime not null default current_timestamp, 
        updated_at datetime not null default current_timestamp 
    );
    ```

2. Products
    product that user can see and select to order
    ```
    CREATE TABLE products (
        `id` integer NOT NULL primary key AUTOINCREMENT, 
        name varchar(50) not null default ``, 
        price integer not null default ``, 
        qty integer not null default ``, 
        description varchar(50) not null default ``, 
        image varchar(50) not null default ``, 
        created_at datetime not null default current_timestamp, 
        updated_at datetime not null default current_timestamp 
    );
    ```

3. Orders
    a record that store users order, this will have many to one relation with `users`.
    ```
    CREATE TABLE orders (
        `id` integer NOT NULL primary key AUTOINCREMENT, 
        customer_id integer not null , 
        status varchar(50) not null default ``, 
        total_qty integer not null, 
        total_amount integer not null, 
        created_at datetime not null default current_timestamp, 
        updated_at datetime not null default current_timestamp, 
        FOREIGN KEY (customer_id) REFERENCES customers(id)
    );
    ```

4. Order Units
    a smaller entity of order that store the breakdown data of one order for each products. Also contains a snapshot of selected product data for historical record purposes. this have relateion many to one with `orders`.
    ```
    CREATE TABLE order_units (
        `id` integer NOT NULL primary key AUTOINCREMENT, 
        order_id integer not null, 
        product_id integer not null, 
        qty integer not null default 1,
        price integer not null default 0, 
        name varchar(50) not null default ``,
        description varchar(50) not null default ``, 
        image varchar(50) not null default ``, 
        created_at datetime not null default current_timestamp,
        FOREIGN KEY (order_id) REFERENCES orders(id)
    );
    ```