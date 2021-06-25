# Run Tests For Admin User
docker run \
    -v $(pwd)/tests:/etc/newman --network=water-api_default \
    --entrypoint /bin/bash postman/newman:ubuntu \
    -c "newman run /etc/newman/water-regression.postman_collection.json \
        --environment=/etc/newman/water-docker-compose.postman_environment.json"