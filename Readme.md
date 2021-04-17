### How to run
- Config `.netrc` so that you can pull private go repo

    Copy `.netrc.example` to create your own `.netrc` and place it
    at the same place as `.netrc.example` so that the docker-compose
    is able to pull the private repo. Also, copy your created `.netrc`
    to `~/.netrc` so that your IDE such as `goland` will be able to
    pull the go mode private repo. Alternative way to allow your IDE to
    pull go mod private repo is to config the git by enabling ssh-key
    authentication. You first let git to replace the `https` url with
    `git` by running.
    
    ```bash
    git config --global url."ssh://git@github.com/secure-for-ai/".insteadOf "https://github.com/secure-for-ai/"
    ```
    
    Then, `go mod download` will be able to pull the private repo if
    ssh-key is correctly configured.

- Setup the environment `GOPRIVATE=github.com/secure-for-ai` in IDE such as 
  goland so that we can pull our private package directly instead of pulling
  through the proxy, which is the default pulling.

- Start docker compose
    ```bash
    docker-compose up
    ```

- Initial Mongo
    ```bash
    # login mongo
    mongo --port 27017 --host=localhost --authenticationDatabase=admin \
        -p password --username test
    ```
    
    ```bash
    # setup replication with single master
    rs.initiate({_id: "rs0", members: [{_id: 0, host: "localhost:27017"}] })
    ```
