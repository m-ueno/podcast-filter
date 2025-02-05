package filter

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

type FilterFeedTaskSet struct {
	Tasks     []FilterFeedTask
	OutputDir string
}

func NewFilterFeedTaskSet(tasks []FilterFeedTask, outputDir string) FilterFeedTaskSet {
	return FilterFeedTaskSet{
		Tasks:     tasks,
		OutputDir: outputDir,
	}
}

func (ts *FilterFeedTaskSet) AddTask(task FilterFeedTask) {
	ts.Tasks = append(ts.Tasks, task)
}

func (ts FilterFeedTaskSet) Run() error {
	if err := os.MkdirAll(ts.OutputDir, 0755); err != nil && err != os.ErrExist {
		return err
	}
	var g errgroup.Group
	for _, task := range ts.Tasks {
		task := task
		g.Go(func() error {
			defer slog.Info("End task", "url", task.FeedUrl)
			slog.Info("Start task", "url", task.FeedUrl)
			return task.Run(ts.OutputDir)
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

type FilterFeedTask struct {
	FeedUrl        string
	Rule           Rule
	OutputFileName string
}

func NewFilterFeedTask(url string, rule Rule, filename string) FilterFeedTask {
	return FilterFeedTask{
		FeedUrl:        url,
		Rule:           rule,
		OutputFileName: filename,
	}
}

func (t FilterFeedTask) Run(basedir string) error {
	resp, err := http.Get(t.FeedUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.OpenFile(filepath.Join(basedir, t.OutputFileName), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return ParseFilterWrite(resp.Body, file, t.Rule)
}
