// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	cmd "github.com/yp05327/teachart/cmd/teachart"
	"github.com/yp05327/teachart/pkg/app"
	"github.com/yp05327/teachart/pkg/options"
)

var Cancel context.CancelFunc
var CleanWaitGroup sync.WaitGroup

func main() {
	var sig os.Signal
	sigs := make(chan os.Signal, 1)
	errChan := make(chan error, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()
	ctx, Cancel = context.WithCancel(ctx)

	go func() {
		// set logger formatter
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableTimestamp:       true,
			DisableLevelTruncation: true,
		})

		// for execute all persistent pre-runs from the root parent till this command.
		cobra.EnableTraverseRunHooks = true

		globalOptions := options.NewGlobalOptions(app.DefaultRepoDir)
		rootCmd, err := cmd.NewRootCmd(ctx, globalOptions)
		if err != nil {
			errChan <- err
			return
		}

		errChan <- rootCmd.Execute()
	}()

	select {
	case sig = <-sigs:
		if sig != nil {
			Cancel()

			// See http://tldp.org/LDP/abs/html/exitcodes.html
			switch sig {
			case syscall.SIGINT:
				os.Exit(130)
			case syscall.SIGTERM:
				os.Exit(143)
			}
		}
	case err := <-errChan:
		if err != nil {
			logrus.Error(err)
		}
	}
}
