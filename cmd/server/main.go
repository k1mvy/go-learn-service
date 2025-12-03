package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron/v2"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	// cron scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	_, err = addJobs(s)
	if err != nil {
		fmt.Println(err)
	}
	//jobs, err := addJobs(s)
	//for _, job := range jobs {
	//	fmt.Printf("running job: %v", job.ID())
	//}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	// http
	go func() {
		defer wg.Done()

		fmt.Println("starting server")
		err = http.ListenAndServe(":3000", r)
		if err != nil {
			panic(err)
		}
	}()

	// cron
	go func() {
		defer wg.Done()

		fmt.Println("starting schedule")
		s.Start()
	}()

	wg.Wait()
}

func addJobs(s gocron.Scheduler) ([]gocron.Job, error) {
	j, err := s.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func() {
				fmt.Println("job completed")
			},
		),
	)
	if err != nil {
		return nil, err
	}

	return []gocron.Job{j}, nil
}
