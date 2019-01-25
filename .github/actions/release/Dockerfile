FROM alpine:latest

LABEL "com.github.actions.name"="guitarrapc"
LABEL "com.github.actions.description"="add release package on the github"
LABEL "com.github.actions.icon"="activity"
LABEL "com.github.actions.color"="yellow"

RUN apk add --no-cache \
	bash \
	ca-certificates \
	coreutils \
	curl \
	jq

COPY entrypoint.sh /usr/local/bin/entrypoint

CMD ["entrypoint"]