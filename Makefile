all:
	CGO_ENABLED=0 go build -v 
ball:
	cd frontend && npm run build && cd .. && CGO_ENABLED=0 go build -v