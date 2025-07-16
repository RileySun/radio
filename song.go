package radio

import(
	"io"
	"time"
	"bytes"
	
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

type Song struct {
	context *audio.Context
	player  *audio.Player
	
	current      time.Duration
	total        time.Duration
	
	volume128    int //128
	paused bool
	
	stopUpdating chan(bool)
	
	OnEnd func()
}

type audioStream interface {
	io.ReadSeeker
	Length() int64
}


//Context
const sampleRate = 44100
var audioContext *audio.Context
func init() {
	audioContext = audio.NewContext(sampleRate)
}

//Constructor
func NewSong(byteData []byte, format string) *Song {	
	song := &Song{
		context:audioContext,
		volume128:128,
		stopUpdating:make(chan bool),
	}
	
	//setup
	const bytesPerSample = 8
	var s audioStream
	var err error
	
	//Audio Type
	switch format {
		case "ogg":
			s, err = vorbis.DecodeF32(bytes.NewReader(byteData))
		case "mp3":
			s, err = mp3.DecodeF32(bytes.NewReader(byteData))
	}
	if err != nil {
		panic(err)
	}
	
	song.player, err = song.context.NewPlayerF32(s)
	if err != nil {
		panic(err)
	}
	
	//Props
	song.total = time.Second * time.Duration(s.Length()) / bytesPerSample / sampleRate
	if song.total == 0 {
		song.total = 1
	}
	
	return song
}

//Update
func (s *Song) startUpdate() {
	go func() {
		for {
			select {
				case <- s.stopUpdating:
					return
				default:
					if s.player.Position() >= s.total {
						if s.OnEnd != nil {
							s.OnEnd()
						}
						
						s.endUpdate()
					}
					time.Sleep(time.Second/2)	
	   		 }
		}
	}()
}

func (s *Song) endUpdate() {
	s.stopUpdating <- true
}

//Actions
func (s *Song) Play() {
	s.startUpdate()
	s.player.Play()
}

func (s *Song) PlayOnce() {
	go func() {
		s.startUpdate()
	}()
	s.player.Rewind()
	s.player.Play()
}

func (s *Song) SetVolume(newVolume float64) {
	if newVolume > 1 {
		newVolume = 1
	}
	if newVolume <= 0{
		newVolume = 0.00001
	}
	s.player.SetVolume(newVolume)
}

func (s *Song) Loop() {
	
}

func (s *Song) Pause() {
	if s.IsPlaying() {
		s.player.Pause()
		go func() {
			s.endUpdate()
		}()//Seperate thread makes it not lag
	}
}

func (s *Song) Mute() {
	s.Pause()
}

func (s *Song) Close() {
	s.player.Close()
}
	
func (s *Song) IsPlaying() bool {
	return s.player.IsPlaying()
}