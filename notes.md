https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html

Needed.

1. total backup size
2. total backup size raw
3. total file count
    total snapshots count
4. latest snapshot short id
5. latest snapshot timestamp
6. latest snapshot duration
7. latest snapshot files processed
7. files added
8. files removed 

Pattern design
runtime data

backup data

stats data

prune data


    restic_stats_total_size
    restic_stats_total_size_raw
    restic_stats_total_file_count

restic_metadata{short_id=} 1

restic_timestamp
- restic_duration
- restic_files_processed
- restic_bytes_processed

- restic_added
    restic_removed

- restic_files_new
- restic_files_changed
- restic_files_unmodified  
   restic_files_removed
restic_dirs_new
restic_dirs_changed
restic_dirs_unmodified
    restic_dirs_removed

restic_running %
restic_scanner ?
