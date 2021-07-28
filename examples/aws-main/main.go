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
package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Main function for the action
func Main(obj map[string]interface{}) map[string]interface{} {
	msg := make(map[string]interface{})

	id, ok1 := obj["id"].(string)
	key, ok2 := obj["key"].(string)
	region, ok3 := obj["region"].(string)
	if ok1 && ok2 && ok3 {
		sess := session.Must(session.NewSession(&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(id, key, ""),
		}))
		service := ec2.New(sess)
		res, err := service.DescribeInstances(&ec2.DescribeInstancesInput{})
		if err != nil {
			msg["error"] = err.Error()
		} else {
			instances := []string{}
			for _, resv := range res.Reservations {
				for _, inst := range resv.Instances {
					instances = append(instances, *inst.InstanceId)
				}
			}
			msg["instances"] = instances
		}
	} else {
		msg["help"] = "required id, key and region"
	}
	return msg
}
