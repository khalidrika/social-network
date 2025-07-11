#!/bin/bash

echo "🐳 Building and starting Social Network application..."

# Stop any existing containers
echo "Stopping existing containers..."
docker-compose down

# Build and start containers
echo "Building and starting containers..."
docker-compose up --build -d

# Wait a moment for containers to start
echo "Waiting for containers to start..."
sleep 5

# Show container status
echo "Container status:"
docker ps -a

echo ""
echo "✅ Application is starting!"
echo "🌐 Frontend: http://localhost"
echo "🔧 Backend API: http://localhost:8081"
echo ""
echo "To view logs:"
echo "  Frontend: docker logs social-network-frontend"
echo "  Backend:  docker logs social-network-backend"
echo ""
echo "To stop: docker-compose down"
