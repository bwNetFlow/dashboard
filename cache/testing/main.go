package main

import (
	"fmt"

	"github.com/mediocregopher/radix"
)

func main() {
	fmt.Printf("Servus.\n")

	pool, err := radix.NewPool("tcp", "127.0.0.1:6379", 10)
	if err != nil {
		fmt.Printf("error while connecting to redis: %s", err)
		return
	}
	defer pool.Close()

	/*
		label := counters.NewLabel(map[string]string{
			"ipversion":   "ipv4",
			"application": "web",
			"protoname":   "tcp",
			"direction":   "incoming",
			"cid":         "10109",
			"peer":        "TestPeer",
		})
		item := counters.CounterItems{
			Label: label,
			Value: uint64(4711),
		}

		itemByte, err := json.Marshal(item)
		err = pool.Do(radix.FlatCmd(nil, "SET", "structtest", &itemByte))
		if err != nil {
			fmt.Printf("error: %+v", err)
		}

		var fooVal string
		err = pool.Do(radix.Cmd(&fooVal, "GET", "structtest"))
		if err != nil {
			fmt.Printf("failed to get value from redis: %s", err)
			return
		}
		fmt.Printf("%+v\n", fooVal)

		var parsedItem counters.CounterItems
		json.Unmarshal([]byte(fooVal), &parsedItem)
		fmt.Printf("%+v\n", parsedItem)
	*/

	var fooValX bool
	err = pool.Do(radix.Cmd(&fooValX, "EXISTS", "structtest"))
	if err != nil {
		fmt.Printf("failed to get value from redis: %s", err)
		return
	}
	fmt.Printf("%+v\n", fooValX)

	/*	// This example retrieves the current value of `key` and then sets a new
		// value on it in an atomic transaction.
		key := "someKey"
		var prevVal string

		err = pool.Do(radix.WithConn(key, func(c radix.Conn) error {

			// Begin the transaction with a MULTI command
			if err := c.Do(radix.Cmd(nil, "MULTI")); err != nil {
				return err
			}

			// If any of the calls after the MULTI call error it's important that
			// the transaction is discarded. This isn't strictly necessary if the
			// error was a network error, as the connection would be closed by the
			// client anyway, but it's important otherwise.
			var err error
			defer func() {
				if err != nil {
					// The return from DISCARD doesn't matter. If it's an error then
					// it's a network error and the Conn will be closed by the
					// client.
					c.Do(radix.Cmd(nil, "DISCARD"))
				}
			}()

			// queue up the transaction's commands
			if err = c.Do(radix.Cmd(nil, "GET", key)); err != nil {
				return err
			}
			if err = c.Do(radix.Cmd(nil, "SET", key, "someOtherValue")); err != nil {
				return err
			}

			// execute the transaction, capturing the result
			var result []string
			if err = c.Do(radix.Cmd(&result, "EXEC")); err != nil {
				return err
			}

			// capture the output of the first transaction command, i.e. the GET
			prevVal = result[0]
			return nil
		}))
		if err != nil {
			// handle error
		}

		fmt.Printf("the value of key %q was %q\n", key, prevVal)
	*/
}
