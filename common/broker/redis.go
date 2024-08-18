package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"my-stocks/common/aggregator"
	"my-stocks/common/queues"
)

type RedisBroker struct {
	server *redis.Client
	ctx    context.Context
}

func NewRedisBroker(server *redis.Client, ctx context.Context) *RedisBroker {
	return &RedisBroker{server: server, ctx: ctx}
}

func (r *RedisBroker) Consume(queue queues.QueueName, group queues.QueueNameGroup, sender chan<- *Message) error {
	//_, err := r.server.XGroupDestroy(r.ctx, string(queue), string(group)).Result()
	result, err := r.server.XInfoGroups(r.ctx, string(queue)).Result()
	found := false
	if err == nil {
		for _, infoGroup := range result {
			if infoGroup.Name == string(group) {
				found = true
				break
			}
		}
	}
	if !found || err != nil {
		err = r.server.XGroupCreateMkStream(r.ctx, string(queue), string(group), "0").Err()
		if err != nil {
			fmt.Println(queue, group)
			log.Fatalf("error x %v", err)
		}
	}

	for {
		entries, err := r.server.XReadGroup(r.ctx, &redis.XReadGroupArgs{
			Group:   string(group),
			Streams: []string{string(queue), ">"},
			Count:   1,
			Block:   0,
			NoAck:   true,
		}).Result()

		if err != nil {
			log.Fatalf("error x %v", err)
		}

		for i := 0; i < len(entries[0].Messages); i++ {
			var tmp any
			_ = json.Unmarshal([]byte(entries[0].Messages[i].Values["Data"].(string)), &tmp)
			fmt.Println(entries[0].Messages[i].Values["Data"])
			sender <- &Message{
				ID:    entries[0].Messages[i].ID,
				Group: group,
				Queue: queue,
				Data: &aggregator.Correlation{
					CorrelationId: entries[0].Messages[i].Values["CorrelationId"].(string),
					Data:          tmp,
				},
			}
		}
	}
}

func (r *RedisBroker) ConsumeWithHandler(queue queues.QueueName, group queues.QueueNameGroup, handler ConsumerOnly) error {
	//_, err := r.server.XGroupDestroy(r.ctx, string(queue), string(group)).Result()
	result, err := r.server.XInfoGroups(r.ctx, string(queue)).Result()
	found := false
	if err == nil {
		for _, infoGroup := range result {
			if infoGroup.Name == string(group) {
				found = true
				break
			}
		}
	}
	if !found || err != nil {
		err = r.server.XGroupCreateMkStream(r.ctx, string(queue), string(group), "0").Err()
		if err != nil {
			fmt.Println(queue, group)
			log.Fatalf("error x %v", err)
		}
	}

	for {
		entries, err := r.server.XReadGroup(r.ctx, &redis.XReadGroupArgs{
			Group:   string(group),
			Streams: []string{string(queue), ">"},
			Count:   1,
			Block:   0,
			NoAck:   true,
		}).Result()

		if err != nil {
			log.Fatalf("error x %v", err)
		}

		for i := 0; i < len(entries[0].Messages); i++ {
			var tmp any
			_ = json.Unmarshal([]byte(entries[0].Messages[i].Values["Data"].(string)), &tmp)
			handler(&Message{
				ID:    entries[0].Messages[i].ID,
				Group: group,
				Queue: queue,
				Data: &aggregator.Correlation{
					CorrelationId: entries[0].Messages[i].Values["CorrelationId"].(string),
					Data:          tmp,
				},
			})
		}
	}
}

func (r *RedisBroker) Ack(queue queues.QueueName, group queues.QueueNameGroup, messageId string) error {
	return r.server.XAck(r.ctx, string(queue), string(group), messageId).Err()
}

func (r *RedisBroker) Push(queue queues.QueueName, data *aggregator.Correlation) error {
	data.Data, _ = json.Marshal(data.Data)
	marshal := structs.Map(data)
	err := r.server.XAdd(r.ctx, &redis.XAddArgs{
		Stream: string(queue),
		Values: marshal,
	}).Err()
	return err
}
