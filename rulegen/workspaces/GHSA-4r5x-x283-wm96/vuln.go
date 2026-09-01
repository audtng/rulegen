package main

	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"


func startMongoDBCommand(opt *sqlOption) (lcmd *localcommand.LocalCommand, err error) {
	cmd := opt.MongoDBCommandArgs()
	lcmd, err = localcommand.New("mongosh", cmd, localcommand.WithPtyWin(opt.win.Width, opt.win.Height))
	if err != nil {
		return nil, err
	}
	"os"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/jumpserver/koko/pkg/localcommand"
func startPostgreSQLCommand(opt *sqlOption) (lcmd *localcommand.LocalCommand, err error) {
	argv := opt.PostgreSQLCommandArgs()
	//psql 是启动postgresql的客户端
	lcmd, err = localcommand.New("psql", argv, localcommand.WithPtyWin(opt.win.Width, opt.win.Height))
	if err != nil {
		return nil, err
	}
