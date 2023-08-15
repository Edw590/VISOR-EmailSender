package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"VISOR_S_L/Utils"
)

// Email Sender //

var modInfo_GL Utils.ModInfo = Utils.ModInfo{}
var modFileMainInfo_GL Utils.ModFileMainInfo = Utils.ModFileMainInfo{}

type _FileInfo struct {
	file_name string
	mod_time int64
}

//PUT A MAX EMAILS COUNTER!!!!!!!! 20 PER HOUR AS SAID ON GOOGLE!!!!!!!!!!!!!!!

var realMain Utils.RealMain = nil
func main() {Utils.ModStartup(Utils.NUM_MOD_EmailSender, realMain)}
func init() {realMain =
	func(realMain_param_1 Utils.ModInfo, realMain_param_2 Utils.ModFileMainInfo) {
		modInfo_GL = realMain_param_1
		modFileMainInfo_GL = realMain_param_2

		var to_send_dir Utils.GPath = modInfo_GL.Data_dir.Add(Utils.TO_SEND_REL_FOLDER)

		fmt.Println("Checking for emails to send in \"" + to_send_dir.GPathToStringConversion() + "\"...")

		for {
			files, err := os.ReadDir(to_send_dir.GPathToStringConversion())
			if nil != err {
				continue
			}

			var files_to_send []_FileInfo = make([]_FileInfo, 0, len(files))
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".eml") {
					file_stats, _ := os.Stat(to_send_dir.Add(file.Name()).GPathToStringConversion())
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
					if "" != files_to_send[i].file_name && files_to_send[i].mod_time < file_to_send.mod_time {
						file_to_send = files_to_send[i]
						idx_to_remove = i
					}
				}

				var file_path Utils.GPath = to_send_dir.Add(file_to_send.file_name)

				// ... and send it.
				var mail_to string = strings.TrimSuffix(file_to_send.file_name, ".eml")
				mail_to = mail_to[Utils.RAND_STR_LEN:]

				fmt.Println("--------------------")
				fmt.Println("Sending email file " + file_to_send.file_name + " to " + mail_to + "...")

				if Utils.UEmail.SendEmail(*file_path.ReadFile(), mail_to) {
					fmt.Println("Email sent successfully.")

					// Remove the file
					Utils.USlices.DelElem(&files_to_send, idx_to_remove)
					if nil == os.Remove(file_path.GPathToStringConversion()) {
						fmt.Println("File deleted successfully.")
					} else {
						fmt.Println("Error deleting file.")
					}
				}
			}

			os.Exit(0)

			modFileMainInfo_GL.LoopSleep(5)
		}
	}
}
