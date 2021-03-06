// Copyright 2016-2019 Yincheng Fang, Alex Stocks
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

import (
	"github.com/dubbo/go-for-apache-dubbo/common/logger"
	_ "github.com/dubbo/go-for-apache-dubbo/common/proxy/proxy_factory"
	"github.com/dubbo/go-for-apache-dubbo/config"
	_ "github.com/dubbo/go-for-apache-dubbo/protocol/jsonrpc"
	_ "github.com/dubbo/go-for-apache-dubbo/registry/protocol"

	_ "github.com/dubbo/go-for-apache-dubbo/filter/impl"

	_ "github.com/dubbo/go-for-apache-dubbo/cluster/cluster_impl"
	_ "github.com/dubbo/go-for-apache-dubbo/cluster/loadbalance"
	_ "github.com/dubbo/go-for-apache-dubbo/registry/zookeeper"
)

var (
	survivalTimeout = int(3e9)
)

// they are necessary:
// 		export CONF_PROVIDER_FILE_PATH="xxx"
// 		export APP_LOG_CONF_FILE="xxx"
func main() {

	_, proMap := config.Load()
	if proMap == nil {
		panic("proMap is nil")
	}

	initSignal()
}

func initSignal() {
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		logger.Infof("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
		// reload()
		default:
			go time.AfterFunc(time.Duration(float64(survivalTimeout)*float64(time.Second)), func() {
				logger.Warnf("app exit now by force...")
				os.Exit(1)
			})

			// 要么fastFailTimeout时间内执行完毕下面的逻辑然后程序退出，要么执行上面的超时函数程序强行退出
			fmt.Println("provider app exit now...")
			return
		}
	}
}
