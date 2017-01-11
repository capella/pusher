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
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/topfreegames/pusher/extensions"
	"github.com/topfreegames/pusher/interfaces"
)

type statsReporterInitializer func(string, *logrus.Logger, string) (interfaces.StatsReporter, error)
type feedbackReporterInitializer func(string, *logrus.Logger) (interfaces.FeedbackReporter, error)

//AvailableStatsReporters contains functions to initialize all stats reporters
var AvailableStatsReporters = map[string]statsReporterInitializer{
	"statsd": func(configFile string, logger *logrus.Logger, appName string) (interfaces.StatsReporter, error) {
		return extensions.NewStatsD(configFile, logger, appName)
	},
}

//AvailableFeedbackReporters contains functions to initialize all feedback reporters
var AvailableFeedbackReporters = map[string]feedbackReporterInitializer{
	"kafka": func(configFile string, logger *logrus.Logger) (interfaces.FeedbackReporter, error) {
		return extensions.NewKafkaProducer(configFile, logger)
	},
}

func configureStatsReporters(configFile string, logger *logrus.Logger, appName string, config *viper.Viper) ([]interfaces.StatsReporter, error) {
	reporters := []interfaces.StatsReporter{}
	reporterNames := config.GetStringSlice("stats.reporters")
	for _, reporterName := range reporterNames {
		reporterFunc, ok := AvailableStatsReporters[reporterName]
		if !ok {
			return nil, fmt.Errorf("Failed to initialize %s. Stats Reporter not available.", reporterName)
		}

		r, err := reporterFunc(configFile, logger, appName)
		if err != nil {
			return nil, fmt.Errorf("Failed to initialize %s. %s", reporterName, err.Error())
		}
		reporters = append(reporters, r)
	}

	return reporters, nil
}

func configureFeedbackReporters(configFile string, logger *logrus.Logger, config *viper.Viper) ([]interfaces.FeedbackReporter, error) {
	reporters := []interfaces.FeedbackReporter{}
	reporterNames := config.GetStringSlice("feedback.reporters")
	for _, reporterName := range reporterNames {
		reporterFunc, ok := AvailableFeedbackReporters[reporterName]
		if !ok {
			return nil, fmt.Errorf("Failed to initialize %s. Feedback Reporter not available.", reporterName)
		}

		r, err := reporterFunc(configFile, logger)
		if err != nil {
			return nil, fmt.Errorf("Failed to initialize %s. %s", reporterName, err.Error())
		}
		reporters = append(reporters, r)
	}

	return reporters, nil
}
