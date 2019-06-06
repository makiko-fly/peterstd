package env

import "os"

//// 在不同的运行环境中，会将环境变量CONFIGOR设置为不同的值
//// 根据CONFIGOR的值来加载不同的配置文件
//// 配置文件命名格式应该为：{service_name}.{config_name}.yaml
//// local:      本地开发环境
//// testing:    线上测试环境
//// staging:    灰度测试环境
//// production: 正式生产环境
//// See: https://en.wikipedia.org/wiki/Deployment_environment
//
//func IsLocalEnv() bool {
//	return os.Getenv("CONFIGOR_ENV") == "local"
//}
//
//func IsTestingEnv() bool {
//	return os.Getenv("CONFIGOR_ENV") == "testing"
//}
//
//func IsStagingEnv() bool {
//	return os.Getenv("CONFIGOR_ENV") == "staging"
//}
//
//func IsProductionEnv() bool {
//	return os.Getenv("CONFIGOR_ENV") == "production"
//}

// use ivanka standard

func IsLocalEnv() bool {
	return os.Getenv("CONFIGOR_ENV") == "local"
}

func IsTestingEnv() bool {
	return os.Getenv("CONFIGOR_ENV") == "ivktest"
}

func IsProductionEnv() bool {
	return os.Getenv("CONFIGOR_ENV") == "ivkprod"
}
