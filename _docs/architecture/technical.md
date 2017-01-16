# Server
## Entry Point
The entry point for the server program is in `server/main.go`. This file contains a main method. In this method an `AppContext` 
structure is created with all the required values (See the server routing section of this file for more detail on 
`AppContext`).

After the `AppContext` has been created a new `HandlerLoader` is created. This specific `HandlerLoader` is implemented in 
`server/handlers/handlers.go`. The goal of the `HandlerLoader` interface is to bootstrap a series of handlers onto an `http.ServeMux` 
without having to worry about the handlers themselves.

After the `HandlerLoader` is created it's `Load()` method is called (which registers all of our applications HTTP Handlers 
with the provided `http.ServeMux`) we call `http.ListenAndServe` to start serving our registered handlers with the Golang 
HTTP stack.

## Routing
This section discusses the contents of `server/handlers/handlers.go`.  

The most important method in this file is the `Load` function for the `Loader`. This is the method that gets called to 
register all the application handlers with the provided `http.ServeMux`. 

Several helper functions including `registerDir`, `registerFile` and `registerEndpoint` are provided.

### Endpoint Handlers
The meat of the Squad Up routing system comes into play with the `registerEndpoint` method. This method allows you to 
implement a simpler version of an `http.Handler` called an `EndpointHandler`.
 
 All you have to do to implement an `EndpointHandler` is create a struct which implements the method:     
 
 ```
 Serve (ctx *models.AppContext, r *http.Request) (interface{}, *models.APIError)
 ```
 
 Then every time the url you registered your handler for gets a request the `Serve` method is called. The routing system 
 expects an `EndpointHandler` to return an interface to marshall and send to the client and a reference to an error. This 
 allows one to program HTTP routes with much ease.


# Event planning
The core of Squad Up is its ability to plan events given a set of constraints. This section describes 
how the the event planning algorithm operates.

## Steps
The event planning process is broken down into a series of steps. This allows for a solid data processing 
structure and also happens to mimic the flow in which events are planned by groups.  

The steps are as follows:

- 1. Data Sourcing
    - The first step is to gather data. After all we can't run any fancy filters or algorithms if we don't 
      have any data to run it on.
    - Data is gathered and then stored as rows of the Event Proposal table.
        - Each entry represents a *proposed event schedule*. This schedule include the location and a time block.
        - These entries are filtered later in the second step
- 2. Data Filtering
    - Now that we have data to work with we have to filter it.
    - We do this by assigning each Event Proposal a weight, on a scale of 0 - 100.
    - The weight of each entry is determined by the implementation.
- 3. Commit
    - After the ideal proposal has been chosen and verified by the group the proposal details must be copied to the main 
      Event entry. 
      
## Implementation
A general description of each of the following steps can be found in the [models documentation](/docs/architecture/models.md) 
under the `type` field of the Event model.
      

