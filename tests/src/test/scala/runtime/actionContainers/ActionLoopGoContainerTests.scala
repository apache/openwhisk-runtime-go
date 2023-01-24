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

import spray.json.{JsArray, JsObject, JsString}

abstract class ActionLoopGoContainerTests
    extends ActionProxyContainerTestUtils
    with WskActorSystem {

  import GoResourceHelpers._

  val goCompiler: String
  val image: String

  def withActionLoopContainer(code: ActionContainer => Unit) =
    withContainer(image)(code)

  behavior of image

  def helloGo(main: String, pkg: String = "main") = {
    val func = main.capitalize
    s"""|package ${pkg}
        |
        |import "fmt"
        |
        |func ${func}(obj map[string]interface{}) map[string]interface{} {
        |	 name, ok := obj["name"].(string)
        |	 if !ok {
        |	  	name = "Stranger"
        |	 }
        |	 fmt.Printf("name=%s\\n", name)
        |  msg := make(map[string]interface{})
        |	 msg["${pkg}-${main}"] = "Hello, " + name + "!"
        |	 return msg
        |}
        |""".stripMargin
  }

  def helloSrc(main: String) = Seq(
    Seq(s"${main}.go") -> helloGo(main)
  )

  def helloMsg(name: String = "Demo") =
    runPayload(JsObject("name" -> JsString(name)))

  def okMsg(key: String, value: String) =
    200 -> Some(JsObject(key -> JsString(value)))

  it should "run sample with init that does nothing" in {
    val (out, err) = withActionLoopContainer { c =>
      c.init(JsObject())._1 should be(403)
      c.run(JsObject())._1 should be(500)
    }
  }

  it should "accept a binary main" in {
    val exe = ExeBuilder.mkBase64Zip(goCompiler, helloSrc("main"), "main")

    withActionLoopContainer { c =>
      c.init(initPayload(exe))._1 shouldBe (200)
      c.run(helloMsg()) should be(okMsg("main-main", "Hello, Demo!"))
    }
  }

  it should "accept a src main action " in {
    var src = ExeBuilder.mkBase64Src(helloSrc("main"))
    withActionLoopContainer { c =>
      c.init(initPayload(src))._1 shouldBe (200)
      c.run(helloMsg()) should be(okMsg("main-main", "Hello, Demo!"))
    }
  }

  it should "accept a src not-main action " in {
    var src = ExeBuilder.mkBase64Src(helloSrc("hello"))
    withActionLoopContainer { c =>
      c.init(initPayload(src, "hello"))._1 shouldBe (200)
      c.run(helloMsg()) should be(okMsg("main-hello", "Hello, Demo!"))
    }
  }

  it should "accept a zipped src main action" in {
    var src = ExeBuilder.mkBase64SrcZip(helloSrc("main"))
    withActionLoopContainer { c =>
      c.init(initPayload(src))._1 shouldBe (200)
      c.run(helloMsg()) should be(okMsg("main-main", "Hello, Demo!"))
    }
  }

  it should "accept a zipped src not-main action" in {
    var src = ExeBuilder.mkBase64SrcZip(helloSrc("hello"))
    withActionLoopContainer { c =>
      c.init(initPayload(src, "hello"))._1 shouldBe (200)
      c.run(helloMsg()) should be(okMsg("main-hello", "Hello, Demo!"))
    }
  }

  it should "deploy a zip main src with subdir" in {
    var src = ExeBuilder.mkBase64SrcZip(
      Seq(
        Seq("hello", "hello.go") -> helloGo("Hello", "hello"),
        Seq("hello", "go.mod") -> "module hello\n",
        Seq("main.go") ->
          """
          |package main
          |import "hello"
          |func Main(args map[string]interface{})map[string]interface{} {
          | return hello.Hello(args)
          |}
        """.stripMargin,
        Seq("go.mod") -> "module action\nreplace hello => ./hello\nrequire hello v0.0.0-00010101000000-000000000000\n"
      )
    )
    withActionLoopContainer { c =>
      c.init(initPayload(src))._1 shouldBe (200)
      c.run(helloMsg()) should be(okMsg("hello-Hello", "Hello, Demo!"))
    }
  }

  it should "support return array result" in {
    val helloArrayGo = {
      s"""
         |package main
         |
         |func Main(obj map[string]interface{}) []interface{} {
         |    result := []interface{}{"a", "b"}
         |    return result
         |}
         |
       """.stripMargin
    }
    val src = ExeBuilder.mkBase64SrcZip(
      Seq(
        Seq(s"main.go") -> helloArrayGo
      ))
    withActionLoopContainer { c =>
      c.init(initPayload(src))._1 shouldBe (200)
      val result = c.runForJsArray(JsObject())
      result._1 shouldBe (200)
      result._2 shouldBe Some(JsArray(JsString("a"), JsString("b")))
    }
  }

  it should "support array as input param" in {
    val helloArrayGo = {
      s"""
         |package main
         |
         |func Main(obj []interface{}) []interface{} {
         |    return obj
         |}
         |
       """.stripMargin
    }
    val src = ExeBuilder.mkBase64SrcZip(
      Seq(
        Seq(s"main.go") -> helloArrayGo
      ))
    withActionLoopContainer { c =>
      c.init(initPayload(src))._1 shouldBe (200)
      val result =
        c.runForJsArray(runPayload(JsArray(JsString("a"), JsString("b"))))
      result._1 shouldBe (200)
      result._2 shouldBe Some(JsArray(JsString("a"), JsString("b")))
    }
  }
}
