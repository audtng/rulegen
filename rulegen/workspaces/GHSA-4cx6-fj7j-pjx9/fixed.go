package main

	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	exec "golang.org/x/sys/execabs"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"

import (
	"fmt"
	"runtime"

	exec "golang.org/x/sys/execabs"
)

var execCommand = exec.Command
