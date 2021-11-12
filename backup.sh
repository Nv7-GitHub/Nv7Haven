tar -zcf /backups/backup_$(date +%Y%m%d).tar.gz -C /data
find /home/tony/backup/daily/* -mtime +7 -delete