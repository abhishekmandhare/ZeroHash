# ZeroHash VWAP Calculator

The goal of this project is to create a real-time VWAP (volume-weighted average price) calculation engine. App uses the coinbase websocket feed to stream in trade executions and update the VWAP for each trading pair as updates become available. 

## Requirements

1. **Docker**
2. **Golang 1.18**


### Commands

* Run app locally
    
    ```make run ```
    
For windows:
    
    make run-win

* Run in a docker
    
    ```make docker-run```

* Vendor 
    
    ```make vendor```

* Lint
    
    ```make lint```

    Run 
    ```dep-install```
    if the lint command fails.

* Test

    ```make test```

Instructions
-----

1. Clone this repository.
2. Create a new branch called `xxx`.
3. Create a pull request from your `xxx` branch to the main branch.

### Design

![VWAP design](/doc/design/design.PNG?raw=true "Data flow")

## VWAP calculation

VWAP is calculated using below formula
    
    - Sum(price*size) / (size)
    
VWAP is calculated using a maximum of 200 data points. A queue is used to maintain a window size of 200. If 201the element is pushed in, then the 1st element is ppoped out to maintain the window size.

### Future improvements:

1. Add integration testing by mocking the webscoket connection
2. Improve efficiency of queue by replacing datastructure with a more efficient one.

