Go MongoDB Project
==================

Overview
--------

This project demonstrates how to build a Go server that interacts with MongoDB using the `go-mongo-driver`. The server manages a MongoDB database called `movieDB`, which includes two collections: **movies** and **casts**.

-   **Movies Collection**: Contains information about movies, including a reference to the **casts** collection.
-   **Casts Collection**: Contains information about the cast members of each movie.

Data Models
-----------


```
type Movie struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Title           string             `json:"title"`
	Year            int                `json:"year"`
	Genres          []string           `json:"genres"`
	Href            string             `json:"href"`
	Extract         string             `json:"extract"`
	Thumbnail       string             `json:"thumbnail"`
	ThumbnailWidth  int                `json:"thumbnail_width"`
	ThumbnailHeight int                `json:"thumbnail_height"`
	CastID          primitive.ObjectID `bson:"cast_id,omitempty"`
}
```


```
type Cast struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	MovieID primitive.ObjectID `bson:"movie_id"`
	Cast    []string           `json:"cast"`
}
```

Setup
-----

### Prerequisites

-   Go (1.17 or later)
-   Docker
-   Docker Compose

### Running the Server

#### 1\. Start MongoDB Containers

Before running the server, make sure MongoDB is running. Use Docker Compose to start MongoDB containers:

`make mongodb-up`


#### 2\. Seed MongoDB

After the MongoDB container is up, seed your MongoDB with initial data:

`make mongodb-seed`

#### 3\. Run the Go Server

Start the Go server:

`make run`

#### 4\. Stopping MongoDB

To stop and remove MongoDB containers:

`make mongodb-down`

Postman API Collection
----------------------

For API testing, a Postman collection is provided in the `/docs` folder of the project. Import the collection into Postman, and you'll be able to test the available API endpoints easily.

**Steps to import in Postman**:

1.  Open Postman.
2.  Click on "Import" in the top left.
3.  Select "File" and upload the Postman collection from the `/docs` folder.
4.  Once imported, you can try out the API requests directly in Postman.

Configuration
-------------

The server configuration is specified in the `config.json` file, which includes necessary details such as MongoDB connection parameters and server settings. Ensure that this file is set up correctly before running the server.

Makefile Targets
----------------

| Target | Description |
| --- | --- |
| `run` | Runs the Go server using the configuration file. |
| `mongodb-seed` | Seeds MongoDB with initial data. |
| `mongodb-up` | Starts MongoDB containers using Docker Compose. |
| `mongodb-down` | Stops and removes MongoDB containers. |

Troubleshooting
---------------

If you encounter issues, ensure the following:

1.  **MongoDB is running**: Use the `make mongodb-up` command to start the MongoDB containers.
2.  **Configuration file is correct**: Double-check the `config.json` file for any misconfigurations.
3.  **Docker and Docker Compose are properly installed**: Ensure that Docker and Docker Compose are installed and set up on your system.