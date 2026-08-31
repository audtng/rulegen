package main

		}

	case "Rename":
		targetPath, err := sanitizePath(r.Target, root)
		if err != nil {
			logger.LogSFTPRequestBlocked(r, ip, err)
			sftpServer.HandleWebhookSend("sftp", r, ip, true)
			return err
		}
		err = os.Rename(fullPath, targetPath)
		if err != nil {
			logger.LogSFTPRequestBlocked(r, ip, err)
			sftpServer.HandleWebhookSend("sftp", r, ip, true)
func NewSFTPServer(opts *options.Options, wl *httpserver.Whitelist, webhook webhook.Webhook) *SFTPServer {
	return &SFTPServer{
		IP:          opts.IP,
		Port:        opts.SFTPPort,
		KeyFile:     opts.SFTPKeyFile,
		Username:    opts.Username,
		Password:    opts.Password,
