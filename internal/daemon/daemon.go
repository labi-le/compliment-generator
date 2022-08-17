package daemon

import (
	"github.com/sirupsen/logrus"
	"time"
)

var (
	cDaemon *Daemon
)

type Daemon struct {
	Tasks *Tasks
	*logrus.Logger
}

type Goal func(t *Task) error

type Task struct {
	ID          string
	Goal        Goal
	IsStarted   bool
	NeedStop    bool
	RepeatEvery time.Duration
}

type Tasks []*Task

func NewDaemon(logger *logrus.Logger) *Daemon {
	if cDaemon == nil {
		cDaemon = &Daemon{
			Tasks:  &Tasks{},
			Logger: logger,
		}
	}

	return cDaemon
}

func Add(d *Task) {
	cDaemon.Add(d)
}

func Delete(id string) {
	cDaemon.Delete(id)
}

func (d *Daemon) Add(t *Task) {
	if d.Find(t.ID) != nil {
		d.Infof("Tasks daemon user %s already exist", t.ID)
		return
	}

	*d.Tasks = append(*d.Tasks, t)
}

func (d *Daemon) Find(id string) *Task {
	for _, task := range *d.Tasks {
		if task.ID == id {
			return task
		}
	}
	return nil
}

func (d *Daemon) Delete(u string) {
	for i, task := range *d.Tasks {
		if d.Find(u) != nil {
			*d.Tasks = append((*d.Tasks)[:i], (*d.Tasks)[i+1:]...)
			task.NeedStop = true

			d.Infof("Tasks daemon delete user %s", u)
		}
	}
}

func (d *Daemon) RunDaemon() {
	d.Info("Tasks daemon started")

	go taskCounter(d)

	for {
		for _, t := range *d.Tasks {
			if t == nil {
				time.Sleep(1 * time.Second)
				continue
			}

			go task(t, d.Logger)
		}
	}
}

func taskCounter(d *Daemon) {
	for {
		d.Infof("Tasks daemon users count: %d", len(*d.Tasks))
		time.Sleep(10 * time.Second)
	}
}

func task(t *Task, l *logrus.Logger) {
	if t.IsStarted {
		return
	}

	l.Infof("Tasks daemon start: %s", t.ID)
	t.IsStarted = true

	for {
		if t.NeedStop {
			break
		}

		err := t.Goal(t)
		if err != nil {
			l.Errorf("Tasks daemon error: %s", err)
		}

		time.Sleep(t.RepeatEvery)
	}
}
