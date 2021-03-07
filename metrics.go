package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func (b *backup) initMetrics() {
	b.runsSuccessful = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "restic_runs_successful",
		Help: "The total number of backups that succeeded during runtime.",
	})
	b.runsFailed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "restic_runs_failed",
		Help: "The total number of backups that failed during runtime.",
	})
	b.runsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "restic_runs_total",
		Help: "The total number of backups attempted during runtime, including failures.",
	})

	b.runtimeActive = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "restic_runtime_active",
		Help: "Shows if job is running.",
	})
	b.runtimePercentDone = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "restic_runtime_percent_done",
		Help: "Percent done for ongoing backup job.",
	})
	b.runtimeElapsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "restic_runtime_elapsed",
		Help: "Seconds running.",
	})
	b.runtimeRemaining = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "restic_runtime_remaining",
		Help: "Seconds remaining.",
	})

	// b.backupDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
	// 	Name: "backup_run_duration",
	// 	Help: "The duration of backups in milliseconds.",
	// })

	// b.filesTotal = prometheus.NewGauge(prometheus.GaugeOpts{
	// 	Name: "backup_files_total",
	// 	Help: "Amount of total files.",
	// })
	// b.filesNew = prometheus.NewHistogram(prometheus.HistogramOpts{
	// 	Name: "backup_files_new",
	// 	Help: "Amount of new files.",
	// })
	// b.filesChanged = prometheus.NewHistogram(prometheus.HistogramOpts{
	// 	Name: "backup_files_changed",
	// 	Help: "Amount of files with changes.",
	// })
	// b.filesUnmodified = prometheus.NewHistogram(prometheus.HistogramOpts{
	// 	Name: "backup_files_unmodified",
	// 	Help: "Amount of files unmodified since last backup.",
	// })
	// b.filesProcessed = prometheus.NewHistogram(prometheus.HistogramOpts{
	// 	Name: "backup_files_processed",
	// 	Help: "Total number of files scanned by the backup for changes.",
	// })
	// b.bytesAdded = prometheus.NewHistogram(prometheus.HistogramOpts{
	// 	Name: "backup_added_bytes",
	// 	Help: "Total number of bytes added to the repository.",
	// })
	// b.bytesProcessed = prometheus.NewHistogram(prometheus.HistogramOpts{
	// 	Name: "backup_processed_bytes",
	// 	Help: "Total number of bytes scanned by the backup for changes",
	// })

	// b.secondsElapsed = prometheus.NewGauge(prometheus.GaugeOpts{
	// 	Name: "backup_run_seconds_elapsed",
	// 	Help: "Seconds elapsed since last backup start.",
	// })

	// b.percentDone = prometheus.NewGauge(prometheus.GaugeOpts{
	// 	Name: "backup_run_percent_done",
	// 	Help: "Percent done for ongoing backup job.",
	// })

	prometheus.MustRegister(
		b.runsSuccessful,
		b.runsFailed,
		b.runsTotal,

		b.runtimeActive,
		b.runtimePercentDone,
		b.runtimeElapsed,
		b.runtimeRemaining,

		// b.preHookDuration // not yet
		// b.postHookDuration // not yet

		// b.backupFilesNew,
		// b.backupFilesChanged,
		// b.backupFilesUnmodified,
		// b.backupDirsNew,
		// b.backupDirsChanged,
		// b.backupDirsUnmodified,
		// b.backupDataBlobs,
		// b.backupTreeBlobs,
		// b.backupDataAdded,
		// b.backupTotalFilesProcessed,
		// b.backupTotalBytesProcessed,
		// b.backupTotalDuration,
		// b.backupShortId,
		// -> timestamp?

		// b.statsDuration, // not yet
		// b.statsSize, // not yet
		// b.statsSizeRaw, // not yet
		// b.statsFilesTotal, // not yet
		// b.statsSnapshotsTotal, // not yet

		// b.snapshotsDuration, // not yet
		// b.snapshotTimestamp, // not yet
		// b.snapshotHostname, // not yet
		// b.snapshotPath, // not yet
		// b.snapshotShortId, // not yet

		// b.diffFilesAdded, // not yet
		// b.difFilesChanged, // not yet
		// b.diffFilesRemoved, // not yet
		// b.diffDirsAdded, // not yet
		// b.diffDirsChanged, // not existend? // not yet
		// b.diffDirsRemoved, // not yet

		// prune + remove? // not yet
	)
}

func (b *backup) startMetricsServer() {
	http.Handle(b.PrometheusEndpoint, promhttp.Handler())
	err := http.ListenAndServe(b.PrometheusAddress, nil)
	logger.Fatal("metrics server closed", zap.Error(err))
}
