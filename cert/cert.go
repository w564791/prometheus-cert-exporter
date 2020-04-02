package cert

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/prometheus/common/log"
	//"net"
	"strings"

	//"golang.org/x/net/http2"
	"io/ioutil"
	//"strconv"

	//"math/rand"
	"os"
	"time"
)

var (
	lay_out = "2006-01-02 15:04:05"
)

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

func ParsePem(path string,

//temp map[string]string
) (lables []string, value float64) {

	if IsFile(path) {

		log.Info("accept file: ", path)
		certPEMBlock, err := ioutil.ReadFile(path)

		if err != nil {
			log.Error(err.Error())
			return
		}
		//获取证书信息 -----BEGIN CERTIFICATE-----   -----END CERTIFICATE-----
		//这里返回的第二个值是证书中剩余的 block, 一般是rsa私钥 也就是 -----BEGIN RSA PRIVATE KEY 部分
		//一般证书的有效期，组织信息等都在第一个部分里
		roots := x509.NewCertPool()
		if ok := roots.AppendCertsFromPEM(certPEMBlock); !ok {
			log.Warnln(fmt.Sprintf("%s not a eff certificate", path))
			return
		}

		certDERBlock, _ := pem.Decode(certPEMBlock)

		if certDERBlock == nil {
			log.Error("err", err.Error())
			//return
		}

		x509Cert, err := x509.ParseCertificate(certDERBlock.Bytes)
		if err != nil {
			log.Info(err.Error())

		}

		Tnow := time.Now()
		After := x509Cert.NotAfter
		From := x509Cert.NotBefore
		Domain := strings.Join(x509Cert.DNSNames, "/")
		//Ip:=net.IPNet(x509Cert.IPAddresses,"")
		//Domain:=x509Cert.do

		tn := time.Date(Tnow.Year(), Tnow.Month(), Tnow.Day(), 0, 0, 0, 0, time.Local)
		tf := time.Date(After.Year(), After.Month(), After.Day(), 0, 0, 0, 0, time.Local)

		lables = []string{After.Format(lay_out), From.Format(lay_out), path, Domain}

		return lables, tf.Sub(tn).Hours()

	}

	return
}

//func (c *ClusterManager) Describe(ch chan<- *prometheus.Desc) {
//	ch <- c.OOMCountDesc
//	//ch <- c.RAMUsageDesc
//}
//func (c *ClusterManager) Collect(ch chan<- prometheus.Metric) {
//	//oomCountByHost, ramUsageByHost := c.ReallyExpensiveAssessmentOfTheSystemState()
//	oomCountByHost:= c.ReallyExpensiveAssessmentOfTheSystemState()
//	for host, oomCount := range oomCountByHost {
//		ch <- prometheus.MustNewConstMetric(
//			c.OOMCountDesc,
//			prometheus.GaugeValue,
//			float64(oomCount),
//			host,
//		)
//	}
//	//for host, ramUsage := range ramUsageByHost {
//	//	ch <- prometheus.MustNewConstMetric(
//	//		c.RAMUsageDesc,
//	//		prometheus.GaugeValue,
//	//		ramUsage,
//	//		host,
//	//	)
//	//}
//}
//func NewClusterManager(name string) *ClusterManager {
//	return &ClusterManager{
//		Name: name,
//		OOMCountDesc: prometheus.NewDesc(
//			"cert_exp_date",
//			"certficatye expire date",
//			[]string{"from","after"},
//			prometheus.Labels{"file": name},
//		),
//		//RAMUsageDesc: prometheus.NewDesc(
//		//	"clustermanager_ram_usage_bytes",
//		//	"RAM usage as reported to the cluster manager.",
//		//	[]string{"host"},
//		//	prometheus.Labels{"Name": zone},
//		//),
//	}
//}
