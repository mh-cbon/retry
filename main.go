package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {

	var opts = struct {
		timeout       time.Duration
		retryInterval time.Duration
		tail          bool
		quiet         bool
		okText        string
	}{}

	flag.StringVar(&opts.okText, "ok-text", "", "a text to lookup for inthe stdout/stderr that confirms did run correctly.")
	flag.DurationVar(&opts.timeout, "timeout", time.Minute, "maximum timeout duration")
	flag.DurationVar(&opts.retryInterval, "retry-interval", time.Second, "Retry interval duration")
	flag.BoolVar(&opts.tail, "tail", false, "Apply interval at tail or head")
	flag.BoolVar(&opts.quiet, "quiet", false, "stfu")

	flag.Parse()
	args := flag.Args()

	if len(args) < 0 {
		panic("a command line must be provided")
	}

	ctx, cancel := context.WithTimeout(context.Background(), opts.timeout)
	notifyCancel(cancel)

	var program string
	var programArgs []string
	if len(args) > 0 {
		program = args[0]
	}
	if len(args) > 1 {
		programArgs = args[1:]
	}

	var exitErr error
	combinedOutput := new(bytes.Buffer)
	//keep executin the program until is succeeded
	for {
		combinedOutput.Reset()
		if !opts.tail {
			<-time.After(opts.retryInterval)
		}
		if !opts.quiet {
			sargs := ""
			for _, a := range programArgs {
				if strings.HasPrefix(a, "-") {
					sargs += fmt.Sprintf(" %v", a)
				} else if strings.ContainsAny(a, " ") {
					sargs += fmt.Sprintf(" %q", a)
				} else {
					sargs += fmt.Sprintf(" %v", a)
				}
			}
			log.Print("=> ", program, " ", sargs)
		}
		cmd := exec.CommandContext(ctx, program, programArgs...)
		cmd.Stdout = io.MultiWriter(os.Stdout, combinedOutput)
		cmd.Stderr = io.MultiWriter(os.Stderr, combinedOutput)
		cmd.Stdin = os.Stdin
		// fail if the call to execute does not succeeed
		if err := cmd.Start(); err != nil {
			log.Fatalf("failed to execute %q %v, err=%v", program, programArgs, err)
		}
		//wait for completion,
		exitErr = cmd.Wait()
		if opts.okText == "" {
			// continue unitl it works.
			if exitErr != nil {
				if opts.tail {
					<-time.After(opts.retryInterval)
				}
				continue
			}
		} else {
			// continue unitl it contains desired text.
			if !strings.Contains(combinedOutput.String(), opts.okText) {
				if opts.tail {
					<-time.After(opts.retryInterval)
				}
				continue
			}
		}
		//yeah, it worked! exit asap.
		break
	}

	if exitErr != nil {
		log.Fatalf("failed to execute the command %v %v, err=%v", program, programArgs, exitErr)
	}
}

//notifyCancel context on ctrl+c. This call is non blocking.
func notifyCancel(cancel func(), sigs ...os.Signal) {
	if len(sigs) < 1 {
		sigs = append(sigs, os.Interrupt, os.Kill, syscall.SIGTERM) // super important to use the signal selections,
		// i observed that under load (heavy uncontrolled allocations) the system would send signals
		// to the program, i guess to claim for memory or something similar at hardware level,
		// anyway, it kills the process inadvertely.
		// best is to avoid allocations.
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, sigs...)

	go func() {
		// defer close(sig) // dont, the signal chan might receive many signals instances,
		// that would trigger an error "write on close channel"
		// Anyway, the program is exiting, it will garbage.
		<-sig
		log.Println("got cancellation signal")
		cancel()
	}()
}
