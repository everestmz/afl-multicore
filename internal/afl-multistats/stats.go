package multistats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type AFLStats struct {
	StartTime           int     `json:"start_time"`
	LastUpdate          int     `json:"last_update"`
	FuzzerPid           int     `json:"fuzzer_pid"`
	CyclesDone          int     `json:"cycles_done"`
	ExecsDone           int     `json:"execs_done"`
	ExecsPerSecond      float64 `json:"execs_per_sec"`
	PathsTotal          int     `json:"paths_total"`
	PathsFavored        int     `json:"paths_favored"`
	PathsFound          int     `json:"paths_found"`
	PathsImported       int     `json:"paths_imported"`
	MaxDepth            int     `json:"max_depth"`
	CurPath             int     `json:"cur_path"`
	PendingFavs         int     `json:"pending_favs"`
	PendingTotal        int     `json:"pending_total"`
	VariablePaths       int     `json:"variable_paths"`
	Stability           float64 `json:"stability"`
	BitmapCvg           float64 `json:"bitmap_cvg"`
	UniqueCrashes       int     `json:"unique_crashes"`
	UniqueHangs         int     `json:"unique_hangs"`
	LastPath            int     `json:"last_path"`
	LastCrash           int     `json:"last_crash"`
	LastHang            int     `json:"last_hang"`
	ExecsSinceLastCrash int     `json:"execs_since_crash"`
	ExecTimeout         int     `json:"exec_timeout"`
	AFLBanner           string  `json:"afl_banner"`
	AFLVersion          string  `json:"afl_version"`
	TargetMode          string  `json:"target_mode"`
	CommandLine         string  `json:"command_line"`
}

func (s *AFLStats) JSON() ([]byte, error) {
	return json.Marshal(s)
}

func (s *AFLStats) Basic() *basicStats {
	return &basicStats{
		ExecsDone:      s.ExecsDone,
		ExecsPerSecond: s.ExecsPerSecond,
		PathsTotal:     s.PathsTotal,
		PathsFavored:   s.PathsFavored,
		PendingFavs:    s.PendingFavs,
		PendingTotal:   s.PendingTotal,
		UniqueCrashes:  s.UniqueCrashes,
		UniqueHangs:    s.UniqueHangs,
	}
}

type basicStats struct {
	ExecsDone      int     `json:"execs_done"`
	ExecsPerSecond float64 `json:"execs_per_sec"`
	PathsTotal     int     `json:"paths_total"`
	PathsFavored   int     `json:"paths_favored"`
	PendingFavs    int     `json:"pending_favs"`
	PendingTotal   int     `json:"pending_total"`
	UniqueCrashes  int     `json:"unique_crashes"`
	UniqueHangs    int     `json:"unique_hangs"`
}

func (s *basicStats) JSON() ([]byte, error) {
	return json.Marshal(s)
}

func (s *basicStats) Human() []byte {
	var buf bytes.Buffer
	buf.WriteString("Summary stats\n")
	buf.WriteString("=============\n")
	buf.WriteString(fmt.Sprintf("          Total execs : %v\n", s.ExecsDone))
	buf.WriteString(fmt.Sprintf("     Execs per second : %.2f\n", s.ExecsPerSecond))
	buf.WriteString(fmt.Sprintf("          Total paths : %v\n", s.PathsTotal))
	buf.WriteString(fmt.Sprintf("       Favoured paths : %v\n", s.PathsFavored))
	buf.WriteString(fmt.Sprintf("Pending favored paths : %v\n", s.PendingFavs))
	buf.WriteString(fmt.Sprintf("  Total pending paths : %v\n", s.PendingTotal))
	buf.WriteString(fmt.Sprintf("              Crashes : %v\n", s.UniqueCrashes))
	buf.WriteString(fmt.Sprintf("                Hangs : %v\n", s.UniqueHangs))
	return buf.Bytes()
}

