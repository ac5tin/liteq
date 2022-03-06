package main

import (
	"liteq/lib"
	"liteq/queue"
	"log"
	"sync"
	"time"
)

func main() {
	log.Println("Initialising new client")
	client, err := lib.NewClient("0.0.0.0:8080")
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Finished initialising client")
	defer client.Close()

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// receive tasks
	go func() {
		defer wg.Done()
		ch, err := client.GetTasks(queue.TaskStatusCreated)
		if err != nil {
			log.Println(err.Error())
		}

		tasks := new([]*queue.Task)
		for {
			t := <-ch
			*tasks = append(*tasks, t)
			log.Printf("Received task %s\n", t.ID)
			if len(*tasks) == 2 {
				break
			}
		}
	}()

	// add tasks
	go func() {
		defer wg.Done()
		values := []queue.Task{
			{
				ID:           "1",
				Data:         []byte("test"),
				Status:       queue.TaskStatusCreated,
				CreationDate: time.Now(),
			},
			{
				ID:           "2",
				Data:         []byte("test2"),
				Status:       queue.TaskStatusCreated,
				CreationDate: time.Now(),
			},
		}

		log.Println("Adding tasks")
		for _, v := range values {
			value := v
			if err := client.AddTask(&value); err != nil {
				log.Fatal(err.Error())
			}
		}
	}()

	wg.Wait()
}
