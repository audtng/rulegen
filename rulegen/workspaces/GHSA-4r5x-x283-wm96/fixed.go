package main

package srvconn

import (
	"os/user"
	"strconv"
	"syscall"

	"github.com/jumpserver/koko/pkg/localcommand"
)

func BuildNobodyWithOpts(opts ...localcommand.Option) (nobodyOpts []localcommand.Option, err error) {
	nobody, err := user.Lookup("nobody")
	if err != nil {
		return nil, err
	}
	uid, _ := strconv.Atoi(nobody.Uid)
	gid, _ := strconv.Atoi(nobody.Gid)
	nobodyOpts = make([]localcommand.Option, 0, len(opts)+1)
	nobodyOpts = append(nobodyOpts, opts...)
	nobodyCredential := localcommand.WithCmdCredential(&syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)})
	nobodyOpts = append(nobodyOpts, nobodyCredential)
	return nobodyOpts, nil
}
	"strconv"
	"time"

	"github.com/jumpserver/koko/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"


func startMongoDBCommand(opt *sqlOption) (lcmd *localcommand.LocalCommand, err error) {
	cmd := opt.MongoDBCommandArgs()
	opts, err := BuildNobodyWithOpts(localcommand.WithPtyWin(opt.win.Width, opt.win.Height))
	if err != nil {
		logger.Errorf("build nobody with opts error: %s", err)
		return nil, err
	}
	lcmd, err = localcommand.New("mongosh", cmd, opts...)
	if err != nil {
		return nil, err
	}
	"os"
	"strconv"

	"github.com/jumpserver/koko/pkg/logger"
	_ "github.com/lib/pq"

	"github.com/jumpserver/koko/pkg/localcommand"
func startPostgreSQLCommand(opt *sqlOption) (lcmd *localcommand.LocalCommand, err error) {
	argv := opt.PostgreSQLCommandArgs()
	//psql 是启动postgresql的客户端
	opts, err := BuildNobodyWithOpts(localcommand.WithPtyWin(opt.win.Width, opt.win.Height))
	if err != nil {
		logger.Errorf("build nobody with opts error: %s", err)
		return nil, err
	}
	lcmd, err = localcommand.New("psql", argv, opts...)
	if err != nil {
		return nil, err
	}
