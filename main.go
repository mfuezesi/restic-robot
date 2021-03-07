package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron"
	"go.uber.org/zap"
)

type backup struct {
	Schedule           string `required:"true"    envconfig:"SCHEDULE"`            // cron schedule
	Repository         string `required:"true"    envconfig:"RESTIC_REPOSITORY"`   // repository name
	Password           string `required:"true"    envconfig:"RESTIC_PASSWORD"`     // repository password
	Args               string `                   envconfig:"RESTIC_ARGS"`         // additional args for backup command
	RunOnBoot          bool   `                   envconfig:"RUN_ON_BOOT"`         // run a backup on startup
	PrometheusEndpoint string `default:"/metrics" envconfig:"PROMETHEUS_ENDPOINT"` // metrics endpoint
	PrometheusAddress  string `default:":8080"    envconfig:"PROMETHEUS_ADDRESS"`  // metrics host:port
	PreCommand         string `                   envconfig:"PRE_COMMAND"`         // command to execute before restic is executed
	PostCommand        string `                   envconfig:"POST_COMMAND"`        // command to execute after restic was executed (successfully)

	runsSuccessful     prometheus.Counter
	runsFailed         prometheus.Counter
	runsTotal          prometheus.Counter
	runtimeActive      prometheus.Gauge
	runtimePercentDone prometheus.Gauge
	runtimeElapsed     prometheus.Gauge
	runtimeRemaining   prometheus.Gauge
	// backupDuration        prometheus.Histogram
	// filesNew              prometheus.Histogram
	// filesChanged          prometheus.Histogram
	// filesUnmodified       prometheus.Histogram
	// filesProcessed        prometheus.Histogram
	// bytesAdded            prometheus.Histogram
	// bytesProcessed        prometheus.Histogram
	// secondsElapsed        prometheus.Gauge
	// percentDone           prometheus.Gauge
	// filesTotal            prometheus.Gauge
}

func main() {
	b := backup{}
	err := envconfig.Process("", &b)
	if err != nil {
		logger.Fatal("failed to configure", zap.Error(err))
	}

	err = b.Ensure()
	if err != nil {
		logger.Fatal("failed to ensure repository", zap.Error(err))
	}

	b.initMetrics()
	go b.startMetricsServer()

	cr := cron.New()
	err = cr.AddJob(b.Schedule, &b)
	if err != nil {
		logger.Fatal("failed to schedule task", zap.Error(err))
	}
	logger.Info("set cron job",
		zap.String("schedule", b.Schedule),
	)

	if b.RunOnBoot {
		b.Run()
	}
	cr.Run()
}

var ob struct {
	MessageType      string  `json:"message_type"`
	SecondsElapsed   int     `json:"seconds_elapsed"`
	SecondsRemaining int     `json:"seconds_remaining"`
	PercentDone      float64 `json:"percent_done"`
	TotalFiles       int     `json:"total_files"`
	FilesDone        int     `json:"files_done"`
	TotalBytes       int     `json:"total_bytes"`
	BytesDone        int     `json:"bytes_done"`
	DataAdded        int     `json:"data_added"`
	SnapshotID       string  `json:"snapshot_id"`
}

