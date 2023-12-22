/*******************************************************************************
 * Copyright 2023-2023 Edw590
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 ******************************************************************************/

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"Utils"
)

// Email Sender //

type _FileInfo struct {
	file_name  string
	modif_time int64
}

// MAX_EMAILS_HOUR is the maximum number of emails that can be sent per hour according to Google, which are 20. But
// we'll go with 15 to be safe about emails sent without this module's help (error emails).
const MAX_EMAILS_HOUR = 15

type _MGIModSpecInfo _ModSpecInfo
var (
	realMain        Utils.RealMain = nil
	modStartInfo_GL Utils.ModStartInfo
	modGenInfo_GL   Utils.ModGenInfo[_MGIModSpecInfo]
)

func main() {Utils.ModStartup[_MGIModSpecInfo](Utils.NUM_MOD_EmailSender, realMain)}
func init() {realMain =
	func(realMain_param_1 Utils.ModStartInfo, realMain_param_2 any) {
		modStartInfo_GL = realMain_param_1
		modGenInfo_GL = realMain_param_2.(Utils.ModGenInfo[_MGIModSpecInfo])

		var to_send_dir Utils.GPath = modStartInfo_GL.UserData_dir.Add2(Utils.TO_SEND_REL_FOLDER)

		fmt.Println("Checking for emails to send in \"" + to_send_dir.GPathToStringConversion() + "\"...")

		for {
			var files_to_send []_FileInfo = nil

			files, err := os.ReadDir(to_send_dir.GPathToStringConversion())
			if nil != err {
				fmt.Println("Error reading directory \"" + to_send_dir.GPathToStringConversion() + "\".")

				goto end_loop
			}

			files_to_send = make([]_FileInfo, 0, len(files))
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".eml") {
					file_stats, _ := os.Stat(to_send_dir.Add2(file.Name()).GPathToStringConversion())
					files_to_send = append(files_to_send, _FileInfo{file.Name(), file_stats.ModTime().UnixNano()})
				}
			}

			for len(files_to_send) > 0 {
				// No mega fast email spamming - don't want the account blocked.
				time.Sleep(1 * time.Second)

				// Look for the file with the oldest modification time until there are no more files to send
				var file_to_send _FileInfo = files_to_send[0]
				var idx_to_remove int = 0
				for i := 1; i < len(files_to_send) - 1; i++ {
					if "" != files_to_send[i].file_name && files_to_send[i].modif_time < file_to_send.modif_time {
						file_to_send = files_to_send[i]
						idx_to_remove = i
					}
				}

				var file_path Utils.GPath = to_send_dir.Add2(file_to_send.file_name)

				// ... and send it.
				var mail_to string = strings.TrimSuffix(file_to_send.file_name, ".eml")
				mail_to = mail_to[Utils.RAND_STR_LEN:]

				fmt.Println("--------------------")
				fmt.Println("Sending email file " + file_to_send.file_name + " to " + mail_to + "...")

				if !reachedMaxEmailsHour() {
					if err = Utils.SendEmailEMAIL(*file_path.ReadFile(), mail_to, false); nil == err {
						if time.Now().Hour() != modGenInfo_GL.ModSpecInfo.Hour {
							modGenInfo_GL.ModSpecInfo.Hour = time.Now().Hour()
							modGenInfo_GL.ModSpecInfo.Num_emails_hour = 0
						}
						modGenInfo_GL.ModSpecInfo.Num_emails_hour++
						_ = modGenInfo_GL.Update()
						fmt.Println("Email sent successfully.")

						// Remove the file
						Utils.DelElemSLICES(&files_to_send, idx_to_remove)
						if nil == os.Remove(file_path.GPathToStringConversion()) {
							fmt.Println("File deleted successfully.")
						} else {
							fmt.Println("Error deleting file.")
						}
					} else {
						fmt.Println("Error sending email with error\n" + Utils.GetFullErrorMsgGENERAL(err))

						panic(err)
					}
				} else {
					fmt.Println("The maximum number of emails per hour has been reached (" +
						strconv.Itoa(MAX_EMAILS_HOUR) + "). Waiting for the next hour...")

					goto end_loop
				}
			}

			end_loop:

			return

			modGenInfo_GL.LoopSleep(5)
		}
	}
}

/*
reachedMaxEmailsHour returns true if the maximum number of emails per hour has been reached.

-----------------------------------------------------------

– Returns:
  - true if the maximum number of emails per hour has been reached, false otherwise.
 */
func reachedMaxEmailsHour() bool {
	return modGenInfo_GL.ModSpecInfo.Num_emails_hour >= MAX_EMAILS_HOUR &&
		time.Now().Hour() == modGenInfo_GL.ModSpecInfo.Hour
}
