#!/bin/bash

IMAGE_NAME="gopherdb"
PORT=5321

echo "🔧 Building Docker image: $IMAGE_NAME..."

docker build -t $IMAGE_NAME .

if [ "$1" == "--run" ]; then
  echo "🧹 Removing old container if it exists..."
  docker rm -f $IMAGE_NAME 2>/dev/null

  echo "🚀 Running container..."
  docker run -d -p $PORT:$PORT --name $IMAGE_NAME $IMAGE_NAME
  echo "✅ GopherDb is running at http://localhost:$PORT"
else
  echo "✅ Docker image '$IMAGE_NAME' built successfully."
  echo "📝 To run the container: ./build.sh --run"
fi
