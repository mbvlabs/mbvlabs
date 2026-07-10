image := "ghcr.io/mbvlabs/mbvlabs:latest"

# Build and publish the latest Docker image to GHCR.
docker-publish:
    docker build --tag {{ image }} .
    docker push {{ image }}
