package util

func HandleErr(err error, tips string) {
	if err != nil {
		panic(tips + " 错误信息:" + err.Error())
	}
}
