deploy:
	GOOS=linux go build -o /tmp/edbd ./cmd/edbd
	gzip /tmp/edbd
	scp /tmp/edbd.gz root@edb.dev:/usr/local/bin/edbd.gz
	ssh -T root@edb.dev "service edbd stop && gunzip -f /usr/local/bin/edbd.gz && service edbd start"
