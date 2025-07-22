package radio

import (
	"log"
	"embed"
)

//Types
type EmbedRadio struct {
	Paused bool
	Volume float64
	Song *Song
	Queue *EmbedQueue
	stopUpdating chan bool
	Format string
	embedFiles embed.FS
}

type EmbedQueue struct {
	Songs []string
	Index int64
	Length int64
}

	//Create
//EmbedRadio
func NewEmbedRadio(songList []string, format string, embedFiles embed.FS) *EmbedRadio {
	embedRadio := &EmbedRadio {
		Paused:false,
		Format:format,
		Volume:1,
		embedFiles:embedFiles,
	}
	
	embedRadio.Queue = embedRadio.NewEmbedQueue(songList)
	embedRadio.Song = NewSong(embedRadio.getFile(songList[0]), embedRadio.Format)
	embedRadio.Song.OnEnd = embedRadio.GetEmbedQueueNext
	
	return embedRadio
}

//EmbedQueue
func (r *EmbedRadio) NewEmbedQueue(songList []string) *EmbedQueue {
	queue := &EmbedQueue {
		Index:0,
		Length:int64(len(songList)),
		Songs:songList,
	}
	
	return queue
}
	
//Utils
func (r *EmbedRadio) IsPlaying() bool {
	return r.Song.IsPlaying()
}

func (r *EmbedRadio) Close() {
	r.Song.Close()
}

//Actions
func (r *EmbedRadio) Play() {
	r.Song.Play()
	r.Paused = false
}

func (r *EmbedRadio) Mute() {
	r.Song.Mute()
}

func (r *EmbedRadio) Pause() {
	if r.IsPlaying() {
		r.Paused = true
		r.Song.Pause()
	}
}

func (r *EmbedRadio) Stop() {
	r.Paused = false
	r.Song.Pause()
}

func (r *EmbedRadio) SetVolume(newVolume float64) {
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

//EmbedQueue
func (r *EmbedRadio) GetEmbedQueueNext() {
	//if queue exists
	if r.Queue.Length != 0 {
		if r.Queue.Index < r.Queue.Length - 1 {
			r.Queue.Index++
		} else {
			r.Queue.Index = 0
		}
		r.Song.Close()
		r.Song = NewSong(r.getFile(r.Queue.Songs[r.Queue.Index]), r.Format)
		r.Song.SetVolume(r.Volume)
		r.Song.OnEnd = r.GetEmbedQueueNext
		r.Song.Play()
	}
}

//Utils
func (r *EmbedRadio) getFile(path string) []byte {
	bytes, err := r.embedFiles.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	return bytes
}