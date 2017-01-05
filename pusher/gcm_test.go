/*
 * Copyright (c) 2016 TFG Co <backend@tfgco.com>
 * Author: TFG Co <backend@tfgco.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package pusher

import (
	"time"

	"github.com/Sirupsen/logrus/hooks/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GCM Pusher", func() {
	configFile := "../config/test.yaml"
	senderID := "sender-id"
	apiKey := "api-key"
	appName := "testapp"
	logger, hook := test.NewNullLogger()

	BeforeEach(func() {
		hook.Reset()
	})

	Describe("Creating new gcm pusher", func() {
		It("should return configured pusher", func() {
			pusher := NewGCMPusher(configFile, senderID, apiKey, appName, logger)
			Expect(pusher).NotTo(BeNil())
			Expect(pusher.apiKey).To(Equal(apiKey))
			Expect(pusher.AppName).To(Equal(appName))
			Expect(pusher.Config).NotTo(BeNil())
			Expect(pusher.ConfigFile).To(Equal(configFile))
			Expect(pusher.MessageHandler).NotTo(BeNil())
			Expect(pusher.Queue).NotTo(BeNil())
			Expect(pusher.run).To(BeFalse())
			Expect(pusher.senderID).To(Equal(senderID))
		})
	})

	Describe("Start gcm pusher", func() {
		It("should launch go routines and run forever", func() {
			pusher := NewGCMPusher(configFile, senderID, apiKey, appName, logger)
			Expect(pusher).NotTo(BeNil())
			defer func() { pusher.run = false }()
			go pusher.Start()
			time.Sleep(50 * time.Millisecond)
			Expect(pusher.run).To(BeTrue())
			// TODO test signals
			// Sending SIGTERM makes the test abort
			// syscall.Kill(os.Getpid(), syscall.SIGTERM)
			// Expect(pusher.run).To(BeFalse())
		})
	})
})