// Run performs the backup
func (b *backup) Run() {
	logger.Info("backup started")
	startTime := time.Now()

	b.runtimeActive.Set(1)

	if len(b.PreCommand) < 0 {
		if stdout, err := b.executePreCommand(); err != nil {
			logger.Error("failed to execute pre-command: " + err.Error())
			b.runsFailed.Inc()
			b.runsTotal.Inc()
			return
		} else {
			logger.Info("output of pre-command: " + *stdout)
		}
	}

	argsx := strings.Fields("restic backup -v -q --json " + b.Args)
	cmd := exec.Command(argsx[0], argsx[1:]...)
	// errbuf := bytes.NewBuffer(nil)
	// outbuf := bytes.NewBuffer(nil)
	// cmd.Stderr = errbuf
	// cmd.Stdout = outbuf

	// if err := cmd.Run(); err != nil {
	// 	logger.Error("failed to run backup",
	// 		zap.Error(err),
	// 		zap.String("output", errbuf.String()))
	// 	b.backupsRunsFailed.Inc()
	// 	b.backupsRunsTotal.Inc()
	// 	return
	// }

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	// s := bufio.NewScanner(stdout)
	// for s.Scan() {
	// 	fmt.Println(s.Text())
	// }

	dec := json.NewDecoder(stdout)

	for {
		err := dec.Decode(&ob)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error decoding: %v", err)
		}

		b.runtimePercentDone.Set(float64(ob.PercentDone))
		b.runtimeElapsed.Set(float64(ob.SecondsElapsed))
		b.runtimeRemaining.Set(float64(ob.SecondsRemaining))

		// fmt.Printf("%v / %v / %v / %v\n", ob.MessageType, ob.PercentDone, ob.SecondsElapsed, ob.TotalFiles)
		// b.filesTotal.Set(float64(ob.TotalFiles))
		// b.percentDone.Set(float64(ob.PercentDone))
		// b.secondsElapsed.Set(float64(ob.SecondsElapsed))
		// b.filesTotal.Observe(float64(ob.TotalFiles))
		// b.filesTotal.Collect()
		// testvar := reflect.TypeOf(ob.TotalFiles)
		// fmt.Println(testvar)
		// fmt.Printf(reflect.TypeOf(ob.TotalFiles))
		// fmt.Printf(ob)
		// b.percentDone.Set(1337)
		// b.backupDuration.Observe(0)
		// b.filesProcessed.Observe(float64(statistics.filesProcessed))
		// result.filesNew, _ = strconv.Atoi(fileStats[0][1])        //nolint:errcheck

		// totalFilesMod, _ := strconv.Atoi(ob.TotalFiles)
		// testvar = reflect.TypeOf(totalFilesMod)
		// fmt.Println(testvar)

		if ob.MessageType == "summary" {
			// fmt.Printf("Finished!\n")
			// fmt.Printf("%v / %v \n", ob.SnapshotID, ob.DataAdded)
			b.runtimePercentDone.Set(1)
		}
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	if len(b.PostCommand) > 0 {
		if stdout, err := b.executePostCommand(); err != nil {
			logger.Error("failed to execute post-command: " + err.Error())
			b.runsFailed.Inc()
			b.runsTotal.Inc()
			return
		} else {
			logger.Info("output of post-command: " + *stdout)
		}
	}

	d := time.Since(startTime)
	fmt.Println(d)

	// statistics, err := extractStats(outbuf.String())
	// if err != nil {
	// 	logger.Warn("failed to extract statistics from command output",
	// 		zap.Error(err))
	// }

	logger.Info("backup completed") // zap.Duration("duration", d),
	// zap.Int("filesNew", statistics.filesNew),
	// zap.Int("filesChanged", statistics.filesChanged),
	// zap.Int("filesUnmodified", statistics.filesUnmodified),
	// zap.Int("filesProcessed", statistics.filesProcessed),
	// zap.Int("bytesAdded", statistics.bytesAdded),
	// zap.Int("bytesProcessed", statistics.bytesProcessed),
	// zap.Int("bytesProcessed", 0),

	// b.backupDuration.Observe(float64(d.Nanoseconds() * 1000))
	// b.secondsElapsed.Set(0)

	// b.SecondsElapsed.Observe(float64(0))
	// b.backupDuration.Observe(float64(d.Nanoseconds() * 1000))
	// b.filesNew.Observe(float64(123))
	// b.filesChanged.Observe(float64(statistics.filesChanged))
	// b.filesUnmodified.Observe(float64(statistics.filesUnmodified))
	// b.filesProcessed.Observe(float64(statistics.filesProcessed))
	// b.bytesAdded.Observe(float64(statistics.bytesAdded))
	// b.bytesProcessed.Observe(float64(statistics.bytesProcessed))

	// b.percentDone.Observe(float64(0))

	b.runtimePercentDone.Set(0)
	b.runtimeElapsed.Set(0)
	b.runtimeRemaining.Set(0)
	b.runtimeActive.Set(0)
	b.runsSuccessful.Inc()
	b.runsTotal.Inc()
}

