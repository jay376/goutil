package algo

import "io"

func fulfil(bytes []byte, reader io.Reader) error {
	offset, length := 0, len(bytes)
	for {
		n, err := reader.Read(bytes[offset:])
		if n+offset == length {
			return nil
		}
		if err != nil {
			return err
		}
		offset += n
	}
}
