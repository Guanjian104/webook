package ioc

import (
    "github.com/Guanjian104/webook/internal/service/sms"
    "github.com/Guanjian104/webook/internal/service/sms/localsms"
    "github.com/Guanjian104/webook/internal/service/sms/tencent"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
    tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
    "os"
)

func InitSMSService() sms.Service {
    return localsms.NewService()
    // return initTencentSMSService()
}

func initTencentSMSService() sms.Service {
    secretId, ok := os.LookupEnv("SMS_SECRET_ID")
    if !ok {
        panic("找不到腾讯 SMS 的 secret id")
    }
    secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
    if !ok {
        panic("找不到腾讯 SMS 的 secret key")
    }
    c, err := tencentSMS.NewClient(
        common.NewCredential(secretId, secretKey),
        "ap-nanjing",
        profile.NewClientProfile(),
    )
    if err != nil {
        panic(err)
    }
    return tencent.NewService(c, "1400842696", "妙影科技")
}
