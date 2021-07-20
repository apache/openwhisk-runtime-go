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

import java.io.File
import java.net.URI
import java.nio.file._
import java.nio.charset.StandardCharsets
import java.nio.file.attribute.BasicFileAttributes

import collection.JavaConverters._
import actionContainers.ResourceHelpers

import scala.util.Random

object GoResourceHelpers {

  /** Create a temporary directory in your /tmp directory.
    * /tmp/openwhisk/random-lowercase-string/prefix+suffix
    *
    * This is needed to use docker volume mounts.
    * On mac I need to use the /tmp directory,
    * because the default folder used by gradle under Mac
    * is not accessible by default by Docker for Mac
    *
    */
  def tmpDirectoryFile(prefix: String, suffix: String = "") =
    new File(
      new File(new File("/tmp", "openwhisk"),
               Random.alphanumeric
                 .take(10)
                 .toArray
                 .mkString
                 .toLowerCase /*random filename alphanumeric and lower case*/ ),
      prefix ++ suffix
    )

  def createTmpDirectory(prefix: String, suffix: String = "") = {
    val tmpDir = tmpDirectoryFile(prefix, suffix)
    tmpDir.mkdirs()
    tmpDir.toPath.toAbsolutePath
  }

  private def makeZipFromDir(dir: Path): Path = makeArchiveFromDir(dir, ".zip")

  private def makeJarFromDir(dir: Path): Path = makeArchiveFromDir(dir, ".jar")

  /**
    * Compresses all files beyond a directory into a zip file.
    * Note that Jar files are just zip files.
    */
  private def makeArchiveFromDir(dir: Path, extension: String): Path = {
    // Any temporary file name for the archive.
    val arPath = createTmpDirectory("output", extension).toAbsolutePath()

    // We "mount" it as a filesystem, so we can just copy files into it.
    val dstUri = new URI("jar:" + arPath.toUri().getScheme(),
                         arPath.toAbsolutePath().toString(),
                         null)
    // OK, that's a hack. Doing this because newFileSystem wants to create that file.
    arPath.toFile().delete()
    val fs = FileSystems.newFileSystem(dstUri, Map(("create" -> "true")).asJava)

    // Traversing all files in the bin directory...
    Files.walkFileTree(
      dir,
      new SimpleFileVisitor[Path]() {
        override def visitFile(path: Path, attributes: BasicFileAttributes) = {
          // The path relative to the src dir
          val relPath = dir.relativize(path)

          // The corresponding path in the zip
          val arRelPath = fs.getPath(relPath.toString())

          // If this file is not top-level in the src dir...
          if (relPath.getParent() != null) {
            // ...create the directory structure if it doesn't exist.
            if (!Files.exists(arRelPath.getParent())) {
              Files.createDirectories(arRelPath.getParent())
            }
          }

          // Finally we can copy that file.
          Files.copy(path, arRelPath)

          FileVisitResult.CONTINUE
        }
      }
    )

    fs.close()

    arPath
  }

  /**
    * Creates a temporary directory in the home and reproduces the desired file structure
    * in it. Returns the path of the temporary directory and the path of each
    * file as represented in it.
    *
    * NOTE this is different from writeSourcesToTempDirectory because it uses
    * a tmp folder in the home directory.
    *
    * */
  private def writeSourcesToHomeTmpDirectory(
      sources: Seq[(Seq[String], String)]): (Path, Seq[Path]) = {
    // A temporary directory for the source files.
    val srcDir = createTmpDirectory("src")
    val srcAbsPaths = for ((sourceName, sourceContent) <- sources) yield {
      // The relative path of the source file
      val srcRelPath = Paths.get(sourceName.head, sourceName.tail: _*)
      // The absolute path of the source file
      val srcAbsPath = srcDir.resolve(srcRelPath)
      // Create parent directories if needed.
      Files.createDirectories(srcAbsPath.getParent)
      // Writing contents
      Files.write(srcAbsPath, sourceContent.getBytes(StandardCharsets.UTF_8))

      srcAbsPath
    }

    (srcDir, srcAbsPaths)
  }

  /**
    * Builds executables using docker images as compilers.
    * Images are assumed to be able to be run as
    * docker <image> -v <sources>:/src -v <output>:/out compile <main>
    * The compiler will read sources from /src and leave the final binary in /out/<main>
    * <main> is also the name of the main function to be invoked
    * Implementations available for swift and go
    */
  object ExeBuilder {

    private lazy val dockerBin: String = {
      List("/usr/bin/docker", "/usr/local/bin/docker").find { bin =>
        new File(bin).isFile
      }.get // This fails if the docker binary couldn't be located.
    }

    // prepare sources, then compile them
    // return the zip File
    private def compile(image: String,
                        sources: Seq[(Seq[String], String)],
                        main: String) = {
      require(!sources.isEmpty)

      // The absolute paths of the source file
      val (srcDir, srcAbsPaths) = writeSourcesToHomeTmpDirectory(sources)
      val src = srcAbsPaths.head.toFile

      // A temporary directory for the destination files.
      // DO NOT CREATE IT IN ADVANCE or you will get a permission denied
      val binDir = tmpDirectoryFile("bin")
      binDir.mkdirs()
      val zip = new File(binDir, main + ".zip")

      // command to compile
      val cmd = s"${dockerBin} run -i ${image} -compile ${main}"

      // compiling
      //println(s"${cmd}\n<${src}\n>${bin}")

      import sys.process._
      (src #> cmd #> zip).!

      // result
      zip
    }

    def mkBase64Zip(image: String,
                    sources: Seq[(Seq[String], String)],
                    main: String) = {
      val zip = compile(image, sources, main)
      ResourceHelpers.readAsBase64(zip.toPath)
    }

    def mkBase64Src(sources: Seq[(Seq[String], String)]) = {
      val (srcDir, srcAbsPaths) = writeSourcesToHomeTmpDirectory(sources)
      ResourceHelpers.readAsBase64(srcAbsPaths.head)
    }

    def mkBase64SrcZip(sources: Seq[(Seq[String], String)]) = {
      val (srcDir, srcAbsPaths) = writeSourcesToHomeTmpDirectory(sources)
      val archive = makeZipFromDir(srcDir)
      //println(s"zip=${archive.toFile.getAbsolutePath}")
      ResourceHelpers.readAsBase64(archive)
    }
  }

  def writeFile(name: String, body: String): Unit = {
    val fw = new java.io.FileWriter(name)
    fw.write(body)
    fw.close
  }
}
