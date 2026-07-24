FROM golang:1.25-bookworm

RUN apt-get update && apt-get install -y \
    gcc \
    pkg-config \
    libgl1-mesa-dev \
    libx11-dev \
    libxcursor-dev \
    libxrandr-dev \
    libxinerama-dev \
    libxi-dev \
    libxxf86vm-dev \
    libxkbcommon-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

CMD ["go", "run", "."]