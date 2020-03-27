package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v1"
	"io/ioutil"
	//"log"
	"net/http"
	"os"
	"time"
)

var (
	Paths = kingpin.Flag("path","cert path,provide file/dir").Required().Strings()
	CertExp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "cert_exp_date",
			Help:      "Number of blob storage operations waiting to be processed, partitioned by user and type.",
		},
		[]string{
			"name",
			"from",
			"after",
		},
	)
	filesA=make(map[string]string)
)
// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}
func parsePemFile(ps []string)  {
	var ss=make(map[string]string)
	for _,p :=range ps{
		parsePem(p,ss)

	}

}

func parsePem(path string,temp map[string]string) {
	//var filesM []string

	if IsFile(path){

		log.Info("found file",path)
		certPEMBlock, err := ioutil.ReadFile(path)

		if err != nil {
			log.Error(err.Error())
			return
		}
		//获取证书信息 -----BEGIN CERTIFICATE-----   -----END CERTIFICATE-----
		//这里返回的第二个值是证书中剩余的 block, 一般是rsa私钥 也就是 -----BEGIN RSA PRIVATE KEY 部分
		//一般证书的有效期，组织信息等都在第一个部分里
		roots := x509.NewCertPool()
		if ok := roots.AppendCertsFromPEM(certPEMBlock);!ok{
			log.Warnln(fmt.Sprintf("%s not a eff certificate",path))
			return
		}
		temp[path] = "true"

		filesA[path]="true"

		certDERBlock, _ := pem.Decode(certPEMBlock)

		if certDERBlock == nil {
			log.Error("err",err.Error())
			//return
		}
		//layout:="2006-01-02 15:04:05"

		x509Cert, err := x509.ParseCertificate(certDERBlock.Bytes)
		if err != nil {
			log.Info(err.Error())

		}

		Tnow:=time.Now()
		After:=x509Cert.NotAfter
		From:=x509Cert.NotBefore

		tn:=time.Date(Tnow.Year(),Tnow.Month(),Tnow.Day(),0,0,0,0,time.Local)
		tf:=time.Date(After.Year(),After.Month(),After.Day(),0,0,0,0,time.Local)
		//CertExp.With(prometheus.Labels{
		//	"name": path,
		//	"from":From.String(),
		//	"after":  After.String(),
		//}).Inc()

		CertExp.WithLabelValues(path,From.String(),After.String()).Set(tf.Sub(tn).Hours())
		//CertExp.WithLabelValues(path,From.String(),After.String())
		//ok:=CertExp.Delete(map[string]string{"name":"ca3.pem"})
		//fmt.Println(ok)
		for k,_:=range filesA{
			if _,ok:=temp[k];!ok{
				CertExp.WithLabelValues(k,From.String(),After.String()).Set(-1)
			}
		}


		//CertExp.Delete("name":"")
	}else if IsDir(path){
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Error(err.Error())
		}
		for _,file :=range files{
			parsePem(fmt.Sprintf("%s/%s",path,file.Name()),temp)

		}
	}


}
func init() {
	prometheus.MustRegister(CertExp)
}

func main(){

	kingpin.Parse()

	go func() {
		for   {

			parsePemFile(*Paths)
			time.Sleep(time.Duration(10 *time.Second))
		}
	}()
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	log.Fatal(http.ListenAndServe(":9090",nil))


}
