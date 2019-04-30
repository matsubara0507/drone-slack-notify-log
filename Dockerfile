FROM plugins/base:multiarch

ADD release/linux/amd64/drone-slack-notify-log /bin/
ENTRYPOINT ["/bin/drone-slack-notify-log"]
