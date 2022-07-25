package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	return output
}

func saveUserInfo(user User) {
	fmt.Printf("WRITING FILE FOR UID %d\n", user.id)

	filename := fmt.Sprintf("users/uid%d.txt", user.id)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(user.getActivityInfo())
	time.Sleep(time.Second)
}

func generateUsers(count int) []User {
	users := make([]User, count)

	for i := 0; i < count; i++ {
		users[i] = User{
			id:    i + 1,
			email: fmt.Sprintf("user%d@company.com", i+1),
			logs:  generateLogs(rand.Intn(1000)),
		}
		fmt.Printf("generated user %d\n", i+1)
		time.Sleep(time.Millisecond * 100)
	}

	return users
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}

	return logs
}

//___________________________________________

func newUser(id int) User {

	u := User{
		id:    id,
		email: fmt.Sprintf("user%d@company.com", id+1),
		logs:  generateLogs(rand.Intn(1000)),
	}
	fmt.Printf("generated user %d\n", id+1)
	time.Sleep(time.Millisecond * 100)

	return u
}

var wg *sync.WaitGroup

func main() {
	const userCount = 100
	const workerCount = 20

	jobs := make(chan int, userCount)

	rand.Seed(time.Now().Unix())

	startTime := time.Now()

	users := generateUsers(userCount)

	for _, user := range users {
		saveUserInfo(user)
	}
	t0 := time.Since(startTime).Seconds()
	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", t0)

	//-----------------------------------------------------------------

	fmt.Printf("\n\nWorker Pool Start!\n")

	startTime = time.Now()

	wg = &sync.WaitGroup{}

	for i := 0; i < userCount; i++ {
		wg.Add(1)
		jobs <- i + 1 // передаём id пользователя в канал
		if i < workerCount {
			go worker(i+1, jobs)
		}
	}
	close(jobs)

	wg.Wait()
	fmt.Printf("Worker Pool DONE!\n")

	fmt.Printf("First  process Time Elapsed: %.2f seconds\n", t0)
	fmt.Printf("Second process Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
	fmt.Printf("At %d workers, the second process is faster than first in %.2f times", workerCount, t0/time.Since(startTime).Seconds())
}

func worker(wid int, job <-chan int) {

	for id := range job {

		fmt.Printf("Worker #%d start: User #%d\n", wid, id)
		user := newUser(id)
		saveUserInfo(user)

		fmt.Printf("Worker #%d done: User #%d\n", wid, id)
		wg.Done()
	}
}
