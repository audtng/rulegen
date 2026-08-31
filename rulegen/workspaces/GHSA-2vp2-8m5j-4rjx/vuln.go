package main

			if _, err := copyFile(filepath.Join(common.GetAppPath(), "nps"), "/usr/local/bin/nps"); err != nil {
				log.Fatalln(err)
			} else {
				os.Chmod("/usr/local/bin/nps", 0777)
				log.Println("Executable files have been copied to", "/usr/local/bin/nps")
			}
		} else {
			os.Chmod("/usr/bin/nps", 0777)
			log.Println("Executable files have been copied to", "/usr/bin/nps")
		}

