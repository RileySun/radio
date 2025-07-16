package radio

import (
	
)

//Types
type Radio struct {
	Paused bool
	Volume float64
	Song *Song
	Queue *Queue
	stopUpdating chan bool
	Format string
}

type Queue struct {
	Songs [][]byte
	Original []string
	Index int64
	Length int64
}

	//Create
//Radio
func NewRadio(songByteList [][]byte, format string) *Radio {
	newRadio := &Radio {
		Paused:false,
		Format:format,
		Volume:1,
	}
	
	newRadio.Queue = newRadio.NewQueue(songByteList)
	newRadio.Song = NewSong(newRadio.Queue.Songs[0], newRadio.Format)
	newRadio.Song.OnEnd = newRadio.GetQueueNext
	
	return newRadio
}

//Queue
func (r *Radio) NewQueue(songList [][]byte) *Queue {
	queue := &Queue {
		Index:0,
		Length:int64(len(songList)),
	}

	//Get Music Paths
	for _, songItem := range songList {
		queue.Songs = append(queue.Songs, songItem)
	}
	
	return queue
}
	
//Utils
func (r *Radio) IsPlaying() bool {
	return r.Song.IsPlaying()
}

func (r *Radio) Close() {
	r.Song.Close()
}

//Actions
func (r *Radio) Play() {
	r.Song.Play()
	r.Paused = false
}

func (r *Radio) Mute() {
	r.Song.Mute()
}

func (r *Radio) Pause() {
	if r.IsPlaying() {
		r.Paused = true
		r.Song.Pause()
	}
}

func (r *Radio) Stop() {
	r.Paused = false
	r.Song.Pause()
}

func (r *Radio) SetVolume(newVolume float64) {
	if newVolume > 1 {
		newVolume = 1
	}
	if newVolume <= 0 {
		newVolume = 0.0001
		r.Paused = true
	}
	
	r.Volume = newVolume
	r.Song.SetVolume(r.Volume)
}

//Queue
func (r *Radio) GetQueueNext() {
	//if queue exists
	if r.Queue.Length != 0 {
		if r.Queue.Index < r.Queue.Length - 1 {
			r.Queue.Index++
		} else {
			r.Queue.Index = 0
		}
		songBytes := r.Queue.Songs[r.Queue.Index]
		r.Song.Close()
		r.Song = NewSong(songBytes, r.Format)
		r.Song.SetVolume(r.Volume)
		r.Song.OnEnd = r.GetQueueNext
		r.Song.Play()
	}
}