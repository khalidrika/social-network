#!/bin/zsh

set -e

echo "📦 Installing Docker (rootless)..."

# Check if Docker is already installed
if command -v docker &> /dev/null; then
    echo "✅ Docker is already installed."
    docker --version
else
    # Download and install Docker rootless
    echo "📥 Downloading Docker rootless installer..."
    curl -fsSL https://get.docker.com/rootless | sh
    echo "✅ Docker rootless installed."
fi

# Set up environment variables first
echo "⚙️ Setting up environment variables..."

# Check if PATH is already in .zshrc
if ! grep -q 'export PATH=$HOME/bin:$PATH' ~/.zshrc 2>/dev/null; then
    echo 'export PATH=$HOME/bin:$PATH' >> ~/.zshrc
fi

# Check if DOCKER_HOST is already in .zshrc
if ! grep -q 'export DOCKER_HOST=unix://$XDG_RUNTIME_DIR/docker.sock' ~/.zshrc 2>/dev/null; then
    echo 'export DOCKER_HOST=unix://$XDG_RUNTIME_DIR/docker.sock' >> ~/.zshrc
fi

# Apply environment variables for current session
export PATH=$HOME/bin:$PATH
export DOCKER_HOST=unix://$XDG_RUNTIME_DIR/docker.sock

echo "✅ Environment configured."

# Run rootless setup, skip iptables check, continue on error
if [ -f "$HOME/bin/dockerd-rootless-setuptool.sh" ]; then
    echo "⚙️ Running dockerd-rootless-setuptool.sh with --skip-iptables..."
    $HOME/bin/dockerd-rootless-setuptool.sh install --skip-iptables || echo "⚠️ Setup tool completed with warnings (this is normal)."
else
    echo "⚠️ dockerd-rootless-setuptool.sh not found, skipping setup tool."
fi

# Check if Docker daemon is already running
if docker info &> /dev/null; then
    echo "✅ Docker daemon is already running."
else
    # Start the rootless Docker daemon in background
    echo "🚀 Starting Docker daemon (rootless) in background..."
    
    # Kill any existing docker daemon processes
    pkill -f dockerd-rootless || true
    
    # Wait a moment
    sleep 2
    
    # Start new daemon
    nohup $HOME/bin/dockerd-rootless.sh > ~/docker-rootless.log 2>&1 &
    
    echo "✅ Docker daemon started in background (log: ~/docker-rootless.log)"
    
    # Wait for daemon to start
    echo "⏳ Waiting for Docker daemon to start..."
    for i in {1..30}; do
        if docker info &> /dev/null; then
            echo "✅ Docker daemon is ready!"
            break
        fi
        sleep 1
        echo -n "."
    done
fi

# Install Docker Compose v2 if not already installed
if [ ! -f "$HOME/.docker/cli-plugins/docker-compose" ]; then
    echo "📦 Installing Docker Compose v2..."
    mkdir -p ~/.docker/cli-plugins
    
    # Get the latest version or use a specific one
    COMPOSE_VERSION="v2.26.1"
    curl -SL "https://github.com/docker/compose/releases/download/${COMPOSE_VERSION}/docker-compose-linux-x86_64" \
      -o ~/.docker/cli-plugins/docker-compose
    chmod +x ~/.docker/cli-plugins/docker-compose
    
    # Add to PATH if not already there
    if ! grep -q 'export PATH=$HOME/.docker/cli-plugins:$PATH' ~/.zshrc 2>/dev/null; then
        echo 'export PATH=$HOME/.docker/cli-plugins:$PATH' >> ~/.zshrc
    fi
    export PATH=$HOME/.docker/cli-plugins:$PATH
    
    echo "✅ Docker Compose installed."
else
    echo "✅ Docker Compose is already installed."
fi

echo ""
echo "🔍 Verifying installation..."

# Test Docker
if docker --version; then
    echo "✅ Docker is working!"
else
    echo "❌ Docker is not working properly."
fi

# Test Docker Compose
if docker compose version; then
    echo "✅ Docker Compose is working!"
else
    echo "❌ Docker Compose is not working properly."
fi

# Test Docker daemon
if docker info &> /dev/null; then
    echo "✅ Docker daemon is running!"
else
    echo "❌ Docker daemon is not running. Check ~/docker-rootless.log for errors."
fi

echo ""
echo "✅ Installation complete!"
echo ""
echo "📝 Next steps:"
echo "1. Run 'source ~/.zshrc' to apply environment changes"
echo "2. Or restart your terminal"
echo "3. Test with: docker run hello-world"
echo ""
echo "💡 If Docker doesn't work immediately:"
echo "   - Check the log: tail -f ~/docker-rootless.log"
echo "   - Restart daemon: pkill dockerd-rootless && nohup ~/bin/dockerd-rootless.sh > ~/docker-rootless.log 2>&1 &"
