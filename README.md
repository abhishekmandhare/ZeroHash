# ZeroHash VWAP Calculator

The goal of this project is to create a real-time VWAP (volume-weighted average price) calculation engine. App uses the coinbase websocket feed to stream in trade executions and update the VWAP for each trading pair as updates become available. 

### Commands

* Run app locally
    
    ```make run ```

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

## Documentation
