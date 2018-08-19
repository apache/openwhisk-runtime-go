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

object GoResourceHelpers {

  /** Create a temporary directory in your home directory.
    *
    * This is needed to use docker volume mounts.
    * On mac I need to use the home directory,
    * because the default folder used by gradle under Mac
    * is not accessible by default by Docker for Mac
    *
    */
  def createHomeTmpDirectory(prefix: String, suffix: String = "") = {
    val srcFileDir = new File(new File(System.getProperty("user.home"), "tmp"), prefix+System.currentTimeMillis().toString+suffix)
    srcFileDir.mkdirs()
    srcFileDir.toPath.toAbsolutePath()
  }

  private def makeZipFromDir(dir: Path): Path = makeArchiveFromDir(dir, ".zip")

  private def makeJarFromDir(dir: Path): Path = makeArchiveFromDir(dir, ".jar")

  /**
    * Compresses all files beyond a directory into a zip file.
    * Note that Jar files are just zip files.
    */
  private def makeArchiveFromDir(dir: Path, extension: String): Path = {
    // Any temporary file name for the archive.
    val arPath = createHomeTmpDirectory("output", extension).toAbsolutePath()

    // We "mount" it as a filesystem, so we can just copy files into it.
    val dstUri = new URI("jar:" + arPath.toUri().getScheme(), arPath.toAbsolutePath().toString(), null)
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
      })

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
  private def writeSourcesToHomeTmpDirectory(sources: Seq[(Seq[String], String)]): (Path, Seq[Path]) = {
    // A temporary directory for the source files.
    val srcDir = createHomeTmpDirectory("src")
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
    // return the exe File  and the output dir Path
    private def compile(image: String, sources: Seq[(Seq[String], String)], main: String) = {
      require(!sources.isEmpty)

      // The absolute paths of the source file
      val (srcDir, srcAbsPaths) = writeSourcesToHomeTmpDirectory(sources)
      val src = srcDir.toFile.getAbsolutePath

      // A temporary directory for the destination files.
      val outDir = createHomeTmpDirectory("out").toAbsolutePath()
      val out = outDir.toFile.getAbsolutePath

      // command to compile
      val exe = new File(out, main)
      val cmd = s"${dockerBin} run -v ${src}:/src -v ${out}:/out  ${image} compile ${main}"

      // compiling
      import sys.process._
      cmd.!

      // result
      exe -> outDir

    }

    def mkBase64Exe(image: String, sources: Seq[(Seq[String], String)], main: String) = {
      val (exe, dir) = compile(image, sources, main)
      //println(s"exe=${exe.getAbsolutePath}")
      ResourceHelpers.readAsBase64(exe.toPath)
    }

    def mkBase64Zip(image: String, sources: Seq[(Seq[String], String)], main: String) = {
      val (exe, dir) = compile(image, sources, main)
      val archive = makeZipFromDir(dir)
      //println(s"zip=${archive.toFile.getAbsolutePath}")
      ResourceHelpers.readAsBase64(archive)
    }

    def mkBase64Src(sources: Seq[(Seq[String], String)], main: String) = {
      val (srcDir, srcAbsPaths) = writeSourcesToHomeTmpDirectory(sources)
      val file = new File(srcDir.toFile, main)
      //println(file)
      ResourceHelpers.readAsBase64(file.toPath)
    }

    def mkBase64SrcZip(sources: Seq[(Seq[String], String)], main: String) = {
      val (srcDir, srcAbsPaths) = writeSourcesToHomeTmpDirectory(sources)
      //println(srcDir)
      val archive = makeZipFromDir(srcDir)
      ResourceHelpers.readAsBase64(archive)
    }
  }

  def writeFile(name: String, body: String): Unit = {
    val fw = new java.io.FileWriter(name)
    fw.write(body)
    fw.close
  }
}
