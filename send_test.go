package journald

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	uniqueSmallMessageID   = "1"
	uniqueBigMessageID     = "2"
	uniquePrintedMessageID = "3"
	uniqueLowerMessageID   = "4"
)

func TestMain(m *testing.M) {
	uniqueSmallMessageID = time.Now().String()
	time.Sleep(time.Millisecond)
	uniqueBigMessageID = time.Now().String()
	time.Sleep(time.Millisecond)
	uniquePrintedMessageID = time.Now().String()
	time.Sleep(time.Millisecond)
	uniqueLowerMessageID = time.Now().String()

	os.Exit(m.Run())
}

func TestIsNotExist(t *testing.T) {
	require.False(t, IsNotExist())
}

type Test struct {
	Field1 string
	Field2 int
	Field3 map[string]float32
}

func TestSendSmallMessage(t *testing.T) {
	s := Test{"field1", 2, map[string]float32{"3": 0.4, "5": 6.7}}
	trace := runtime.ReadTrace()
	err := Send("AnotherMessage", PriorityInfo, map[string]interface{}{"ONE": 1, "TWO": trace, "THREE": s, "FOUR": RandStringRunes(768)})
	require.NoError(t, err)
}

func TestCheckPrint(t *testing.T) {
	err := Print(PriorityNotice, "printed message: %s", uniquePrintedMessageID)
	require.NoError(t, err)

	time.Sleep(time.Second * 5)

	out, err := exec.Command("sh", "-c", fmt.Sprintf("journalctl 'MESSAGE=printed message: %s' -o verbose", uniquePrintedMessageID)).Output()
	require.NoError(t, err)
	require.True(t, strings.Contains(string(out), fmt.Sprintf("MESSAGE=printed message: %s", uniquePrintedMessageID)))
	require.True(t, strings.Contains(string(out), "PRIORITY=5"))
}

func TestCheckSmallMessage(t *testing.T) {
	err := Send("SmallMessage", PriorityWarning, map[string]interface{}{"TEST_ID": uniqueSmallMessageID, "WITH_NEWLINES": "b\n\nsd\n"})
	require.NoError(t, err)

	time.Sleep(time.Second * 5)

	out, err := exec.Command("sh", "-c", fmt.Sprintf("journalctl 'TEST_ID=%s' -o verbose", uniqueSmallMessageID)).Output()
	require.NoError(t, err)
	require.True(t, len(out) > 500)
	require.True(t, strings.Contains(string(out), "WITH_NEWLINES"))
	require.True(t, strings.Contains(string(out), "MESSAGE=SmallMessage"))
	require.True(t, strings.Contains(string(out), "PRIORITY=4"))

	out, err = exec.Command("sh", "-c", fmt.Sprintf("journalctl -o verbose 'TEST_ID=%s' -F WITH_NEWLINES", uniqueSmallMessageID)).Output()
	require.NoError(t, err)
	require.Equal(t, "b\n\nsd\n"+"\n", string(out))
}

func TestCheckBigMessage(t *testing.T) {
	err := Send("BigMessage", PriorityErr, map[string]interface{}{"TEST_ID": uniqueBigMessageID, "BIG_MESSAGE": RandStringRunes(1024 * 512)})
	require.NoError(t, err)

	time.Sleep(time.Second * 5)

	out, err := exec.Command("sh", "-c", fmt.Sprintf("journalctl 'TEST_ID=%s' -o verbose", uniqueBigMessageID)).Output()
	require.NoError(t, err)
	require.True(t, len(out) > 512*1024)
	require.True(t, strings.Contains(string(out), "BIG_MESSAGE"))
	require.True(t, strings.Contains(string(out), "MESSAGE=BigMessage"))
	require.True(t, strings.Contains(string(out), "PRIORITY=3"))
}

func TestCheckLowerMessage(t *testing.T) {
	j := Journal{}
	j.NormalizeFieldNameFn = strings.ToUpper
	err := j.Send("LowerMessage", PriorityCrit, map[string]interface{}{"test_id": uniqueLowerMessageID})
	require.NoError(t, err)
	require.NoError(t, j.Close())

	time.Sleep(time.Second * 5)

	out, err := exec.Command("sh", "-c", fmt.Sprintf("journalctl 'TEST_ID=%s' -o verbose", uniqueLowerMessageID)).Output()
	require.NoError(t, err)
	require.True(t, strings.Contains(string(out), "MESSAGE=LowerMessage"))
	require.True(t, strings.Contains(string(out), "PRIORITY=2"))
}

func TestSendClose(t *testing.T) {
	j := Journal{}
	err := j.Print(PriorityInfo, "SendClose")
	require.NoError(t, err)
	require.NoError(t, j.Close())
	err = j.Print(PriorityInfo, "SendClose")
	assertErr(t, err, " use of closed network connection")
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func assertErr(t *testing.T, e error, rx interface{}) {
	assert.Error(t, e)
	assert.Regexp(t, rx, e.Error())
}
