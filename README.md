# Mekari Test

## How to run
<p align="justify">
Because this application and its database are already containerized using Docker, there are no specific prerequisites for this app. You can get it up and running by simply executing the following command, which is defined in the Makefile:

```azure
make run
```
<p align="justify">
This command will handle all the necessary steps to start the application and its database in their respective containers. Make sure you have Docker installed and running on your system, as it is the only requirement to run this application.

## Project Structure
<p align="justify">
This project follows a clean code architecture, which divides the code into four distinct layers:

- Model/Domain: This layer represents the core domain models and entities of your application.
- Repository: The repository layer is responsible for data access and storage operations.
- Service: This layer contains the business logic and services that operate on the data.
- Handler/Controller: The handler or controller layer manages the application's endpoints and interacts with the outside world.

## Unit Test
<p align="justify">
For unit testing, this project covers two endpoints: 

``
"Get Employee by ID" and "Register Employee."
``


<p align="justify">
Each layer of the application, including the repository, service, and handler, has its own set of unit tests. The testing strategy prioritizes covering test cases rather than aiming for high code coverage.

## API Documentation
<p align="justify">
For detailed information about the API, including input and output specifications, please refer to the following link:

https://app.theneo.io/ef00f0b3-a206-4882-abf7-b7807b8afc83/mekari/employees