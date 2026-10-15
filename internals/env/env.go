package env

type Env struct {
	SEED_PATH	string
	SEED		string
	CONFIG_PATH	string
	LISTEN_ADDR string
	BACKGROUND  bool
}

var ENV = &Env{}