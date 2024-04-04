package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/smtp"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/ivpusic/grpool"
)

type stats struct {
	Durations  []map[string]time.Duration
	SrcIPStats []string
	DstIPStats []string
	ErrorCnt   int
	TotalCnt   int
}

var (
	metric                              stats
	host, from, to, subject, body, helo string
	workers, count, jobs, size, timeout int
	balance, showerror, journal         bool
	outbound                            string
)

const (
	metricTemplate = `` +
		`Bomberman - SMTP Performance Test Tool` + "\n" +
		`--------------------------------------` + "\n" +
		`Count			: %d` + "\n" +
		`Error			: %d` + "\n" +
		`Size			: %dK` + "\n" +
		`Start			: %v` + "\n" +
		`End			: %v` + "\n" +
		`Time			: %v` + "\n"

	bodyTemplate = `Message-ID: <%s.FakeMail.Golang@HappyDay>` + "\r\n" +
		`Date: <%s>` + "\r\n" +
		`From: <%s>` + "\r\n" +
		`To: %s` + "\r\n" +
		`Subject: %s` + "\r\n\r\n" +
		`%s`

	bodyJournalTemplate = `Content-Type: multipart/mixed; boundary="_%s_"` + "\r\n" +
		`MIME-Version: 1.0` + "\r\n" +
		`Message-ID: <%s.FakeJournalMail.Golang@HappyDay>` + "\r\n" +
		`Date: <%s>` + "\r\n" +
		`From: <%s>` + "\r\n" +
		`To: %s` + "\r\n" +
		`Subject: %s` + "\r\n\r\n" +
		`--_%s_` + "\r\n" +
		`Content-Type: text/plain; charset="us-ascii"` + "\r\n" +
		`Content-Transfer-Encoding: 7bit` + "\r\n" +
		`MIME-Version: 1.0` + "\r\n\r\n" +
		`%s` + "\r\n" + // plain text body
		`--_%s_` + "\r\n" +
		`Content-Type: message/rfc822` + "\r\n\r\n" +
		`MIME-Version: 1.0` + "\r\n" +
		`Message-ID: <%s.FakeRedirectedMail.Golang@HappyDay>` + "\r\n" +
		`Date: <%s>` + "\r\n" +
		`From: <%s>` + "\r\n" +
		`To: %s` + "\r\n" +
		`Subject: %s` + "\r\n\r\n" +
		`Content-Type: text/html; charset=utf-8` + "\r\n" +
		`Content-Transfer-Encoding: base64` + "\r\n" +
		`Content-Disposition: attachment; filename= PDF_200KB.pdf` + "\r\n\r\n" +
		`%s` + "\r\n" + // base64 encoded pdf
		`--_%s_--` + "\r\n"

	dialTimeout   = time.Second * 6
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func init() {

	flag.StringVar(&host, "host", "localhost:25", "-host=example.org:25")
	flag.StringVar(&from, "from", "me@example.org", "-from=me@example.org")
	flag.StringVar(&to, "to", "to@example.net", "-to=me@example.net")
	flag.StringVar(&subject, "subject", "Test Email", "-subject=Test Email")
	flag.StringVar(&helo, "helo", "mail.example.org", "-helo=mail.example.org")
	flag.StringVar(&outbound, "outbound", "", "-outbound=0.0.0.0")
	flag.IntVar(&timeout, "timeout", 6, "-timeout=6 (second)")
	flag.IntVar(&count, "count", 10, "-count=10")
	flag.IntVar(&workers, "workers", 10, "-workers=100")
	flag.IntVar(&jobs, "jobs", 10, "-jobs=50")
	flag.IntVar(&size, "size", 5, "size=5 (Kilobyte)")
	flag.BoolVar(&balance, "balance", false, "-balance")
	flag.BoolVar(&showerror, "showerror", true, "-showerror")
	flag.BoolVar(&journal, "journal", false, "-journal to emulate journaling (email inside email)")
	//TODO: timeout

	flag.Usage = usage

	metric = stats{
		Durations:  []map[string]time.Duration{},
		SrcIPStats: []string{},
		DstIPStats: []string{},
	}
}

func usage() {

	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "OPTIONS:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "USAGE:")
	fmt.Fprintln(os.Stderr, "./bomberman -host=mail.server.com:25 -from=test@mydomain.com -to=user@remotedomain.com -workers=100 -jobs=100 -count=100 -helo=mydomain.com -balance -size=2")
	fmt.Fprintln(os.Stderr, "")
}

func main() {

	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)

	flag.Parse()

	startTime, endtime := start()

	printResults(balance, startTime, endtime)
}

