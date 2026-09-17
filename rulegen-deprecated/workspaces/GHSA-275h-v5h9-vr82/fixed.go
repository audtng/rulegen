package main


		// 没有在配置文件中命中时在文件系统上查找
		filePath = filepath.Join(util.SnippetsPath, filePath)

		// 限制只能访问 snippets 目录内的文件，并拦截敏感路径，避免通过路径穿越读取工作空间内的敏感文件
		if !gulu.File.IsSubPath(util.SnippetsPath, filePath) {
			c.Status(http.StatusUnauthorized)
			return
		}
		if util.IsSensitivePath(filePath) {
			logging.LogErrorf("refuse to serve sensitive snippet file [%s]", c.Request.URL.Path)
			c.Status(http.StatusForbidden)
			return
		}

		c.File(filePath)
	})
}
