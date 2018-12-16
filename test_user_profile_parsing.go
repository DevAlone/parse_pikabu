package main

import (
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/results_processor"
	"github.com/streadway/amqp"
	"sync"
)

func main() {
	err := models.InitDb()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// start server
	/*go func() {
		err := server.Run()
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	// start task manager
	go func() {
		err := task_manager.Run()
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()*/

	// start results processor
	go func() {
		err := results_processor.Run()
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	err = pushTaskToQueue([]byte(`{"user":{"current_user_id":0,"user_id":"2561615","user_name":"Pisacavtor","rating":"-3.5","gender":"0","comments_count":3,"stories_count":1,"stories_hot_count":"0","pluses_count":0,"minuses_count":0,"signup_date":"1544846469","is_rating_ban":false,"avatar":"https:\/\/cs8.pikabu.ru\/avatars\/2561\/x2561615-512432259.png","awards":[],"is_subscribed":false,"is_ignored":false,"note":null,"approved":"","communities":[],"subscribers_count":0,"is_user_banned":true,"is_user_fully_banned":false,"public_ban_history":[{"id":"151513","date":1544854692,"moderator_id":"1836690","comment_id":"0","comment_desc":"","story_id":"6354471","user_id":"2561615","reason":"\u041e\u0442\u0441\u0443\u0442\u0441\u0442\u0432\u0438\u0435 \u043f\u0440\u0443\u0444\u0430 \u0438\u043b\u0438 \u043d\u0435\u043f\u043e\u0434\u0442\u0432\u0435\u0440\u0436\u0434\u0451\u043d\u043d\u0430\u044f\/\u0438\u0441\u043a\u0430\u0436\u0451\u043d\u043d\u0430\u044f \u0438\u043d\u0444\u043e\u0440\u043c\u0430\u0446\u0438\u044f (\u0432\u0431\u0440\u043e\u0441)","reason_id":"94","story_url":"https:\/\/pikabu.ru\/story\/3_chasa_pyitok_6354471","moderator_name":"depotato","moderator_avatar":"https:\/\/cs5.pikabu.ru\/avatars\/1836\/s1836690-1399622318.png","reason_limit":null,"reason_count":null,"reason_title":null}],"user_ban_time":1545459492}}`))
	if err != nil {
		panic(err)
	}
	// TODO: test pikago.UserProfile serialization

	wg.Wait()

}

func pushTaskToQueue(message []byte) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"parser_results",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"parser_results",
		"user_profile",
		true,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         message,
		},
	)
	return err
}
