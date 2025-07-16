package radio

import(
	"os"
	"log"
	
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

type Song struct {
	Path string
	file *os.File
	streamer beep.StreamSeekCloser
	format beep.Format
	Ctrl *beep.Ctrl
	
	OnEnd func()
}

var oldSampleRate beep.SampleRate

func init() {
	oldSampleRate = 44100
	speaker.Init(oldSampleRate, 1)
}

//Constructor
func NewSong(filepath string) *Song {	
	currentSong := &Song {
		Path:filepath,
	}
	
	file := OpenSongFile(currentSong.Path)
	
	var err error
	currentSong.streamer, currentSong.format, err = mp3.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	
	return currentSong
}

func NewSongFromBytes(byteData []byte) *Song {	
	currentSong := &Song{}
	
	file := NewVirtualFile()
	_, _ = file.Write(byteData)
	
	var err error
	currentSong.streamer, currentSong.format, err = mp3.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	
	return currentSong
}

//Utils
func OpenSongFile(filepath string) *os.File {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func Resample(s *Song) *beep.Resampler  {
	return beep.Resample(1, oldSampleRate, s.format.SampleRate, s.streamer)
}

//Actions
func (s *Song) Play(restart bool) {
	if restart {
		seekErr := s.streamer.Seek(0)
		if seekErr != nil {
			log.Fatal(seekErr)
		}
	}
	resampled := Resample(s)
	
	s.Ctrl = &beep.Ctrl{Streamer: resampled}
	speaker.Play(beep.Seq(s.Ctrl, beep.Callback(func() {
		if s.OnEnd != nil {
			s.OnEnd()
		}
	})))
}

func (s *Song) PlayOnce() {
	seekErr := s.streamer.Seek(0)
	if seekErr != nil {
		log.Fatal(seekErr)
	}
	resampled := Resample(s)
	
	s.Ctrl = &beep.Ctrl{Streamer: resampled}
	speaker.Play(beep.Seq(s.Ctrl, beep.Callback(func() {
		if s.OnEnd != nil {
			s.OnEnd()
		}
	})))
}

func (s *Song) Loop() {
	_ = s.streamer.Seek(0)
	loop := beep.Loop(-1, s.streamer)
	s.Ctrl = &beep.Ctrl{Streamer: loop}
	speaker.Play(s.Ctrl)
}

func (s *Song) Pause() {
	speaker.Lock()
	s.Ctrl.Paused = true
	speaker.Unlock()
}

func (s *Song) Mute() {
	speaker.Clear()
}
	
func (s *Song) Close() {
	s.streamer.Close()
}