type stats struct {
	runsSuccessful int
	runsFailed     int
	runsTotal      int
	// filesNew        int
	// filesChanged    int
	// filesUnmodified int
	// filesProcessed  int
	// bytesAdded      int
	// bytesProcessed  int
	// secondsElapsed  int
	// percentDone     int
	// filesTotal      int
}

var (
	matchExists = regexp.MustCompile(`.*already (exists|initialized).*`)
	// matchFileStats  = regexp.MustCompile(`Files:\s*([0-9.]*) new,\s*([0-9.]*) changed,\s*([0-9.]*) unmodified`)
	// matchAddedBytes = regexp.MustCompile(`Added to the repo: ([0-9.]+) (\w+)`)
	// matchProcessed  = regexp.MustCompile(`processed ([0-9.]*) files, ([0-9.]+) (\w+)`)
)

// func extractStats(s string) (result stats, err error) {
// 	fileStats := matchFileStats.FindAllStringSubmatch(s, -1)
// 	if len(fileStats[0]) != 4 {
// 		err = errors.Errorf("matchFileStats expected 4, got %d", len(fileStats[0]))
// 		return
// 	}
// 	result.filesNew, _ = strconv.Atoi(fileStats[0][1])        //nolint:errcheck
// 	result.filesChanged, _ = strconv.Atoi(fileStats[0][2])    //nolint:errcheck
// 	result.filesUnmodified, _ = strconv.Atoi(fileStats[0][3]) //nolint:errcheck

// 	addedBytes := matchAddedBytes.FindAllStringSubmatch(s, -1)
// 	if len(addedBytes[0]) != 3 {
// 		err = errors.Errorf("matchAddedBytes expected 3, got %d", len(addedBytes[0]))
// 		return
// 	}
// 	amount, _ := strconv.ParseFloat(addedBytes[0][1], 64) //nolint:errcheck
// 	// restic doesn't use a comma to denote thousands
// 	amount *= 1000
// 	result.bytesAdded = convert(int(amount), addedBytes[0][2])

// 	filesProcessed := matchProcessed.FindAllStringSubmatch(s, -1)
// 	if len(filesProcessed[0]) != 4 {
// 		err = errors.Errorf("filesProcessed expected 4, got %d", len(filesProcessed[0]))
// 		return
// 	}
// 	result.filesProcessed, _ = strconv.Atoi(filesProcessed[0][1]) //nolint:errcheck
// 	amount, _ = strconv.ParseFloat(filesProcessed[0][2], 64)      //nolint:errcheck
// 	amount *= 1000
// 	result.bytesProcessed = convert(int(amount), filesProcessed[0][3])

// 	return
// }

// func convert(b int, unit string) (result int) {
// 	switch unit {
// 	case "TiB":
// 		result = b * (1 << 40)
// 	case "GiB":
// 		result = b * (1 << 30)
// 	case "MiB":
// 		result = b * (1 << 20)
// 	case "KiB":
// 		result = b * (1 << 10)
// 	}
// 	return
// }

// Ensure will create a repository if it does not already exist
func (b *backup) Ensure() (err error) {
	logger.Info("ensuring backup repository exists")
	cmd := exec.Command("restic", "init")
	out := bytes.NewBuffer(nil)
	cmd.Stderr = out
	err = cmd.Run()
	if err != nil {
		if matchExists.MatchString(strings.Trim(out.String(), " \n\r")) {
			logger.Info("repository exists")
			return nil
		}
		return errors.Wrap(err, out.String())
	}
	logger.Info("successfully created repository")
	return
}
