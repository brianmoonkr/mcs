package transcode

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
	"fmt"
	"encoding/json"

	"github.com/gonutz/ftp-client/ftp"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"github.com/fsnotify/fsnotify"
	"github.com/tidwall/gjson"
	"github.com/teamgrit-lab/cojam/config"
)

type Process struct {
	Logger *zap.Logger
}

//rotationMeta ...
func (p *Process) rotationMeta(janus string, videoPath string, orgCtx context.Context) (error, string) {
	// Create a new context and add a timeout to it
	//ctx, cancel := context.WithTimeout(context.Background(), 3*time.Hour)
	ctx, cancel := context.WithTimeout(orgCtx, 3*time.Hour)
	defer cancel() // The cancel should be deferred so resources are cleaned up
	// There are good usages of context.
	///usr/local/go/src/os/exec/exec.go

	cmd := exec.CommandContext(ctx, "sh", "-c", "export JANUS_PPREC_VIDEOORIENTEXT=90;"+janus+" --parse "+videoPath) // "cmd", "/c", "xxx" on windows
	out, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		log.Println("Janus_PPREC :: ParseError: timeout")
		return nil, "transpose=0"
	}

	if err != nil {
		log.Println("Janus_PPREC :: ParseError:", err.Error())
		return nil, "transpose=0"
	}
	output := string(out)
	log.Println("Janus_PPREC :: Parse:", output)

	//https://regex101.com/r/Jw4w3s/2/
	//re := regexp.MustCompile(`(?m)Video orientation:\W?([0-9]+)\W?degrees`)
	re := regexp.MustCompile(`(?m)Video orientation extension ID:\W?([0-9]+)`)

	if re.FindString(output) == "" {
		log.Println("Janus_PPREC :: ParseError: can't find rotation regex....")
		return nil, "transpose=0"
	}

	sub := re.FindStringSubmatch(output)[1]

	i, err := strconv.Atoi(sub)
	if err != nil {
		log.Println("Janus_PPREC :: ParseError: can't convert ascii to integer:", i, sub)
		return nil, "transpose=0"
	}
	log.Println("Janus_PPREC :: Parse :: degree:", i)
	//vflip set
	switch i {
	case 90:
		return nil, "transpose=1"
	case 180:
		return nil, "transpose=2,transpose=2"
	case 270:
		return nil, "transpose=2"
	default:
		return nil, "transpose=0"
	}
}

