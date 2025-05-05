### Building and running your application

When you're ready, start your application by running:
`docker compose up --build`.

Your application will be available at http://localhost:8080.

### Deploying your application to the cloud

Make sure to config YOUR file named '.env' before the following
(you can take exemple from the /projet/exemple.env )

First, build your image at root : 
`docker build --tag goforum -f .\docker\Dockerfile .`
If your cloud uses a different CPU architecture than your development
machine (e.g., you are on a Mac M1 and your cloud provider is amd64),
you'll want to build the image for that platform, e.g.:
`docker build --platform=linux/amd64 --tag goforum -f .\docker\Dockerfile .`.

Then, push it to your registry, e.g. `docker push myregistry.com/goforum`.

Consult Docker's [getting started](https://docs.docker.com/go/get-started-sharing/)
docs for more detail on building and pushing.

### References
* [Docker's Go guide](https://docs.docker.com/language/golang/)