FROM debian:latest
RUN apt-get update && apt-get -y install cron

COPY backup.sh /root/backup.sh
RUN chmod +x /root/backup.sh

COPY backup-cron /etc/cron.d/backup-cron
RUN chmod 0644 /etc/cron.d/backup-cron
RUN crontab /etc/cron.d/backup-cron

RUN touch /var/log/cron.log
CMD cron && tail -f /var/log/cron.log