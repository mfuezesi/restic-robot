restic stats --mode raw-data
    {
      "total_size":2248,
      "total_file_count":0,
      "total_blob_count":9
    }

    repository d1386767 opened successfully, password is correct
    scanning...
    Stats in raw-data mode:
    Snapshots processed:   122
      Total Blob Count:   9
            Total Size:   2.195 KiB


restic stats
    {
      "total_size":5,
      "total_file_count":230
    }

    repository d1386767 opened successfully, password is correct
    scanning...
    Stats in restore-size mode:
    Snapshots processed:   122
      Total File Count:   230
            Total Size:   5 B  


restic snapshots
  [
    {
      "time":"2021-01-11T23:56:22.330806+01:00",
      "tree":"515616f0293a9737cfa632f6f83a340a720cb9a72ad4959063ba4c3e49d2d8db",
      "paths":[
         "/Users/m/Desktop/APZ RESTIC und MINIO/restic-robot/data"
      ],
      "hostname":"m.klybck.gopf.xyz",
      "username":"m",
      "uid":501,
      "gid":20,
      "id":"f496ea67579fbea2becacdff36ef8a84c3e7bc17e94bbebc89ec060991c8f405",
      "short_id":"f496ea67"
    },
  ]

    c43d5679  2021-01-15 21:32:30  m.klybck.gopf.xyz              /Users/m/Desktop/restic-robot/data
    ---------------------------------------------------------------------------------------------------------------------
    122 snapshots

restic snapshots latest
  [
    {
      "time":"2021-01-15T21:32:30.443661+01:00",
      "parent":"9c735c865857b3bcf00792c403b5a072fc7d3b661e0c10239e257dbbfa06e4dd",
      "tree":"6cac15ce015370561eae750bad4db3724f48c928b6d48e6ef84efcf059573299",
      "paths":[
          "/Users/m/Desktop/restic-robot/data"
      ],
      "hostname":"m.klybck.gopf.xyz",
      "username":"m",
      "uid":501,
      "gid":20,
      "id":"c43d567998c6cdd5465215035aa53fe05f76589ff781752eff463b95cc047b21",
      "short_id":"c43d5679"
    }
  ]

    repository d1386767 opened successfully, password is correct
    ID        Time                 Host               Tags        Paths
    ------------------------------------------------------------------------------------------------
    c43d5679  2021-01-15 21:32:30  m.klybck.gopf.xyz              /Users/m/Desktop/restic-robot/data
    ------------------------------------------------------------------------------------------------
    1 snapshots
    
    
restic diff c43d5679 9c735c86
    repository d1386767 opened successfully, password is correct
    comparing snapshot c43d5679 to 9c735c86:


    Files:           0 new,     0 removed,     0 changed
    Dirs:            0 new,     0 removed
    Others:          0 new,     0 removed
    Data Blobs:      0 new,     0 removed
    Tree Blobs:      0 new,     0 removed
      Added:   0 B  



restic backup ./data
    {
      "message_type":"status",
      "percent_done":0,
      "total_files":1,
      "total_bytes":1073741824
    }{
      "message_type":"summary",
      "files_new":0,
      "files_changed":0,
      "files_unmodified":2,
      "dirs_new":0,
      "dirs_changed":0,
      "dirs_unmodified":1,
      "data_blobs":0,
      "tree_blobs":0,
      "data_added":0,
      "total_files_processed":2,
      "total_bytes_processed":1073741829,
      "total_duration":0.516489176,
      "snapshot_id":"533b6766"
    }

    repository d1386767 opened successfully, password is correct

    Files:           1 new,     0 changed,     1 unmodified
    Dirs:            0 new,     1 changed,     0 unmodified
    Added to the repo: 647.025 KiB

    processed 2 files, 1.000 GiB in 0:01
    snapshot 8cb0b81f saved


restic backup --quiet ./data
    {
      "message_type":"status",
      "percent_done":0,
      "total_files":1,
      "total_bytes":1073741824
    }{
      "message_type":"summary",
      "files_new":0,
      "files_changed":0,
      "files_unmodified":2,
      "dirs_new":0,
      "dirs_changed":0,
      "dirs_unmodified":1,
      "data_blobs":0,
      "tree_blobs":0,
      "data_added":0,
      "total_files_processed":2,
      "total_bytes_processed":1073741829,
      "total_duration":0.51291267,
      "snapshot_id":"403a87a8"
    }

    nada?!


json backup runtime
    {"message_type":"status","percent_done":0,"total_files":1,"total_bytes":4080}
    {"message_type":"status","percent_done":0.19733841943961014,"total_files":1033,"files_done":202,"total_bytes":4312926,"bytes_done":851106}
    ...
    {"message_type":"status","seconds_elapsed":1,"percent_done":0.2716393457395461,"total_files":14248,"files_done":3900,"total_bytes":56817984,"bytes_done":15434000}
    ...
    {"message_type":"status","action":"scan_finished","item":"","duration":3.111377127,"data_size":610431843,"metadata_size":0,"total_files":152487}
    ...
    {"message_type":"status","action":"scan_finished","item":"","duration":2.909463954,"data_size":610431843,"metadata_size":0,"total_files":152487}
    ...
    {"message_type":"status","seconds_elapsed":2,"percent_done":0.27880917411446376,"total_files":152487,"files_done":42601,"total_bytes":610431843,"bytes_done":170193998}
    ...
    {"message_type":"status","seconds_elapsed":3,"seconds_remaining":7,"percent_done":0.295835509026681,"total_files":152487,"files_done":45200,"total_bytes":610431843,"bytes_done":180587415}
    ...
    {"message_type":"summary","files_new":0,"files_changed":0,"files_unmodified":152487,"dirs_new":0,"dirs_changed":0,"dirs_unmodified":1542,"data_blobs":0,"tree_blobs":0,"data_added":0,"total_files_processed":152487,"total_bytes_processed":610431843,"total_duration":7.959521653,"snapshot_id":"e8bab5b8"}





total_files_processed":102487 (without dirs)
"total_bytes_processed":410069256

"total_size":410069256,"total_file_count":103524 (incl. dirs!)
{"total_size":410069256,"total_file_count":663089}


https://github.com/sinnwerkstatt/runrestic/blob/master/runrestic/metrics/prometheus.py
https://github.com/restic/restic/blob/master/internal/ui/json/backup.go