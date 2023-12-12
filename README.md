# tic-tac-toe

## What am I looking at?
This application is my attempt at writing a worker-based gameserver that utilizes Go's asynchronicity model. As a POC, I wrote a simple tic-tac-toe implementation, which can be interacted with using the custom websocket-based RPC or the [frontend app](https://github.com/voxelost/tic-tac-toe-frontend/).

## How can I run it?
To run the server you need to have Go or Docker (or Podman) installed.

The Go way:
```sh
go run .
```

The Docker way:
```sh
docker compose up --build
```

OR
```sh
docker build --tag local/tic-tac-toe-backend .
docker run -p 8000:8000 local/tic-tac-toe-backend
```

