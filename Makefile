api:
	CGO_ENABLED=0 go build -v 
full:
	cd frontend && npm run build && cd .. && CGO_ENABLED=0 go build -v
