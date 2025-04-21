package comm

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"maps"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestMustReplaceStringByEnv(t *testing.T) {
	t.Skip("skipping TestMustReplaceStringByEnv test")
	str := `aaaaadsa he ${DEPLOYMENT_NAME} fds ${kk:vv} a
dfasdf
sadf
ds is ${tp:def}a`
	result := MustReplaceStringByEnv(str)
	fmt.Println(result)
}

func TestValidatePhoneNumber(t *testing.T) {
	t.Skip("skipping TestMustReplaceStringByEnv test")

	str := "13223065856"
	result := ValidatePhoneNumber(str)
	assert.Equal(t, result, true)
}

func TestUnique(t *testing.T) {
	str := []string{"1", "2", "3", "1", "2", "3"}
	ustr := Unique(str)
	assert.Equal(t, 3, len(ustr))
}

// Job ...
type Job struct {
	Group   string `yaml:"group"`
	Id      string `yaml:"id"`
	Name    string `yaml:"name"`
	Enabled string `yaml:"scheduleEnabled"`
}

func diff(jobsA, jobsB []Job) {
	toMap := func(jobs []Job) map[string]Job {
		jobM := make(map[string]Job)
		for _, job := range jobs {
			jobM[job.Name] = job
		}
		return jobM
	}
	jmA := toMap(jobsA)
	jmB := toMap(jobsB)

	sameFile, _ := os.Create("same.csv")
	sameCsv := csv.NewWriter(sameFile)
	testFile, _ := os.Create("test.csv")
	testCsv := csv.NewWriter(testFile)
	liveFile, _ := os.Create("live.csv")
	liveCsv := csv.NewWriter(liveFile)
	for name, job := range jmA {
		if jobB, ok := jmB[name]; ok {
			if err := sameCsv.Write([]string{
				job.Name, job.Group, job.Id, job.Enabled, jobB.Enabled,
			}); err != nil {
				panic(err)
			}
		} else {
			if err := testCsv.Write([]string{
				job.Name, job.Group, job.Id, job.Enabled,
			}); err != nil {
				panic(err)
			}
		}
	}

	for name, job := range jmB {
		if _, ok := jmA[name]; !ok {
			if err := liveCsv.Write([]string{
				job.Name, job.Group, job.Id, job.Enabled,
			}); err != nil {
				panic(err)
			}

		}
	}
	sameCsv.Flush()
	testCsv.Flush()
	liveCsv.Flush()
	sameFile.Close()
	testFile.Close()
	liveFile.Close()
}

func TestParse(t *testing.T) {
	jobsFile, err := ioutil.ReadFile("jobT.yaml")
	assert.Nil(t, err)
	var jobsT, jobsL []Job
	yaml.Unmarshal(jobsFile, &jobsT)
	jobsFile, err = ioutil.ReadFile("jobL.yaml")
	yaml.Unmarshal(jobsFile, &jobsL)
	diff(jobsT, jobsL)
}

func TestS3(t *testing.T) {
	assert.Nil(t, list())
}

func TestSelect(t *testing.T) {
	var ch chan int
	realCh := make(chan int, 1000)
	for {
		select {
		case ch <- 1:
			t.Log("case 1")
			fmt.Println("case 1")
			ch = nil
		default:
			t.Log("default")
			ch = realCh
		}
		time.Sleep(time.Second)
	}
}

func TestMaps(t *testing.T) {
	m1 := map[int]string{
		1000: "THOUSAND",
	}
	s1 := []string{"zero", "one", "two", "three"}
	maps.Insert(m1, slices.All(s1))
	fmt.Println("m1 is:", m1)
	min(1, 2)
}
