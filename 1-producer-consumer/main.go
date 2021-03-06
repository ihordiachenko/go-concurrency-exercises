//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"time"
)

func producer(stream Stream) <-chan *Tweet {
	tweets := make(chan *Tweet)

	go (func() {
		for {
			tweet, err := stream.Next()
			if err == ErrEOF {
				close(tweets)
				return
			}

			tweets <- tweet
		}
	})()

	return tweets
}

func consumer(tweets <-chan *Tweet) {
	for {
		tweet, more := <-tweets
		if more {
			if tweet.IsTalkingAboutGo() {
				fmt.Println(tweet.Username, "\ttweets about golang")
			} else {
				fmt.Println(tweet.Username, "\tdoes not tweet about golang")
			}
		} else {
			return
		}
	}
}

func main() {
	start := time.Now()

	stream := GetMockStream()

	// Producer
	tweets := producer(stream)

	// Consumer
	consumer(tweets)

	fmt.Printf("Process took %s\n", time.Since(start))
}
