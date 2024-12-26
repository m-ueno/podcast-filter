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
	RuleSet        RuleSet
	OutputFileName string
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

	return ParseFilterWrite(resp.Body, file, t.RuleSet)
}
