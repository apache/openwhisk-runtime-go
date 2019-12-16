/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package runtime.actionContainers

import actionContainers.ActionContainer.withContainer
import actionContainers.{ActionContainer, BasicActionRunnerTests}
import common.WskActorSystem
import org.junit.runner.RunWith
import org.scalatest.junit.JUnitRunner

@RunWith(classOf[JUnitRunner])
class ActionLoopBasicTests extends BasicActionRunnerTests with WskActorSystem {

  val image = "actionloop-base"

  override def withActionContainer(env: Map[String, String] = Map.empty)(
      code: ActionContainer => Unit) = {
    withContainer(image, env)(code)
  }

  def withActionLoopContainer(code: ActionContainer => Unit) =
    withContainer(image)(code)

  behavior of image

  override val testNoSourceOrExec = TestConfig("")

  override val testNotReturningJson = TestConfig("""#!/bin/bash
      |read line
      |echo '"not json"' >&3
      |read line
      |""".stripMargin)

  override val testEcho = TestConfig("""|#!/bin/bash
       |while read line
       |do
       |    echo "hello stdout"
       |    echo "hello stderr" >&2
       |    echo "$line" | jq -c .value >&3
       |done
    """.stripMargin)

  override val testUnicode = TestConfig(
    """|#!/bin/bash
       |while read line
       |do
       |   delimiter="$(echo "$line" | jq -r ".value.delimiter")"
       |   msg="$delimiter ☃ $delimiter"
       |   echo "$msg"
       |   echo "{\"winter\": \"$msg\"}" >&3
       |done
    """.stripMargin)

  // the __OW_API_HOST should already be in the environment
  // so it is not expected in/read from the input line
  override val testEnv = TestConfig(
    """#!/bin/bash
      |while read line
      |do
      |  __OW_API_KEY="$(echo "$line"        | jq -r .api_key)"
      |  __OW_NAMESPACE="$(echo "$line"      | jq -r .namespace)"
      |  __OW_ACTIVATION_ID="$(echo "$line"  | jq -r .activation_id)"
      |  __OW_ACTION_NAME="$(echo "$line"    | jq -r .action_name)"
      |  __OW_ACTION_VERSION="$(echo "$line" | jq -r .action_version)"
      |  __OW_DEADLINE="$(echo "$line"       | jq -r .deadline)"
      |  echo >&3 "{ \
      |   \"api_host\": \"$__OW_API_HOST\", \
      |   \"api_key\": \"$__OW_API_KEY\", \
      |   \"namespace\": \"$__OW_NAMESPACE\", \
      |   \"activation_id\": \"$__OW_ACTIVATION_ID\", \
      |   \"action_name\": \"$__OW_ACTION_NAME\", \
      |   \"action_version\": \"$__OW_ACTION_VERSION\", \
      |   \"deadline\": \"$__OW_DEADLINE\" }"
      | done
    """.stripMargin)

  val echoSh =
    """|#!/bin/bash
       |while read line
       |do echo "$line" | jq -c .value  >&3
       |done
    """.stripMargin

  override val testInitCannotBeCalledMoreThanOnce = TestConfig(echoSh)

  override val testEntryPointOtherThanMain = TestConfig(echoSh, main = "niam")

  override val testLargeInput = TestConfig(echoSh)
}
