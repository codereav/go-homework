package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages { // Запускаем каждый stage, передаём на вход предыдущий результат
		in = execStage(in, done, stage)
	}

	return in
}

func execStage(in In, done In, stageFunc Stage) Out {
	ch := make(Bi)
	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				// Собираем значения из входящего канала в новый канал,
				// чтобы контролировать stage самостоятельно
				ch <- v
			}
		}
	}()

	return stageFunc(ch)
}
