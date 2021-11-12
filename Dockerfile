FROM debian:latest
ADD backup.sh /root/backup.sh
RUN echo "15 0 * * * sh /root/backup.sh" > /root/crontabfile
RUN crontab /root/crontabfile
RUN systemctl restart crond