package filter

import (
	"net/http"
	"os"
	"path/filepath"
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
	for _, task := range ts.Tasks {
		if err := task.Run(ts.OutputDir); err != nil {
			return err
		}
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
