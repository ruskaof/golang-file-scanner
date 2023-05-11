package queue

type FileMessageQueue interface {
	Produce(filename string) error
	StartConsumer(handle func(filePath string) error) error
}
