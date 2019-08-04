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

import actionContainers.{ActionContainer, ActionProxyContainerTestUtils}
import actionContainers.ActionContainer.withContainer
import common.WskActorSystem
import org.junit.runner.RunWith
import org.scalatest.junit.JUnitRunner

import spray.json.{JsObject, JsString}

@RunWith(classOf[JUnitRunner])
class ActionLoopContainerTests
    extends ActionProxyContainerTestUtils
    with WskActorSystem {

  import GoResourceHelpers._

  val image = "actionloop-base"

  def withActionLoopContainer(code: ActionContainer => Unit) =
    withContainer(image)(code)

  behavior of image

  def shCodeHello(main: String) = Seq(
    Seq(main) ->
      s"""#!/bin/bash
         |while read line
         |do
         |   name="$$(echo $$line | jq -r .value.name)"
         |   if test "$$name" == ""
         |   then exit
         |   fi
         |   echo "name=$$name"
         |   hello="Hello, $$name"
         |   echo '{"${main}":"'$$hello'"}' >&3
         |done
         |""".stripMargin
  )

  private def helloMsg(name: String = "Demo") =
    runPayload(JsObject("name" -> JsString(name)))

  private def okMsg(key: String, value: String) =
    200 -> Some(JsObject(key -> JsString(value)))

  it should "run sample with init that does nothing" in {
    val (out, err) = withActionLoopContainer { c =>
      c.init(JsObject())._1 should be(403)
      c.run(JsObject())._1 should be(500)
    }
  }

  it should "deploy a shell script" in {
    val script = shCodeHello("main")
    val mesg = ExeBuilder.mkBase64Src(script)
    withActionLoopContainer { c =>
      val payload = initPayload(mesg)
      c.init(payload)._1 should be(200)
      c.run(helloMsg()) should be(okMsg("main", "Hello, Demo"))
    }
  }

  it should "deploy a zip based script" in {
    val scr = ExeBuilder.mkBase64SrcZip(shCodeHello("exec"))
    withActionLoopContainer { c =>
      c.init(initPayload(scr))._1 should be(200)
      c.run(helloMsg()) should be(okMsg("exec", "Hello, Demo"))
    }
  }
}
