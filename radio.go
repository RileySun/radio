package radio

import (
	"time"
	"github.com/faiface/beep/speaker"
)

//Types
type Radio struct {
	Paused bool
	Song *Song
	Queue *Queue
	stopUpdating chan bool
}

type Queue struct {
	Songs [][]byte
	Original []string
	Index int64
	Length int64
}

	//Create
//Radio
func NewRadio(songByteList [][]byte) *Radio {
	radio := new(Radio)
	
	radio.Paused = false
	radio.Queue = radio.NewQueue(songByteList)
	radio.Song = NewSongFromBytes(radio.Queue.Songs[0])
	radio.stopUpdating  = make(chan bool, 100)
	
	speaker.Clear()
	
	return radio
}

//Queue
func (r *Radio) NewQueue(songList [][]byte) *Queue {
	queue := new(Queue)
	queue.Index = 0
	queue.Length = int64(len(songList))

	//Get Music Paths
	for _, songItem := range songList {
		queue.Songs = append(queue.Songs, songItem)
	}
	
	return queue
}

	//Update
func (r *Radio) startUpdate() {
	go func() {
		for {
			select {
				case <- r.stopUpdating:
					return
				default:
					speaker.Lock()
					pos := r.Song.streamer.Position()
					end := r.Song.streamer.Len()
					speaker.Unlock()
					
					if pos == end {
						r.GetQueueNext()
					}
					time.Sleep(time.Second/2)
					
	   		 }
		}
	}()
}

func (r *Radio) endUpdate() {
	r.stopUpdating <- true
}
	
//Utils
func (r *Radio) IsPlaying() bool {
	return !r.Song.Ctrl.Paused
}

func (r *Radio) Close() {
	r.Song.Close()
}

//Actions
func (r *Radio) Play() {
	r.startUpdate()
	r.Song.Play(!r.Paused)
	r.Paused = false
}

func (r *Radio) Mute() {
	speaker.Clear()
	r.endUpdate()
}

func (r *Radio) Pause() {
	r.Paused = true
	r.Song.Pause()
	r.endUpdate()
}

func (r *Radio) Stop() {
	r.Paused = false
	r.Song.Pause()
	r.endUpdate()
}

//Queue
func (r *Radio) newQueueSong(songBytes []byte) {
	r.Song.Close()
	r.Song = NewSongFromBytes(songBytes)
	r.Song.Play(true)
}

func (r *Radio) GetQueueNext() {
	//if queue exists
	if r.Queue.Length != 0 {
		if r.Queue.Index < r.Queue.Length - 1 {
			r.Queue.Index++
		} else {
			r.Queue.Index = 0
		}
		newSongBytes := r.Queue.Songs[r.Queue.Index]
		r.newQueueSong(newSongBytes)
	}
}