func (p *Process) mux(avFiles []os.FileInfo, targetDir string, targetName string, orgCtx context.Context) error {
	
	januspprec, err := exec.LookPath(config.CF.Prop.JanusRecCmd)
	if err != nil {
		log.Fatalln("didn't find 'janus-pp-rec' executable:", err)
	}

	log.Println("'janus-pp-rec' executable is in:", januspprec)

	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Fatalln("Didn't find 'ffmpeg' executable:", err)
	}

	log.Println("executable 'ffmpeg' is in:", ffmpeg)

	vCmd := ""
	aCmd := ""
	mp4 := ""
	vMjr := ""
	opus := ""

	for _, f := range avFiles {
		log.Println(f.Name())
		log.Println(targetDir + "/" + f.Name())
		if strings.Contains(f.Name(), "-video.mjr") {

			vCmd = vCmd + " " + targetDir + "/" + f.Name()
			vMjr = targetDir + "/" + f.Name()
			mp4 = targetDir + "/" + strings.Replace(f.Name(), "-video.mjr", "-video.mp4", 1)
			vCmd = vCmd + " " + mp4
		}
		if strings.Contains(f.Name(), "-audio.mjr") {
			aCmd = aCmd + " " + targetDir + "/" + f.Name()
			opus = targetDir + "/" + strings.Replace(f.Name(), "-audio.mjr", "-audio.opus", 1)
			aCmd = aCmd + " " + opus
		}
	}

	err, vf := p.rotationMeta(januspprec, vMjr, orgCtx)
	if err != nil {
		log.Println("JanusVideoParse:parseRotationError:", err)
		return err
	}
	targetFullPath := targetDir + "/" + targetName
	fCmd := "-i " + mp4 + " -i " + opus + " -c:v copy -r 60 -c:a aac -g 120 -f mp4 " + targetFullPath
	tCmd := "-i " + targetFullPath + " -ss " + config.CF.Prop.ThumbOutStart + " -vcodec png -vframes 1 " + targetFullPath + "_%2d.png"

	log.Println("JanusVideoParse:parseRotation:", vf)

	log.Println(vCmd)
	log.Println(aCmd)
	log.Println(fCmd)
	log.Println(tCmd)
	
	ctx, cancel := context.WithTimeout(orgCtx, 3*time.Hour)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	//cmd := exec.CommandContext(ctx, "sh", "-c", januspprec+" "+vCmd+";"+januspprec+" "+aCmd+";"+ffmpeg+" "+fCmd) // "cmd", "/c", "xxx" on windows
	cmd := exec.CommandContext(ctx, "sh", "-c", januspprec+" "+vCmd+";"+januspprec+" "+aCmd+";"+ffmpeg+" "+fCmd+";"+ffmpeg+" "+tCmd)
	out, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		log.Println("command timed out")
		return ctx.Err()
	}

	log.Println("Output:", string(out))

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (p *Process) Process(payload string, ctx context.Context) error {
	
	seq := gjson.Get(payload, "seq").String()
	dir := gjson.Get(payload, "path").String()
	sid := gjson.Get(payload, "sid").String()
	cb := gjson.Get(payload, "cb").String()

	p.Logger.Info(dir)
	p.Logger.Info(sid)

	if sid == "" || dir == "" || len(dir) < 3 {
		p.Logger.Error("Wrong message received")
		return errors.New("Wrong message received")
	}

	
	//
	//watching
	//
	monPeriod := 1 * time.Second
	ticker := time.NewTicker(monPeriod)

	defer func() {
		ticker.Stop()
	}()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		p.Logger.Error("watcher")
		return err
	}
	defer watcher.Close()
	err = watcher.Add(dir)
	if err != nil {
		return err
	}

	cnt := NewCounter()
	cnt.Add(1)
	p.Logger.Error("begins")
	for {
		select {
		case <-ticker.C:
			p.Logger.Error("ticker")
			if cnt.Down() == 0 {
				p.Logger.Error("reject")
				break
			}

			//cnt.value()
		case event, ok := <-watcher.Events:
			if !ok {
				continue
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if path.Ext(event.Name) == ".mjr" {
					cnt.Up()
					p.Logger.Error("mjr")
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}

		if cnt.Get() == 0 {
			break
		}
	}

	if len(dir) > 1 && dir[len(dir)-1] == '/' {
		dir = dir[:len(dir)-2]
	}

	dir = dir + "/" + sid
	//files, err = filepath.Glob(dir + "/*" + sid + "*")
	files, err := ioutil.ReadDir(dir)
	if err != nil || len(files) < 1 {
		p.Logger.Error("file list")
		return err
	}

	// Move files
	targetDir := dir + "/completed"
	os.RemoveAll(targetDir)
	if err := os.MkdirAll(targetDir, 0777); nil != err {
		return err
	}

	// create live directory
	t := time.Now()
	currDate := t.Format("20060102")
	liveAccess := config.CF.Prop.SharedAccessPath + "/" + currDate
	livePath := config.CF.Prop.RepositoryPath + "/" + currDate
	if _, err = os.Stat(livePath); os.IsNotExist(err) {
		if err := os.MkdirAll(livePath, 0777); nil != err {
			return err
		}
	} 

	var largest int64 = 0;
	var largestFileName = ""

	for _, f := range files {
		if path.Ext(f.Name()) == ".mjr" {
			if f.Size() > largest {
				largest = f.Size()	
				largestFileName = f.Name()
			}
			
		}
	}

	bigVideoFileName := largestFileName
	if err := cp(dir+"/"+bigVideoFileName, targetDir+"/"+bigVideoFileName, 4096); err != nil {
		p.Logger.Error(err.Error())
		return err
	}

	bigAudioFileName := p.getAudioFileName(largestFileName)
	if err := cp(dir+"/"+bigAudioFileName, targetDir+"/"+bigAudioFileName, 4096); err != nil {
		p.Logger.Error(err.Error())
		return err
	}

	files, err = ioutil.ReadDir(targetDir)
	if err != nil || len(files) < 1 {
		p.Logger.Error("target file list")
		return err
	}
/* 위에서 최대 크기 화일 2개를 갖고 오는것으로 대체, 2020-05-08
	for _, f := range files {
		p.Logger.Info(f.Name())
		if path.Ext(f.Name()) == ".mjr" {
			if err := cp(dir+"/"+f.Name(), targetDir+"/"+f.Name(), 4096); err != nil {
				p.Logger.Error(err.Error())
				return err
			}
		}
	} 
*/	
	// mux
	if err := p.mux(files, targetDir, "result.mp4", ctx); err != nil {
		p.Logger.Error(err.Error())
		//this case would be crated witout audio , will pass to upload video only
		//return err
	}

	// Result file check
	uploadFile := targetDir + "/result.mp4"
	if _, err := os.Stat(uploadFile); os.IsNotExist(err) {
		files, err := ioutil.ReadDir(dir)
		if err != nil || len(files) < 1 {
			p.Logger.Error("file list")
			return err
		}
		for _, f := range files {
			if path.Ext(f.Name()) == ".mp4" {
				uploadFile = f.Name()
			}
		}
	}

	p.Logger.Info("Done")


	files, err = ioutil.ReadDir(targetDir)
	if err != nil || len(files) < 1 {
		p.Logger.Error("target file list")
		return err
	}

	conn, err := ftp.Connect(config.CF.Prop.CdnInfo.UploadUrl, 21)
	if err != nil {
		p.Logger.Error(err.Error())
		return err
	}
	defer conn.Close()

	err = conn.Login(config.CF.Prop.CdnInfo.FtpId, config.CF.Prop.CdnInfo.FtpPwd)
	if err != nil {
		p.Logger.Error(err.Error())
		return err
	}
	defer conn.Quit()

	var mp4RepoPath []string
	var thumbRepoPath []string

	for _, f := range files {
		p.Logger.Info(f.Name())

		ftpFilename := currDate + "_" + sid + "_" + f.Name()

		repoFilename := sid + "_" + f.Name()
		repoPath := livePath+"/"+repoFilename
		sharedAcessPath := liveAccess+"/"+repoFilename

		if f.Name() == "result.mp4" {		
			
			if err := cp(targetDir+"/"+f.Name(), repoPath, 4096); err != nil {	
				p.Logger.Error(err.Error())
				return err
			}
			//mp4RepoPath = append(mp4RepoPath, sharedAcessPath)

			source, err := os.Open(targetDir+"/"+f.Name())
			if err != nil {
				p.Logger.Error(err.Error())
				return err
			}
			defer source.Close()

			err = conn.Upload(source, "/" + ftpFilename)
        	if err != nil {
                p.Logger.Error(err.Error())
				return err
			}

			//cdnAccessUrl := config.CF.Prop.CdnInfo.CdnUrl + "/" + ftpFilename + "/playlist.m3u8"
			cdnAccessUrl := config.CF.Prop.CdnInfo.CdnUrl + "/" + ftpFilename
			
			mp4RepoPath = append(mp4RepoPath, cdnAccessUrl)

		}

		if path.Ext(f.Name()) == ".png" {
		
			if err := cp(targetDir+"/"+f.Name(), repoPath, 4096); err != nil {
				p.Logger.Error(err.Error())
				return err

			}

			thumbRepoPath = append(thumbRepoPath, sharedAcessPath)
		}

		if path.Ext(f.Name()) == ".mjr" {
			if err := os.Remove(targetDir+"/"+f.Name()); err != nil {
				panic(err)
			}	
		}
	} 
	
	//go p.callback(cb, repoPath)

	errChan := make(chan error, 1)
	
	go func() {

		mp4Path := "[]"
		thumbPath := "[]"
		if len(mp4RepoPath) > 0 {
			byteMp4Path, _ := json.Marshal(mp4RepoPath)
			mp4Path = string(byteMp4Path)
		}

		if len(thumbRepoPath) > 0 {
			byteThumbPath, _ := json.Marshal(thumbRepoPath)
			thumbPath = string(byteThumbPath)
		}
		
		err := p.callback(cb, mp4Path, thumbPath, seq)

		errChan <- err
	}()
	
	err = <-errChan

	return err

}

func (p *Process) callback(cb string, mp4Path string, thumbPath string, seq string) error {
	fmt.Println("##### callback url = " + cb)
	client := resty.New()

	cbBody := `{"mp4_path":` + mp4Path + `, "thumb_path":` + thumbPath +  `, "seq": "` + seq + `" }`
	fmt.Println("##### callback body = " + cbBody)
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(cbBody).
		Post(fmt.Sprintf("%s", cb))

	if err != nil {
		p.Logger.Error("transcode callback")
		fmt.Println(err)
		return err
	}

	log.Println("Callback After live broadcast :", cbBody)

	return nil
}

func (p *Process) getAudioFileName(vFile string) string {
	idx := strings.LastIndex(vFile, "-")
	aFile := vFile[0:idx] + "-audio.mjr"

	return aFile
}