func printResults(balanced bool, startTime, endtime time.Time) {

	fmt.Printf(metricTemplate,
		metric.TotalCnt,
		metric.ErrorCnt,
		size,
		startTime,
		endtime,
		endtime.Sub(startTime))

	if balanced {
		fmt.Println("")
		fmt.Println("Source IP Stats:")
		fmt.Println("")
		printSlice(metric.SrcIPStats, "%s\t\t: %d\n")
	}

	if len(metric.DstIPStats) > 1 {
		fmt.Println("")
		fmt.Println("Destination IP Stats:")
		fmt.Println("")
		printSlice(metric.DstIPStats, "%s\t: %d\n")
	}

	fmt.Println("")
	fmt.Println("SMTP Commands:")
	fmt.Println("")

	mkeys := metricKeys(metric.Durations)

	for i := 0; i < len(mkeys); i++ {
		m := mkeys[i]
		min, max, me := getMetric(m, metric.Durations)
		cnt := countMetric(m, metric.Durations)
		fmt.Printf("%s (%d)\t: min. %-15v, max. %-15v, med. %-15v\n", m, cnt, min, max, me)
	}
	fmt.Println("")
}

func start() (startTime time.Time, endTime time.Time) {

	pool := grpool.NewPool(workers, jobs)
	defer pool.Release()
	pool.WaitCount(count)

	iplist, err := ipv4list()

	if err != nil {
		log.Fatal("pool not created:", err)
	}

	if !journal {
		body = createBodyFixedSize(size)
	} else {
		var errPDF error
		body, errPDF = createPDFWithTargetSize(size)
		if errPDF != nil {
			log.Fatal("PDF not created:", errPDF)
		}
	}

	uuids := prepareUUIDs(count)
	emails := prepareFakeEmails(count, strings.Split(from, "@")[1])
	subjects := prepareFakeSubjects(count)
	messageDate := time.Now().Format(time.RFC1123Z)

	startTime = time.Now()
	for i := 0; i < count; i++ {

		if balance {
			outbound = sequental(i, iplist)
			metric.SrcIPStats = append(metric.SrcIPStats, outbound)
		}

		iLocal := i

		pool.JobQueue <- func() {

			metric.TotalCnt++

			defer pool.JobDone()

			durs, remoteip, err := sendMail(outbound,
				host,
				emails[iLocal],
				to,
				subjects[iLocal],
				body,
				helo,
				uuids[iLocal],
				messageDate)

			if err != nil {
				if showerror {
					fmt.Printf("%d: %v\n", metric.TotalCnt, err)
				}
				metric.ErrorCnt++
			}

			if remoteip != "" {
				metric.DstIPStats = append(metric.DstIPStats, remoteip)
			}

			metric.Durations = append(metric.Durations, durs)
		}
	}

	pool.WaitAll()
	endTime = time.Now()

	return
}

