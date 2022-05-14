package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ChromaMinecraft/email-service/config"
	"github.com/ChromaMinecraft/email-service/rabbitmq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/streadway/amqp"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}

	os.Setenv("TZ", cfg.TZ)

	fmt.Printf("%v", cfg)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	amqpAddr := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.AMQP.User, cfg.AMQP.Password, cfg.AMQP.Host, cfg.AMQP.Port)

	conn, err := amqp.Dial(amqpAddr)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	// mailReq := mailer.NewRequest(
	// 	mailer.WithHost(cfg.MailHost),
	// 	mailer.WithPort(cfg.MailPort),
	// 	mailer.WithFrom(cfg.MailUser),
	// 	mailer.WithSubject("Test"),
	// )

	// mailAuth := smtp.PlainAuth("", cfg.MailUser, cfg.MailPassword, cfg.MailHost)
	// mailClient := mailer.NewMailer(mailAuth, mailReq)

	// fmt.Println(mailClient)

	ampqClient := rabbitmq.NewRabbitMQ(cfg, conn)

	ch, err := ampqClient.GetChannel()
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
	}

	defer ch.Close()

	err = ch.Qos(
		cfg.AMQP.PrefetchCount, // prefetch count
		cfg.AMQP.PrefetchSize,  // prefetch size
		cfg.AMQP.IsGlobal,      // global
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to set QoS", err)
	}

	q, err := ch.QueueDeclare(
		cfg.AMQP.AMQPQueue.Name,         // name
		cfg.AMQP.AMQPQueue.IsDurable,    // durable
		cfg.AMQP.AMQPQueue.IsAutoDelete, // delete when unused
		cfg.AMQP.AMQPQueue.IsExclusive,  // exclusive
		cfg.AMQP.AMQPQueue.IsNoWait,     // no-wait
		nil,                             // arguments
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	p := NewPool(cfg.WorkerCount)

	p.Run()

	// Subscribe to the queue
	var counter int = 1
	go ampqClient.Subscribe(ch, q, func(d amqp.Delivery) {
		log.Printf("Received a message: %s", d.Body)
		p.JobQueue <- Job{ID: counter, Resources: string(d.Body)}
		counter++
	})

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan

	// Handle shutdown
	fmt.Println("------- Shutdown Signal Received -------")
	fmt.Println("All workers done, shutting down!")
}

// Job represents the job payload
type Job struct {
	ID        int
	Resources string
}

// Pool represents the worker pool orchestrator
type Pool struct {
	NumWorkers int
	JobsChan   chan chan Job
	JobQueue   chan Job
	Stopped    chan bool
}

// Worker represents the actual worker who does the job
type Worker struct {
	ID          int
	JobChannel  chan Job
	JobChannels chan chan Job
	Quit        chan bool
}

func NewPool(num int) Pool {
	return Pool{
		NumWorkers: num,
		JobsChan:   make(chan chan Job),
		JobQueue:   make(chan Job),
		Stopped:    make(chan bool),
	}
}

func (p *Pool) Run() {
	log.Println("Spawn workers")
	for i := 0; i <= p.NumWorkers; i++ {
		worker := Worker{
			ID:          i,
			JobChannel:  make(chan Job),
			JobChannels: p.JobsChan,
			Quit:        make(chan bool),
		}
		worker.Start()
	}
	p.Allocate()
}

func (p *Pool) Allocate() {
	q := p.JobQueue
	s := p.Stopped

	go func(queue chan Job) {
		for {
			select {
			case job := <-queue:
				log.Printf("Worker %d got a job: %v", job.ID, job)
				availChannel := <-p.JobsChan
				availChannel <- job
			case <-s:
				log.Printf("Worker %d stopping", p.NumWorkers)
				return
			}
		}
	}(q)
}

func (w *Worker) Start() {
	log.Printf("Starting Worker ID [%d]", w.ID)
	go func() {
		for {
			w.JobChannels <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				w.work(job)
			case <-w.Quit:
				return
			}
		}
	}()
}

// work is the task performer
func (w *Worker) work(job Job) {
	log.Printf("-------")
	log.Printf("Processed by Worker [%d]", w.ID)
	log.Printf("Processed Job With ID [%d] & content: [%s]", job.ID, job.Resources)
	log.Printf("-------")
}
