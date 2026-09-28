package wav

import (
	"encoding/binary"
	"io"
)

func WritePCM16(w io.Writer, pcm []int16, sampleRate, channels int) error {
	dataLen := uint32(len(pcm) * 2)
	byteRate := uint32(sampleRate * channels * 2)
	blockAlign := uint16(channels * 2)
	if _, err := w.Write([]byte("RIFF")); err != nil {
		return err
	}
	for _, v := range []any{uint32(36) + dataLen} {
		if err := binary.Write(w, binary.LittleEndian, v); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte("WAVEfmt ")); err != nil {
		return err
	}
	for _, v := range []any{uint32(16), uint16(1), uint16(channels), uint32(sampleRate), byteRate, blockAlign, uint16(16)} {
		if err := binary.Write(w, binary.LittleEndian, v); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte("data")); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, dataLen); err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, pcm)
}
