echo "Backing up..."
tar -zcvf /backups/backup_$(date +%Y%m%d).tar.gz /data
find /backups/* -mtime +7 -delete