func ReadStats(syncDir string) ([]*AFLStats, error) {
	dirs, err := ioutil.ReadDir(syncDir)
	if err != nil {
		return nil, err
	}

	statsList := []*AFLStats{}

	for _, dir := range dirs {
		stats := &AFLStats{}
		statsFile := filepath.Join(filepath.Join(syncDir, dir.Name(), "fuzzer_stats"))
		statsBytes, err := ioutil.ReadFile(statsFile)
		if err != nil {
			return nil, err
		}

		// Get all the stats into a map
		statsMap := map[string]string{}
		statsLines := strings.Split(string(statsBytes), "\n")
		for _, line := range statsLines {
			lineSplit := strings.Split(line, ":")
			if len(lineSplit) != 2 {
				continue
			}
			k := strings.TrimSpace(lineSplit[0])
			v := strings.TrimSpace(lineSplit[1])
			statsMap[k] = v
		}

		// Map names to indexes of the struct
		tagMap := map[string]int{}
		t := reflect.Indirect(reflect.ValueOf(stats)).Type()
		for i := 0; i < t.NumField(); i++ {
			tag := t.Field(i).Tag.Get("json")
			tagMap[tag] = i
		}

		v := reflect.ValueOf(stats)
		s := v.Elem()
		for k, v := range statsMap {
			if idx, ok := tagMap[k]; ok {
				switch s.Field(idx).Kind() {
				case reflect.String:
					s.Field(idx).SetString(v)
				case reflect.Int:
					val, err := strconv.ParseInt(v, 10, 0)
					if err != nil {
						return nil, err
					}
					s.Field(idx).SetInt(val)
				case reflect.Float64:
					val, err := strconv.ParseFloat(
						strings.ReplaceAll(v, "%", ""), 64)
					if err != nil {
						return nil, err
					}
					s.Field(idx).SetFloat(val)
				}
			}
		}
		statsList = append(statsList, stats)
	}

	return statsList, nil
}

func MergeStats(stats []*AFLStats) *AFLStats {
	finalStats := *(stats[0])
	finalStats.AFLBanner = strings.Split(finalStats.AFLBanner, "_")[0]

	for _, s := range stats[1:] {
		if s.StartTime < finalStats.StartTime {
			finalStats.StartTime = s.StartTime
		}
		if s.LastUpdate > finalStats.LastUpdate {
			finalStats.LastUpdate = s.LastUpdate
		}
		finalStats.CyclesDone += s.CyclesDone
		finalStats.ExecsDone += s.ExecsDone
		finalStats.ExecsPerSecond += s.ExecsPerSecond
		if s.PathsTotal > finalStats.PathsTotal {
			finalStats.PathsTotal = s.PathsTotal
		}
		if s.PathsFavored > finalStats.PathsFavored {
			finalStats.PathsFavored = s.PathsFavored
		}
		if s.PathsFound > finalStats.PathsFound {
			finalStats.PathsFound = s.PathsFound
		}
		// Leave paths imported alone
		if s.MaxDepth > finalStats.MaxDepth {
			finalStats.MaxDepth = s.MaxDepth
		}
		// Leave cur_path, variable_paths alone
		finalStats.PendingFavs += s.PendingFavs
		finalStats.PendingTotal += s.PendingTotal
		if s.Stability < finalStats.Stability {
			finalStats.Stability = s.Stability
		}
		if s.BitmapCvg > finalStats.BitmapCvg {
			finalStats.BitmapCvg = s.BitmapCvg
		}
		finalStats.UniqueCrashes += s.UniqueCrashes
		finalStats.UniqueHangs += s.UniqueHangs
		if s.LastPath > finalStats.LastPath {
			finalStats.LastPath = s.LastPath
		}
		if s.LastCrash > finalStats.LastCrash {
			finalStats.LastCrash = s.LastCrash
		}
		if s.LastHang > finalStats.LastHang {
			finalStats.LastHang = s.LastHang
		}
		finalStats.ExecsSinceLastCrash += s.ExecsSinceLastCrash
	}

	return &finalStats
}
