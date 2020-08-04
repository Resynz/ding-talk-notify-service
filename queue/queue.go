/**
 * @Author: Resynz
 * @Date: 2020/8/4 10:36
 */
package queue

type NotifyTask struct {
}

var (
	notifyQueue    chan *NotifyTask
	notifyExitChan chan bool
	stopQueue      = false
)
