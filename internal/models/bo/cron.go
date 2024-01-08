package bo

type CronTask interface {
	Schedule() string
	Run()
}
