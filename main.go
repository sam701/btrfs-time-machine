package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"sort"
	"time"
)

var config *Config
var today = time.Now().Truncate(24 * time.Hour)

func main() {
	configPath := flag.String("c", "config.yaml", "Path to the config.yaml")
	flag.Parse()

	if *configPath == "" {
		flag.Usage()
		return
	}

	config = getConfig(*configPath)

	createSnapshot()

	toDelete := getDatesToDelete(getSnapshotDates())
	for _, date := range toDelete {
		deleteSnapshot(date)
	}
}

func createSnapshot() {
	dest := snapshotPath(time.Now())
	if _, err := os.Stat(dest); err == nil {
		log.Println(dest, "already exists")
		return
	}

	execCommand(exec.Command("btrfs", "subvolume", "snapshot", "-r", config.ContentDir, dest))
	fmt.Println("Created snapshot", dest)
}

func snapshotPath(date time.Time) string {
	return path.Join(config.SnapshotsDir, date.Format("2006-01-02"))
}

func deleteSnapshot(date time.Time) {
	sp := snapshotPath(date)
	execCommand(exec.Command("btrfs", "subvolume", "delete", sp))
	fmt.Println("Delete snapshot", sp)
}

func getSnapshotDates() []time.Time {
	f, err := os.Open(config.SnapshotsDir)
	if err != nil {
		log.Fatalln("ERROR: cannot open file:", err)
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(names)))

	out := []time.Time{}
	for _, n := range names {
		t, err := time.Parse("2006-01-02", n)
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		out = append(out, t)
	}
	return out
}

func execCommand(cmd *exec.Cmd) {
	var buf bytes.Buffer
	cmd.Stderr = &buf
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		os.Stdout.Write(buf.Bytes())
		log.Fatalln("ERROR", err)
	}
}

func getSlotNumber(date time.Time) int {
	diff := int(today.Sub(date).Hours()) / 24

	acc := 0
	for _, x := range config.Retentions {
		no := diff / x.EveryNDays
		if no <= x.SnapshotsToKeep {
			return acc + no
		}

		acc += x.SnapshotsToKeep
		diff -= x.EveryNDays * x.SnapshotsToKeep
	}
	return -1
}

func getDatesToDelete(dates []time.Time) []time.Time {
	prevSlot := -1
	var prevDate time.Time
	out := []time.Time{}
	for _, date := range dates {
		slot := getSlotNumber(date)
		if slot == prevSlot {
			out = append(out, prevDate)
		}

		prevSlot = slot
		prevDate = date
	}

	return out
}
