package radio

import(
	"log"
	"os"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type Song struct {
	Path string
	file *os.File
	streamer beep.StreamSeekCloser
	format beep.Format
	Current int
	Length int
	Ctrl *beep.Ctrl
}

var oldSampleRate beep.SampleRate

func init() {
	oldSampleRate = 44100
	speaker.Init(oldSampleRate, 1)
}

//Constructor
func NewSong(filepath string) *Song {	
	currentSong := new(Song)
	
	currentSong.Path = filepath
	file := OpenSongFile(filepath)
	
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

func (s *Song) IsEnded() bool {
	return s.Current >= s.Length
}

//Actions
func (s *Song) Play() {
	_ = s.streamer.Seek(0)
	resampled := Resample(s)
	
	s.Ctrl = &beep.Ctrl{Streamer: resampled}
	speaker.Play(resampled)
}

func (s *Song) Mute() {
	speaker.Clear()
}
	
func (s *Song) Close() {
	s.streamer.Close()
}