// sendMail sends an email using the provided SMTP server.
// It takes the following parameters:
// - outbound: The outbound IP address or hostname.
// - smtpServer: The SMTP server address.
// - from: The email address of the sender.
// - to: The email address of the recipient.
// - subject: The subject of the email.
// - body: The body of the email.
// - helo: The HELO/EHLO hostname.
// It returns the following values:
// - metric: A map containing the duration of each step in the email sending process.
// - remoteip: The IP address of the remote SMTP server.
// - err: An error, if any occurred during the email sending process.
func sendMail(outbound, smtpServer, from, to, subject, body, helo string, messageID string, messageDate string) (metric map[string]time.Duration, remoteip string, err error) {

	var wc io.WriteCloser
	var msg string

	startTime := time.Now()

	metric = map[string]time.Duration{}
	host, _, _ := net.SplitHostPort(smtpServer)
	conn, err := newDialer(outbound, smtpServer, dialTimeout)

	if err != nil {
		err = fmt.Errorf("DIAL: %v (out:%s)", err, outbound)
		metric["DIAL"] = time.Since(startTime)
		return
	}

	remoteip = conn.RemoteAddr().String() //remoteip
	metric["DIAL"] = time.Since(startTime)

	newclientTime := time.Now()
	c, err := smtp.NewClient(conn, host)

	if err != nil {
		err = fmt.Errorf("TOUCH: %v", err)
		metric["TOUCH"] = time.Since(newclientTime)
		return
	}

	metric["TOUCH"] = time.Since(newclientTime)
	defer c.Close()

	helloTime := time.Now()
	err = c.Hello(helo)

	if err != nil {
		err = fmt.Errorf("HELO: %v", err)
		metric["HELO"] = time.Since(helloTime)

		return
	}

	metric["HELO"] = time.Since(helloTime)

	mailTime := time.Now()
	err = c.Mail(from)

	if err != nil {
		err = fmt.Errorf("MAIL: %v", err)
		metric["MAIL"] = time.Since(mailTime)

		return
	}

	metric["MAIL"] = time.Since(mailTime)

	rcptTime := time.Now()
	err = c.Rcpt(to)

	if err != nil {
		err = fmt.Errorf("RCPT: %v", err)
		metric["RCPT"] = time.Since(rcptTime)

		return
	}

	metric["RCPT"] = time.Since(rcptTime)

	dataTime := time.Now()

	// template the message body with the message ID and other parameters
	if !journal {
		msg = fmt.Sprintf(bodyTemplate, messageID, messageDate, from, to, subject, body)
	} else {
		msg = fmt.Sprintf(bodyJournalTemplate, messageID, messageID, messageDate, from, to, subject, messageID, "Test email, sorry, no message to analyze", messageID, messageID, messageDate, from, to, subject, body, messageID)
	}

	wc, err = c.Data()

	if err != nil {
		err = fmt.Errorf("DATA: %v", err)
		metric["DATA"] = time.Since(dataTime)

		return
	}

	_, err = fmt.Fprintf(wc, msg)

	err = wc.Close()

	if err != nil {
		err = fmt.Errorf("DATA: %v", err)
		metric["DATA"] = time.Since(dataTime)

		return
	}

	metric["DATA"] = time.Since(dataTime)

	quitTime := time.Now()
	err = c.Quit()

	if err != nil {
		err = fmt.Errorf("QUIT: %v", err)
		metric["QUIT"] = time.Since(quitTime)

		return
	}

	metric["SUCCESS"] = time.Since(startTime)

	return
}

func getMetric(name string, metrics []map[string]time.Duration) (max, min, med time.Duration) {

	totaldur, _ := time.ParseDuration("0ms")
	list := []time.Duration{}

	for i := 0; i < len(metrics); i++ {
		m := metrics[i]

		if t, ok := m[name]; ok {
			totaldur += t
			list = append(list, t)
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i] > list[j]
	})

	min = list[0]
	max = list[len(list)-1]
	med = totaldur / time.Duration(len(list))

	return
}

func countMetric(name string, metrics []map[string]time.Duration) (cnt int) {

	for i := 0; i < len(metrics); i++ {
		m := metrics[i]

		for mkey := range m {

			if mkey == name {
				cnt++
			}
		}
	}

	return
}

func metricKeys(metrics []map[string]time.Duration) (keys []string) {

	for i := 0; i < len(metrics); i++ {
		m := metrics[i]

		for mkey := range m {
			contain := isContain(mkey, keys)

			if !contain {
				keys = append(keys, mkey)
			}
		}
	}

	sort.Strings(keys)

	return
}

func isContain(key string, keys []string) bool {

	exists := false

	for z := 0; z < len(keys); z++ {
		if keys[z] == key {
			exists = true
			break
		}
	}

	return exists
}

func ipv4list() (iplist []string, err error) {

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return
	}

	for i := 0; i < len(addrs); i++ {

		addr := addrs[i].String()

		if strings.Contains(addr, ":") {
			continue
		}

		nt, _, err := net.ParseCIDR(addr)

		if err != nil {
			continue
		}

		if !nt.IsGlobalUnicast() {
			continue
		}

		iplist = append(iplist, nt.String())
	}

	return
}

func newDialer(outboundip, remotehost string, timeout time.Duration) (conn net.Conn, err error) {

	if outboundip == "" {
		return net.Dial("tcp", remotehost)
	}

	dialer := &net.Dialer{Timeout: timeout}
	dialer.LocalAddr = &net.TCPAddr{IP: net.ParseIP(outboundip)}

	conn, err = dialer.Dial("tcp", remotehost)

	return
}

func sequental(index int, list []string) string {

	var ob string
	ln := len(list)

	if index < ln {
		ob = list[index]
	} else {
		li := index % ln
		ob = list[li]
	}

	return ob
}

func printSlice(list []string, format string) {

	m := map[string]int{}

	for i := 0; i < len(list); i++ {
		item := list[i]

		if item == "" {
			continue
		}

		if _, ok := m[item]; !ok {
			m[item] = 1
		} else {
			m[item] = m[item] + 1
		}
	}

	for k, v := range m {
		fmt.Printf(format, k, v)
	}
}

func createBodyFixedSize(n int) string {

	n = n * 1024

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)

	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
