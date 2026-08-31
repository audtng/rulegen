package main

		}

	case "Rename":
		err := os.Rename(fullPath, r.Target)
		if err != nil {
			logger.LogSFTPRequestBlocked(r, ip, err)
			sftpServer.HandleWebhookSend("sftp", r, ip, true)
func NewSFTPServer(opts *options.Options, wl *httpserver.Whitelist, webhook webhook.Webhook) *SFTPServer {
	return &SFTPServer{
		IP:          opts.IP,
		Port:        opts.Port,
		KeyFile:     opts.SFTPKeyFile,
		Username:    opts.Username,
		Password:    opts.Password,
